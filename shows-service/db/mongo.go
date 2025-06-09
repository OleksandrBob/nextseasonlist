package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	DbName               = "shows_db"
	SerialsCollection    = "serials"
	EpisodesCollection   = "episodes"
	CategoriesCollection = "categories"
)

// client is singleton and is initialized when program starts
var mongoDbClient *mongo.Client

func ConnectDb(uri string) error {
	ctx, cancell := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancell()

	client, err := mongo.Connect(ctx, options.Client().SetReadPreference(readpref.Secondary()).ApplyURI(uri))
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
	if mongoDbClient != nil {
		mongoDbClient.Disconnect(context.Background())
		log.Println("Disconnected from mongo")
	} else {
		log.Println("Cannot disconect from mongo - client is uninitialized")
	}
}

func GetSession() (mongo.Session, error) {
	sessionOptions := options.Session().SetDefaultReadPreference(readpref.Primary())
	return mongoDbClient.StartSession(sessionOptions)
}

func GetCollection(collectionName string) *mongo.Collection {
	return mongoDbClient.Database(DbName).Collection(collectionName)
}
