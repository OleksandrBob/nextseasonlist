package main

import (
	"fmt"
	"log"
	"os"

	"github.com/OleksandrBob/nextseasonlist/users-service/db"
	"github.com/OleksandrBob/nextseasonlist/users-service/db/migrations"
	"github.com/OleksandrBob/nextseasonlist/users-service/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	sharedMiddlewares "github.com/OleksandrBob/nextseasonlist/shared/middlewares"
)

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
	defer db.DisconnectDb()
	if err != nil {
		return
	}

	if err = migrations.Migrate_v1(); err != nil {
		log.Println(err.Error())
		return
	}

	userCollection := db.GetCollection(db.UsersCollection)
	tokenBlacklistCollection := db.GetCollection(db.BlacklistedTokensCollection)

	profileHandler := handlers.NewProfileHandler(userCollection)
	authHandler := handlers.NewAuthHandler(userCollection, tokenBlacklistCollection)

	router := gin.Default()

	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/login", authHandler.LoginUser)
		authRoutes.POST("/logout", authHandler.LogOut)
		authRoutes.POST("/register", authHandler.RegisterUser)
		authRoutes.POST("/refreshToken", authHandler.RefreshToken)
	}

	profileRoutes := router.Group("/profile", sharedMiddlewares.AuthMiddleware([]byte(os.Getenv("ACCESS_TOKEN_SECRET"))))
	{
		profileRoutes.GET("/", profileHandler.GetUserData)
		profileRoutes.PUT("/", profileHandler.UpdatePersonalInfo)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// go jobs.StartJobs(tokenBlacklistCollection)

	log.Println("Users-Server running on port: ", port)
	router.Run(":" + port)
}
