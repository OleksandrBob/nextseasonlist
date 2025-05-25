package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/OleksandrBob/nextseasonlist/users-service/models"
	"github.com/OleksandrBob/nextseasonlist/users-service/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProfileHandler struct {
	UserCollection *mongo.Collection
}

func NewProfileHandler(userCollection *mongo.Collection) *ProfileHandler {
	return &ProfileHandler{UserCollection: userCollection}
}

func (h *ProfileHandler) GetUserData(c *gin.Context) {
	userId, exists := c.Get(utils.UserIdClaim)

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User is unauthorized"})
		return
	}

	id, err := primitive.ObjectIDFromHex(userId.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userId in token"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var userInDb models.User
	err = h.UserCollection.FindOne(ctx, bson.M{"_id": id}, options.FindOne().SetProjection(bson.D{{Key: "password", Value: 0}})).Decode(&userInDb)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to decode user from DB"})
	}

	c.JSON(http.StatusOK, userInDb)
}

func (h *ProfileHandler) UpdatePersonalInfo(c *gin.Context) {
	var updateInfo models.UserUpdateInfo
	err := c.ShouldBindJSON(&updateInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	ctx, cancelGetFromDb := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelGetFromDb()

	var userInDb models.User
	err = h.UserCollection.FindOne(ctx, bson.D{{Key: "_id", Value: updateInfo.ID}}).Decode(&userInDb)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User doesn't exist"})
		return
	}
	cancelGetFromDb()

	if updateInfo.FirstName != "" {
		userInDb.FirstName = updateInfo.FirstName
	}
	if updateInfo.LastName != "" {
		userInDb.LastName = updateInfo.LastName
	}

	ctx, cancelUpdate := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelUpdate()

	_, err = h.UserCollection.UpdateOne(ctx, bson.D{{Key: "_id", Value: updateInfo.ID}}, bson.D{{Key: "$set", Value: bson.D{{Key: "firstName", Value: updateInfo.FirstName}, {Key: "lastName", Value: updateInfo.LastName}}}})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Update failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Update succeded"})
}
