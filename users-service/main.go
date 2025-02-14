package main

import (
	"fmt"
	"log"
	"os"

	"github.com/OleksandrBob/nextseasonlist/users-service/db"
	"github.com/OleksandrBob/nextseasonlist/users-service/handlers"
	"github.com/OleksandrBob/nextseasonlist/users-service/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
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

	err = db.ConnectDb(mongoURI)
	if err != nil {
		return
	}

	userCollection = db.GetCollection("users_db", "users")
	authHandler := handlers.NewAuthHandler(userCollection)
	profileHandler := handlers.NewProfileHandler(userCollection)

	router := gin.Default()
	router.Use(gin.Logger())

	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/login", authHandler.LoginUser)
		authRoutes.POST("/register", authHandler.RegisterUser)
		authRoutes.POST("/refreshToken", authHandler.RefreshToken)
	}

	profileRoutes := router.Group("/profile", middlewares.AuthMiddleware())
	{
		profileRoutes.GET("/", profileHandler.GetUserData)
		profileRoutes.PUT("/", profileHandler.UpdatePersonalInfo)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server running on port", port)
	router.Run(":" + port)
}
