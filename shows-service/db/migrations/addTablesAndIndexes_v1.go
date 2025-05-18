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

	_, err := db.GetCollection(db.SerialsCollection).Indexes().CreateMany(ctx,
		[]mongo.IndexModel{
			{
				Keys:    bson.D{{Key: "title", Value: 1}},
				Options: options.Index().SetName("title_idx"),
			},
			{
				Keys:    bson.D{{Key: "categories", Value: 1}},
				Options: options.Index().SetName("categories_idx"),
			}})

	if err != nil {
		return err
	}

	_, err = db.GetCollection(db.CategoriesCollection).Indexes().CreateMany(ctx,
		[]mongo.IndexModel{
			{
				Keys:    bson.D{{Key: "name", Value: 1}},
				Options: options.Index().SetName("name_idx"),
			}})

	if err != nil {
		return err
	}

	return nil
}
