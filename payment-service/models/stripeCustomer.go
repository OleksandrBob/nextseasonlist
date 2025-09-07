package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type PaymentCustomer struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StripeCustomerID   string             `bson:"stripeCustomerId" json:"stripeCustomerId"`
	SubscriptionID     string             `bson:"subscriptionId" json:"subscriptionId"`
	PlanID             string             `bson:"planId" json:"planId"`
	PaymentMethodID    string             `bson:"paymentMethodId" json:"paymentMethodId"`
	SubscriptionStatus string             `bson:"subscriptionStatus" json:"subscriptionStatus"`
	CurrentPeriodEnd   int64              `bson:"currentPeriodEnd" json:"currentPeriodEnd"`
	TrialEnd           int64              `bson:"trialEnd" json:"trialEnd"`
	CreatedAt          int64              `bson:"createdAt" json:"createdAt"`
	UpdatedAt          int64              `bson:"updatedAt" json:"updatedAt"`
}
