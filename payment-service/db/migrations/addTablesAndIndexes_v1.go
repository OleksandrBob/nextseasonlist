package migrations

import (
	"context"
	"time"

	"github.com/OleksandrBob/nextseasonlist/payment-service/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Migrate_v1() error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	return db.RunTransaction(ctx, func(sc mongo.SessionContext) error {
		_, err := db.GetCollection(db.PaymentCustomersCollection).Indexes().CreateMany(sc,
			[]mongo.IndexModel{
				{
					Keys:    bson.D{{Key: "stripeCustomerId", Value: 1}},
					Options: options.Index().SetName("stripeCustomerId_idx").SetUnique(true),
				},
				{
					Keys:    bson.D{{Key: "subscriptionId", Value: 1}},
					Options: options.Index().SetName("subscriptionId_idx").SetUnique(false),
				},
				{
					Keys:    bson.D{{Key: "userId", Value: 1}},
					Options: options.Index().SetName("userId_idx").SetUnique(true),
				},
				{
					Keys:    bson.D{{Key: "email", Value: 1}},
					Options: options.Index().SetName("email_idx"),
				},
			})
		if err != nil {
			return err
		}
		return nil
	})
}
