package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/OleksandrBob/nextseasonlist/users-service/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userCollection *mongo.Collection

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		fmt.Println("mongo uri is unset")
		return
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("MongoDB connection error:", err)
		return
	}

	userCollection = client.Database("users_db").Collection("users")

	authHandler := handlers.NewAuthHandler(userCollection)

	router := gin.Default()
	router.Use(gin.Logger())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to Users Service!"})
	})
	router.POST("/register", authHandler.RegisterUser)
	router.POST("/login", loginUser)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server running on port", port)
	router.Run(":" + port)
}

func loginUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "User logged in!"})
}
