package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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
	hand := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello Sambo"})
	}

	router := gin.Default()
	router.GET("/", hand)

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
		port = "80"
	}

	log.Println("Users-Server running on port: ", port)
	router.Run(":" + port)
}
