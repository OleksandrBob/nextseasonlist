package handlers

import (
	"net/http"

	"github.com/OleksandrBob/nextseasonlist/payment-service/models"
	"github.com/OleksandrBob/nextseasonlist/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/customersession"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HttpHandler struct {
	PaymentCustomersCollection *mongo.Collection
}

func NewHttpHandler(paymentCustomersCollection *mongo.Collection) *HttpHandler {
	return &HttpHandler{PaymentCustomersCollection: paymentCustomersCollection}
}

func (h *HttpHandler) GetPaymentSession(c *gin.Context) {
	userIdClaim, exists := c.Get(utils.UserIdClaim)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized. Can't access user's Id in context"})
		return
	}

	userId, err := primitive.ObjectIDFromHex(userIdClaim.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var customer models.PaymentCustomer
	err = h.PaymentCustomersCollection.FindOne(c.Request.Context(), bson.M{"userId": userId}).Decode(&customer)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	csp := &stripe.CustomerSessionParams{
		Customer: stripe.String(customer.StripeCustomerID),
		Components: &stripe.CustomerSessionComponentsParams{
			PricingTable: &stripe.CustomerSessionComponentsPricingTableParams{
				Enabled: stripe.Bool(true),
			},
		},
	}

	sess, err := customersession.New(csp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"client_secret": sess.ClientSecret})
}
