package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		fmt.Println("mongo uri is unset")
		//return
	}
	fmt.Println("Mongo URI:", mongoURI)

	// err = db.ConnectDb(mongoURI)
	// defer db.DisconnectDb()
	// if err != nil {
	// 	return
	// }

	// if err = migrations.Migrate_v1(); err != nil {
	// 	log.Println(err.Error())
	// 	return
	// }

	// userCollection := db.GetCollection(db.UsersCollection)
	// tokenBlacklistCollection := db.GetCollection(db.BlacklistedTokensCollection)

	// profileHandler := handlers.NewProfileHandler(userCollection)
	// authHandler := handlers.NewAuthHandler(userCollection, tokenBlacklistCollection)

	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(1000000)

	hand := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":      "Hello Sambo",
			"randomNumber": randomNumber,
		})
	}

	router := gin.Default()
	router.GET("/", hand)

	paymentServUri := os.Getenv("PAYMENT_SERVICE_URI")

	router.GET("/payment-check", func(c *gin.Context) {
		resp, err := http.Get(paymentServUri + "/google-check")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"accessed": false, "error": err.Error()})
			return
		}
		defer resp.Body.Close()
		fmt.Println("successfull execution in users service")

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"accessed": false, "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"accessed": resp.StatusCode == http.StatusOK,
			"body":     string(bodyBytes),
		})
	})

	router.GET("/health", func(c *gin.Context) {
		fmt.Println("Triggered health check in users service")

		c.JSON(http.StatusOK, gin.H{
			"body": "all is fine in users service",
		})
	})

	// authRoutes := router.Group("/auth")
	// {
	// 	authRoutes.POST("/login", authHandler.LoginUser)
	// 	authRoutes.POST("/logout", authHandler.LogOut)
	// 	authRoutes.POST("/register", authHandler.RegisterUser)
	// 	authRoutes.POST("/refreshToken", authHandler.RefreshToken)
	// }

	// profileRoutes := router.Group("/profile", sharedMiddlewares.AuthMiddleware([]byte(os.Getenv("ACCESS_TOKEN_SECRET"))))
	// {
	// 	profileRoutes.GET("/", profileHandler.GetUserData)
	// 	profileRoutes.PUT("/", profileHandler.UpdatePersonalInfo)
	// }

	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("Port was unset")
		port = "8080"
	}

	log.Println("Users-Server running on port: ", port)
	router.Run(":" + port)
}
