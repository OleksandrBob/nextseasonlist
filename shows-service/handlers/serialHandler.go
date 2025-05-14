package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/OleksandrBob/nextseasonlist/shows-service/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SerialHandler struct {
	serialsCollection *mongo.Collection
}

func NewSerialHandler(serialsCollection *mongo.Collection) *SerialHandler {
	return &SerialHandler{serialsCollection: serialsCollection}
}

func (h *SerialHandler) SearchSerials(c *gin.Context) {
	var searchRequest models.SearchSerialsQuery
	err := c.ShouldBindJSON(&searchRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid search request"})
		return
	}

}

func (h *SerialHandler) AddSerial(c *gin.Context) {
	var addCommand models.AddSerialCommand
	err := c.ShouldBindJSON(&addCommand)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	serialToAdd := models.Serial{
		ID:          primitive.NewObjectID(),
		Name:        addCommand.Name,
		Description: addCommand.Description,
		Categories:  addCommand.Categories,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := h.serialsCollection.InsertOne(ctx, serialToAdd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("Inserted document with _id: ", result.InsertedID)
	c.JSON(http.StatusCreated, gin.H{})
}
