package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/OleksandrBob/nextseasonlist/payment-service/db"
	"github.com/OleksandrBob/nextseasonlist/payment-service/db/migrations"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Println("Warning: .env file not found, using system enviromant variables")
	}

	mongoUri := os.Getenv("MONGO_URI")
	if mongoUri == "" {
		fmt.Println("mongo uri is unset")
		return
	}

	err = db.ConnectDb(mongoUri)
	defer db.DisconnectDb()
	if err != nil {
		log.Fatalf("failed to connect to mongo: %v", err)
		return
	}

	if err = migrations.Migrate_v1(); err != nil {
		log.Println(err.Error())
		return
	}

	router := gin.Default()
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Hello world")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Println("Users-Server running on port: ", port)
	router.Run(":" + port)
}
