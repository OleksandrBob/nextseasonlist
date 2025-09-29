package handlers

import (
	context "context"
	"time"

	"github.com/OleksandrBob/nextseasonlist/payment-service/models"
	paymentpb "github.com/OleksandrBob/nextseasonlist/payment-service/proto/payment"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/customer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GrpcHandler struct {
	paymentpb.UnimplementedPaymentServiceServer
	PaymentCustomersCollection *mongo.Collection
}

func NewGrpcHandler(paymentCollection *mongo.Collection) *GrpcHandler {
	return &GrpcHandler{PaymentCustomersCollection: paymentCollection}
}

func (h *GrpcHandler) CreateStripeCustomer(ctx context.Context, req *paymentpb.CreateStripeCustomerRequest) (*paymentpb.CreateStripeCustomerResponse, error) {
	userID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return &paymentpb.CreateStripeCustomerResponse{Error: "Invalid user ID"}, nil
	}

	var existingCustomer models.PaymentCustomer
	err = h.PaymentCustomersCollection.FindOne(ctx, bson.M{"userId": userID}).Decode(&existingCustomer)
	if err == nil {
		return &paymentpb.CreateStripeCustomerResponse{StripeCustomerId: existingCustomer.StripeCustomerID}, nil
	}
	if err != mongo.ErrNoDocuments {
		return &paymentpb.CreateStripeCustomerResponse{Error: err.Error()}, err
	}

	params := &stripe.CustomerParams{
		Email: stripe.String(req.Email),
		Name:  stripe.String(req.FirstName + " " + req.LastName),
	}
	cust, err := customer.New(params)
	if err != nil {
		return &paymentpb.CreateStripeCustomerResponse{Error: err.Error()}, nil
	}

	now := time.Now().UTC().Unix()
	paymentCustomer := models.PaymentCustomer{
		UserID:           userID,
		Email:            req.Email,
		StripeCustomerID: cust.ID,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	_, insertErr := h.PaymentCustomersCollection.InsertOne(ctx, paymentCustomer)
	if insertErr != nil {
		return &paymentpb.CreateStripeCustomerResponse{Error: insertErr.Error()}, nil
	}

	return &paymentpb.CreateStripeCustomerResponse{StripeCustomerId: cust.ID}, nil
}
