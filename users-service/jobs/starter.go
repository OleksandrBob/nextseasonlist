package jobs

import (
	"log"
	"time"

	"github.com/OleksandrBob/nextseasonlist/users-service/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func StartJobs(tokenBlacklistCollection *mongo.Collection) {
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
