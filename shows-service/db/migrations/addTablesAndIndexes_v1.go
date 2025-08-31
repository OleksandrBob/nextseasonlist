package migrations

import (
	"context"
	"time"

	"github.com/OleksandrBob/nextseasonlist/shows-service/db"
	"github.com/OleksandrBob/nextseasonlist/shows-service/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Migrate_v1() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return db.RunTransaction(ctx, func(sc mongo.SessionContext) error {
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
			return err
		}

		_, err = db.GetCollection(db.CategoriesCollection).Indexes().CreateMany(sc, []mongo.IndexModel{
			{
				Keys:    bson.D{{Key: "name", Value: 1}},
				Options: options.Index().SetName("name_idx").SetUnique(true),
			},
		})
		if err != nil {
			return err
		}

		_, err = db.GetCollection(db.EpisodesCollection).Indexes().CreateMany(sc, []mongo.IndexModel{
			{
				Keys:    bson.D{{Key: "serialId", Value: 1}},
				Options: options.Index().SetName("serialId_idx"),
			},
		})
		if err != nil {
			return err
		}

		catCount, err := db.GetCollection(db.CategoriesCollection).CountDocuments(sc, bson.D{})
		if err != nil {
			return err
		}

		if catCount == 0 {
			_, err = db.GetCollection(db.CategoriesCollection).InsertMany(sc, []interface{}{
				models.Category{ID: 1, Name: "comedy"},
				models.Category{ID: 2, Name: "drama"},
				models.Category{ID: 3, Name: "horror"},
				models.Category{ID: 4, Name: "adventure"},
				models.Category{ID: 5, Name: "survival"},
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
}
