package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/OleksandrBob/nextseasonlist/shows-service/db"
	"github.com/OleksandrBob/nextseasonlist/shows-service/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type episodeHandler struct {
	episodesCollection *mongo.Collection
	serialsCollection  *mongo.Collection
}

func NewEpisodeHandler(episodesCollection *mongo.Collection, serialsCollection *mongo.Collection) *episodeHandler {
	return &episodeHandler{episodesCollection: episodesCollection, serialsCollection: serialsCollection}
}

func (h *episodeHandler) AddEpisode(c *gin.Context) {
	var aec models.AddEpisodeCommand
	err := c.ShouldBindJSON(&aec)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newEpId := primitive.NewObjectID()

	err = db.RunTransaction(ctx, func(sc mongo.SessionContext) error {
		e := models.Episode{
			ID:          newEpId,
			Name:        aec.Name,
			Season:      aec.Season,
			Number:      aec.Number,
			SerialId:    aec.SerialId,
			ReleaseDate: aec.ReleaseDate,
		}

		if _, err := h.episodesCollection.InsertOne(sc, e); err != nil {
			return err
		}

		var s models.Serial
		err := h.serialsCollection.FindOne(sc, bson.M{"_id": aec.SerialId}).Decode(&s)
		if err != nil {
			return fmt.Errorf("could not find serial with specified id: %w", err)
		}

		if s.Seasons < aec.Season {
			s.Seasons = aec.Season
			_, err = h.serialsCollection.ReplaceOne(sc, bson.M{"_id": s.ID}, s)
			if err != nil {
				return fmt.Errorf("failed to update serial: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": newEpId})
}

func (h *episodeHandler) GetEpisodeById(c *gin.Context) {
	id := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var e models.Episode
	err = h.episodesCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&e)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "episiode not found"})
		return
	}

	c.JSON(http.StatusOK, e)
}
