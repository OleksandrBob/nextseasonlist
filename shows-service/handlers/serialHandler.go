package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/OleksandrBob/nextseasonlist/shows-service/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SerialHandler struct {
	serialsCollection    *mongo.Collection
	categoriesCollection *mongo.Collection
}

func NewSerialHandler(serialsCollection *mongo.Collection, categoriesCollection *mongo.Collection) *SerialHandler {
	return &SerialHandler{serialsCollection: serialsCollection, categoriesCollection: categoriesCollection}
}

func (h *SerialHandler) SearchSerials(c *gin.Context) {
	var searchRequest models.SearchSerialsQuery
	err := c.ShouldBindJSON(&searchRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid search request"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	serialsCursor, err := h.serialsCollection.Find(ctx, bson.M{"$or": []interface{}{bson.M{"title": bson.M{"$regex": searchRequest.Title}}, bson.M{"categories": bson.M{"$in": searchRequest.Categories}}}})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error receiving serials from DB"})
		return
	}

	var foundSerials []models.SearchSerialsQueryResponse
	err = serialsCursor.All(ctx, &foundSerials)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error receiving serials"})
		return
	}

	c.JSON(http.StatusOK, foundSerials)
}

func (h *SerialHandler) AddSerial(c *gin.Context) {
	var addCommand models.AddSerialCommand
	err := c.ShouldBindJSON(&addCommand)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if len(addCommand.Categories) > 0 {
		categoriesCursor, err := h.categoriesCollection.Find(ctx, bson.M{"name": bson.M{"$in": addCommand.Categories}}, options.Find().SetProjection(bson.D{{Key: "_id", Value: 0}}))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error receiving categories from DB"})
			return
		}

		var foundCategories []models.Category
		err = categoriesCursor.All(ctx, &foundCategories)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error checking categories"})
			return
		}

		if len(foundCategories) != len(addCommand.Categories) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Used unexisting category"})
			return
		}
	}

	serialToAdd := models.Serial{
		ID:          primitive.NewObjectID(),
		Title:       addCommand.Title,
		Description: addCommand.Description,
		Categories:  addCommand.Categories,
	}

	result, err := h.serialsCollection.InsertOne(ctx, serialToAdd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("Inserted serial with _id: ", result.InsertedID)
	c.JSON(http.StatusCreated, gin.H{"CreatedSerialId": result.InsertedID})
}
