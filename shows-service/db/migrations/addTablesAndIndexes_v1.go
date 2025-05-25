package migrations

import (
	"context"
	"time"

	"github.com/OleksandrBob/nextseasonlist/shows-service/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Migrate_v1() error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	mongoSession, err := db.GetSession()
	if err != nil {
		return err
	}

	defer mongoSession.EndSession(ctx)

	err = mongo.WithSession(ctx, mongoSession, func(sc mongo.SessionContext) error {
		if err := mongoSession.StartTransaction(); err != nil {
			return err
		}

		_, err := db.GetCollection(db.SerialsCollection).Indexes().CreateMany(sc, []mongo.IndexModel{
			{
				Keys:    bson.D{{Key: "title", Value: 1}},
				Options: options.Index().SetName("title_idx"),
			},
			{
				Keys:    bson.D{{Key: "categories", Value: 1}},
				Options: options.Index().SetName("categories_idx"),
			},
		})

		if err != nil {
			_ = mongoSession.AbortTransaction(sc)
			return err
		}

		_, err = db.GetCollection(db.CategoriesCollection).Indexes().CreateMany(sc, []mongo.IndexModel{
			{
				Keys:    bson.D{{Key: "name", Value: 1}},
				Options: options.Index().SetName("name_idx"),
			},
		})

		if err != nil {
			_ = mongoSession.AbortTransaction(sc)
			return err
		}

		if err := mongoSession.CommitTransaction(sc); err != nil {
			return err
		}

		return nil
	})

	return err
}
