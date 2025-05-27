package main

import (
	"fmt"
	"log"
	"os"

	"github.com/OleksandrBob/nextseasonlist/shows-service/db"
	"github.com/OleksandrBob/nextseasonlist/shows-service/db/migrations"
	"github.com/OleksandrBob/nextseasonlist/shows-service/handlers"
	"github.com/OleksandrBob/nextseasonlist/shows-service/middlewares"
	"github.com/OleksandrBob/nextseasonlist/shows-service/utils"
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

	serialsCollection := db.GetCollection(db.SerialsCollection)
	categoriesCollection := db.GetCollection(db.CategoriesCollection)
	//episodesCollection := db.GetCollection(db.EpisodesCollection)

	serialHandler := handlers.NewSerialHandler(serialsCollection, categoriesCollection)

	router := gin.Default()
	serialRoutes := router.Group("/serial", middlewares.AuthMiddleware())
	{
		serialRoutes.POST("/", middlewares.AllowRoleMiddleware(utils.AdminRole), serialHandler.AddSerial)
		serialRoutes.GET("/", serialHandler.SearchSerials)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Shows-Server running on port: ", port)
	router.Run(":" + port)
}
