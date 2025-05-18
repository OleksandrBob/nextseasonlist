package migrations

import (
	"context"
	"time"

	"github.com/OleksandrBob/nextseasonlist/users-service/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Migrate_v1() error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := db.GetCollection(db.UsersCollection).Indexes().CreateMany(ctx,
		[]mongo.IndexModel{
			{
				Keys:    bson.D{{Key: "email", Value: 1}},
				Options: options.Index().SetName("email_idx"),
			}})

	if err != nil {
		return err
	}

	_, err = db.GetCollection(db.BlacklistedTokensCollection).Indexes().CreateMany(ctx,
		[]mongo.IndexModel{
			{
				Keys:    bson.D{{Key: "expires_at", Value: 1}},
				Options: options.Index().SetName("expires_at_idx").SetExpireAfterSeconds(0),
			},
			{
				Keys:    bson.D{{Key: "token", Value: 1}},
				Options: options.Index().SetName("token_idx"),
			}})

	if err != nil {
		return err
	}

	return nil
}
