package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WebhookHandler struct {
	PaymentCustomersCollection *mongo.Collection
	endpointSecret             string
}

func NewWebhookHandler(paymentCustomersCollection *mongo.Collection) *WebhookHandler {

	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	if endpointSecret == "" {
		log.Println("STRIPE_WEBHOOK_SECRET environment variable not set")
		panic("STRIPE_WEBHOOK_SECRET environment variable not set")
	}

	return &WebhookHandler{PaymentCustomersCollection: paymentCustomersCollection, endpointSecret: endpointSecret}
}

func (h *WebhookHandler) HandleStripeWebhook(c *gin.Context) {
	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Error reading request body"})
		return
	}

	event, err := webhook.ConstructEvent(body, c.GetHeader("Stripe-Signature"), h.endpointSecret)
	if err != nil {
		log.Printf("Webhook signature verification failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signature"})
		return
	}
	log.Printf("Constructed event: %+v", event)

	switch event.Type {
	case "customer.subscription.created":
		h.handleSubscriptionCreated(c.Request.Context(), event)
	case "customer.subscription.updated":
		h.handleSubscriptionUpdated(c.Request.Context(), event)
	case "customer.subscription.deleted":
		h.handleSubscriptionDeleted(c.Request.Context(), event)
	case "invoice.payment_succeeded":
		h.handlePaymentSucceeded(c.Request.Context(), event)
	case "invoice.payment_failed":
		h.handlePaymentFailed(c.Request.Context(), event)
	default:
		log.Printf("Unhandled event type: %s", event.Type)
	}

	c.JSON(http.StatusOK, gin.H{"received": true})
}

func (h *WebhookHandler) handleSubscriptionCreated(ctx context.Context, event stripe.Event) {
	var subscription stripe.Subscription
	err := json.Unmarshal(event.Data.Raw, &subscription)
	if err != nil {
		log.Printf("Error parsing subscription: %v", err)
		return
	}

	log.Printf("Subscription created: %s for customer: %s", subscription.ID, subscription.Customer.ID)

	h.updateCustomerSubscription(ctx, subscription)
}

func (h *WebhookHandler) handleSubscriptionUpdated(ctx context.Context, event stripe.Event) {
	var subscription stripe.Subscription
	err := json.Unmarshal(event.Data.Raw, &subscription)
	if err != nil {
		log.Printf("Error parsing subscription: %v", err)
		return
	}

	log.Printf("Subscription updated: %s for customer: %s, status: %s",
		subscription.ID, subscription.Customer.ID, subscription.Status)

	h.updateCustomerSubscription(ctx, subscription)
}

func (h *WebhookHandler) handleSubscriptionDeleted(ctx context.Context, event stripe.Event) {
	var subscription stripe.Subscription
	err := json.Unmarshal(event.Data.Raw, &subscription)
	if err != nil {
		log.Printf("Error parsing subscription: %v", err)
		return
	}

	log.Printf("Subscription deleted: %s for customer: %s", subscription.ID, subscription.Customer.ID)

	filter := bson.M{"stripeCustomerId": subscription.Customer.ID}
	update := bson.M{
		"$set": bson.M{
			"subscriptionId":     "",
			"planId":             "",
			"subscriptionStatus": "canceled",
			"updatedAt":          time.Now().Unix(),
		},
	}

	_, err = h.PaymentCustomersCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error updating customer subscription: %v", err)
	}
}

func (h *WebhookHandler) handlePaymentSucceeded(ctx context.Context, event stripe.Event) {
	var invoice stripe.Invoice
	err := json.Unmarshal(event.Data.Raw, &invoice)
	if err != nil {
		log.Printf("Error parsing invoice: %v", err)
		return
	}

	// Only handle subscription invoices
	if invoice.BillingReason != stripe.InvoiceBillingReasonSubscription && // TODO: move this check to a separate method
		invoice.BillingReason != stripe.InvoiceBillingReasonSubscriptionCreate &&
		invoice.BillingReason != stripe.InvoiceBillingReasonSubscriptionCycle &&
		invoice.BillingReason != stripe.InvoiceBillingReasonSubscriptionThreshold &&
		invoice.BillingReason != stripe.InvoiceBillingReasonSubscriptionUpdate {
		return
	}

	log.Printf("Payment succeeded for customer: %s",
		invoice.Customer.ID)

	// Update subscription status to active if it was past_due
	filter := bson.M{"stripeCustomerId": invoice.Customer.ID}
	update := bson.M{
		"$set": bson.M{
			"subscriptionStatus": "active",
			"updatedAt":          time.Now().Unix(),
		},
	}

	_, err = h.PaymentCustomersCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error updating customer subscription status: %v", err)
	}
}

func (h *WebhookHandler) handlePaymentFailed(ctx context.Context, event stripe.Event) {
	var invoice stripe.Invoice
	err := json.Unmarshal(event.Data.Raw, &invoice)
	if err != nil {
		log.Printf("Error parsing invoice: %v", err)
		return
	}

	// Only handle subscription invoices
	if invoice.BillingReason != stripe.InvoiceBillingReasonSubscription &&
		invoice.BillingReason != stripe.InvoiceBillingReasonSubscriptionCreate &&
		invoice.BillingReason != stripe.InvoiceBillingReasonSubscriptionCycle &&
		invoice.BillingReason != stripe.InvoiceBillingReasonSubscriptionThreshold &&
		invoice.BillingReason != stripe.InvoiceBillingReasonSubscriptionUpdate {
		return
	}

	log.Printf("Payment failed for customer: %s",
		invoice.Customer.ID)

	// Update subscription status to past_due
	filter := bson.M{"stripeCustomerId": invoice.Customer.ID}
	update := bson.M{
		"$set": bson.M{
			"subscriptionStatus": "past_due",
			"updatedAt":          time.Now().Unix(),
		},
	}

	_, err = h.PaymentCustomersCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error updating customer subscription status: %v", err)
	}
}

func (h *WebhookHandler) updateCustomerSubscription(ctx context.Context, subscription stripe.Subscription) {
	var planId string
	var usagePlan int16

	if len(subscription.Items.Data) > 0 {
		planId = subscription.Items.Data[0].Price.ID

		switch planId {
		case "price_1SCH6JLXxKn9DoPxSLsZLqLf": // TODO: move to env vars
			usagePlan = 1 // TODO: make consts (like enum)
		case "price_1SCH6YLXxKn9DoPxZEL9bALL":
			usagePlan = 2
		default:
			usagePlan = 0
		}
	}

	// Get payment method ID if available
	var paymentMethodId string
	if subscription.DefaultPaymentMethod != nil {
		paymentMethodId = subscription.DefaultPaymentMethod.ID
	}

	filter := bson.M{"stripeCustomerId": subscription.Customer.ID}
	update := bson.M{
		"$set": bson.M{
			"subscriptionId":     subscription.ID,
			"planId":             planId,
			"paymentMethodId":    paymentMethodId,
			"subscriptionStatus": string(subscription.Status),
			"usagePlan":          usagePlan,
			"updatedAt":          time.Now().Unix(),
		},
	}

	opts := options.Update().SetUpsert(true)
	result, err := h.PaymentCustomersCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Printf("Error updating customer subscription: %v", err)
		return
	}

	if result.UpsertedCount > 0 {
		log.Printf("Created new customer record for Stripe customer: %s", subscription.Customer.ID)
	} else {
		log.Printf("Updated customer record for Stripe customer: %s", subscription.Customer.ID)
	}
}
