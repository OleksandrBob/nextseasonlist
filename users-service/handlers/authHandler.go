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
	UserCollection           *mongo.Collection
	TokenBlacklistCollection *mongo.Collection
}

func NewAuthHandler(userCollection *mongo.Collection, tokenBlacklistCollection *mongo.Collection) *AuthHandler {
	return &AuthHandler{UserCollection: userCollection, TokenBlacklistCollection: tokenBlacklistCollection}
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

	user.Password = hashedPassword
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

	c.SetCookie(utils.RefreshTokenName, refreshToken, 7*24*60*60, "", "localhost", true, true) //TODO remove localhost
	c.JSON(http.StatusOK, gin.H{utils.AccessTokenName: accessToken})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken := ""
	for _, cookie := range c.Request.Cookies() {
		if cookie.Name == utils.RefreshTokenName { //TODO && cookie.HttpOnly
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

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	var blacklistedToken models.BlacklistedToken
	err = h.TokenBlacklistCollection.FindOne(ctx, bson.M{"token": refreshToken}).Decode(&blacklistedToken)
	if err == nil { // found in db
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User was loged out"})
		return
	}
	if err.Error() != "mongo: no documents in result" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to decode value from db"})
		return
	}

	userID := claims[utils.UserIdClaim].(string)
	newAccessToken, err := utils.GenerateAccessToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		utils.AccessTokenName: newAccessToken,
	})
}

func (h *AuthHandler) LogOut(c *gin.Context) {
	refreshToken, err := c.Cookie(utils.RefreshTokenName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing refresh token. Can't perform log out"})
		return
	}

	claims, err := utils.ValidateRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	tokenExpirationTime := int64(claims[utils.ExpirationClaim].(float64))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	h.TokenBlacklistCollection.InsertOne(ctx, models.BlacklistedToken{
		Token:     refreshToken,
		ExpiresAt: time.Unix(tokenExpirationTime, 0),
	})

	c.SetCookie(utils.RefreshTokenName, "", -1, "", "localhost", true, true) //TODO remove localhost
}
