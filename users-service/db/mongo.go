package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DbName                      = "users_db"
	UsersCollection             = "users"
	BlacklistedTokensCollection = "blacklistedTokens"
)

var mongoDbClient *mongo.Client

func ConnectDb(uri string) error {
	ctx, cancell := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancell()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("could not connect to mongoDB: %v", err)
		return err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("ping to mongoDB failed: %v", err)
		return err
	}

	mongoDbClient = client
	return nil
}

func DisconnectDb() {
	if mongoDbClient == nil {
		return
	}

	mongoDbClient.Disconnect(context.TODO())
	log.Println("Disconnected from mongo")
}

func GetCollection(collectionName string) *mongo.Collection {
	return mongoDbClient.Database(DbName).Collection(collectionName)
}
