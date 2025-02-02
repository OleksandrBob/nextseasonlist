package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func main() {
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

	mongoClient = client
	router := gin.Default()

	router.Use(gin.Logger())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to Users Service!"})
	})

	router.POST("/register", registerUser)
	router.POST("/login", loginUser)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server running on port", port)
	router.Run(":" + port)
}

func registerUser(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "User registered!"})
}

func loginUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "User logged in!"})
}
