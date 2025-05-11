package jobs

import (
	"context"
	"log"
	"time"

	"github.com/OleksandrBob/nextseasonlist/users-service/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func StartJobs(tokenBlacklistCollection *mongo.Collection) {
	blacklistTokenTTLIndex := mongo.IndexModel{
		Keys:    bson.M{"expires_at": 1},
		Options: options.Index().SetExpireAfterSeconds(0),
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	_, err := tokenBlacklistCollection.Indexes().CreateOne(ctx, blacklistTokenTTLIndex)
	if err != nil {
		log.Fatal(err)
		return
	}

	blacklist := &blacklistedTokenDeletor{tokenBlacklistCollection}
	blacklistT := time.NewTicker(utils.JobRecuuringPeriod)
	defer blacklistT.Stop()

	for {
		select {
		case <-blacklistT.C:
			log.Println("Triggering blacklisted tokens deletion job")
			go blacklist.deleteExpiredBlacklistedTokens()
		}
	}
}
