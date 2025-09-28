package jobs

import (
	"context"
	"log"
	"time"

	"github.com/OleksandrBob/nextseasonlist/users-service/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type blacklistedTokenDeletor struct {
	tokenBlacklistCollection *mongo.Collection
}

func (d *blacklistedTokenDeletor) deleteExpiredBlacklistedTokens() {
	log.Println("Starting blacklisted tokens deletion")

	ctx, cancel := context.WithTimeout(context.Background(), utils.JobMaxDurationTime)
	defer cancel()

	_, err := d.tokenBlacklistCollection.DeleteMany(ctx, bson.D{{Key: "expires_at", Value: bson.D{{Key: "$lt", Value: time.Now().UTC()}}}})

	if err != nil {
		log.Println("Failed to cleanup expired tokens:", err)
	} else {
		log.Println("Expired tokens cleaned up successfully!")
	}
}
