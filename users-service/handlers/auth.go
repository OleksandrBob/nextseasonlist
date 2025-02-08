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
)

type AuthHandler struct {
	UserCollection *mongo.Collection
}

func NewAuthHandler(userCollection *mongo.Collection) *AuthHandler {
	return &AuthHandler{UserCollection: userCollection}
}

func (h *AuthHandler) RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingUser models.User
	err := h.UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	hashedPassword, err := utils.GenerateFromPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	user.Password = string(hashedPassword)
	user.ID = primitive.NewObjectID()
	_, err = h.UserCollection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error user inserting into DB"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (h *AuthHandler) LoginUser(c *gin.Context) {
	var userLoginInput models.UserLoginInput
	err := c.ShouldBindJSON(&userLoginInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	var userInDb models.User
	err = h.UserCollection.FindOne(ctx, bson.M{"email": userLoginInput.Email}).Decode(&userInDb)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if !utils.CheckPasswordHash(userLoginInput.Password, userInDb.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	userId := userInDb.ID.Hex()
	accessToken, err := utils.GenerateAccessToken(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate access token"})
		return
	}
	refreshToken, err := utils.GenerateRefreshToken(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate refresh token"})
		return
	}

	c.SetCookie("refreshToken", refreshToken, 7*24*60*60, "", "localhost", true, true)
	c.JSON(http.StatusOK, gin.H{"accessToken": accessToken})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken := ""
	for _, cookie := range c.Request.Cookies() {
		if cookie.Name == "refreshToken" { //TODO && cookie.HttpOnly
			refreshToken = cookie.Value
			break
		}
	}
	if refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token missing"})
		return
	}

	claims, err := utils.ValidateRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	userID := claims["user_id"].(string)
	newAccessToken, err := utils.GenerateAccessToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
	})
}
