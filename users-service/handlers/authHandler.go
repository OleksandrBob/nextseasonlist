package handlers

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	paymentpb "github.com/OleksandrBob/nextseasonlist/users-service/proto/payment"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/OleksandrBob/nextseasonlist/users-service/models"
	"github.com/OleksandrBob/nextseasonlist/users-service/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	sharedUtils "github.com/OleksandrBob/nextseasonlist/shared/utils"
)

type AuthHandler struct {
	UserCollection           *mongo.Collection
	TokenBlacklistCollection *mongo.Collection
}

var HostUrl string = os.Getenv("HOST_URL")

func NewAuthHandler(userCollection *mongo.Collection, tokenBlacklistCollection *mongo.Collection) *AuthHandler {
	return &AuthHandler{UserCollection: userCollection, TokenBlacklistCollection: tokenBlacklistCollection}
}

func (h *AuthHandler) RegisterUser(c *gin.Context) {
	var registerUserDto models.RegisterUserCommand
	if err := c.ShouldBindJSON(&registerUserDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "description": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var existingUser models.User
	err := h.UserCollection.FindOne(ctx, bson.M{"email": registerUserDto.Email}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	psUrl := os.Getenv("PAYMENT_SERVICE_GRPC")
	log.Println(psUrl)

	conn, err := grpc.NewClient(psUrl, grpc.WithTransportCredentials(insecure.NewCredentials())) // TODO: check what is this
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not connect to payment service"})
		return
	}
	defer conn.Close()

	client := paymentpb.NewPaymentServiceClient(conn)
	grpcCtx, grpcCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer grpcCancel()

	resp, err := client.CreateStripeCustomer(grpcCtx, &paymentpb.CreateStripeCustomerRequest{
		Email:     registerUserDto.Email,
		FirstName: registerUserDto.FirstName,
		LastName:  registerUserDto.LastName,
	})
	if err != nil || resp.Error != "" {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create stripe customer"})
		return
	}

	stripeCustomerId := resp.StripeCustomerId

	hashedPassword, err := utils.GenerateFromPassword(registerUserDto.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	userToCreate := models.User{
		ID:               primitive.NewObjectID(),
		Email:            registerUserDto.Email,
		FirstName:        registerUserDto.FirstName,
		LastName:         registerUserDto.LastName,
		Roles:            []string{sharedUtils.UserRole},
		UsagePlan:        sharedUtils.FreePlan,
		StripeCustomerId: stripeCustomerId,
		Password:         hashedPassword,
	}

	_, err = h.UserCollection.InsertOne(ctx, userToCreate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error user inserting into DB"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "stripeCustomerId": stripeCustomerId})
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
	accessToken, err := utils.GenerateAccessToken(userId, userInDb.Roles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate access token"})
		return
	}
	refreshToken, err := utils.GenerateRefreshToken(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate refresh token"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(utils.RefreshTokenName, refreshToken, 7*24*60*60, "", HostUrl, true, true)
	c.JSON(http.StatusOK, gin.H{utils.AccessTokenName: accessToken})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken := ""
	for _, cookie := range c.Request.Cookies() {
		if cookie.Name == utils.RefreshTokenName {
			refreshToken = cookie.Value
			break
		}
	}
	if refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token missing"})
		return
	}

	refreshTokenClaims, err := utils.ValidateRefreshToken(refreshToken)
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	userId := refreshTokenClaims[utils.UserIdClaim].(string)

	mongoId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userId in token"})
		return
	}

	var userInDb models.User
	err = h.UserCollection.FindOne(ctx, bson.M{"_id": mongoId}, options.FindOne().SetProjection(bson.D{{Key: "password", Value: 0}})).Decode(&userInDb)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not found user for refresh"})
		return
	}

	newAccessToken, err := utils.GenerateAccessToken(userId, userInDb.Roles)

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

	c.SetCookie(utils.RefreshTokenName, "", -1, "", HostUrl, true, true)
}
