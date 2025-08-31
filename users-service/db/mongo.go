package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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

	mongoDbClient.Disconnect(context.Background())
	log.Println("Disconnected from mongo")
}

func GetCollection(collectionName string) *mongo.Collection {
	return mongoDbClient.Database(DbName).Collection(collectionName)
}

func RunTransaction(ctx context.Context, txnFunc func(sessCtx mongo.SessionContext) error) error {
	if mongoDbClient == nil {
		return fmt.Errorf("MongoDB client is not initialized")
	}

	sessionOptions := options.Session().SetDefaultReadPreference(readpref.Primary())
	session, err := mongoDbClient.StartSession(sessionOptions)
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	err = session.StartTransaction()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	sessCtx := mongo.NewSessionContext(ctx, session)
	err = txnFunc(sessCtx)
	if err != nil {
		if abortErr := session.AbortTransaction(sessCtx); abortErr != nil {
			return fmt.Errorf("failed to abort transaction: %w", abortErr)
		}
		return fmt.Errorf("transaction failed: %w", err)
	}

	if commitErr := session.CommitTransaction(sessCtx); commitErr != nil {
		return fmt.Errorf("failed to commit transaction: %w", commitErr)
	}

	return nil
}
