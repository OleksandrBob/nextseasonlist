package handlers

import (
	context "context"
	"os"

	paymentpb "github.com/OleksandrBob/nextseasonlist/payment-service/proto/payment"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/customer"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentHandler struct {
	paymentpb.UnimplementedPaymentServiceServer
	PaymentCollection *mongo.Collection
}

func NewPaymentHandler(paymentCollection *mongo.Collection) *PaymentHandler {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	return &PaymentHandler{PaymentCollection: paymentCollection}
}

func (h *PaymentHandler) CreateStripeCustomer(ctx context.Context, req *paymentpb.CreateStripeCustomerRequest) (*paymentpb.CreateStripeCustomerResponse, error) {
	params := &stripe.CustomerParams{
		Email: stripe.String(req.Email),
		Name:  stripe.String(req.FirstName + " " + req.LastName),
	}
	cust, err := customer.New(params)
	if err != nil {
		return &paymentpb.CreateStripeCustomerResponse{Error: err.Error()}, nil
	}
	return &paymentpb.CreateStripeCustomerResponse{StripeCustomerId: cust.ID}, nil
}
