package main

import (
	"fmt"
	"log"
	"os"

	"github.com/OleksandrBob/nextseasonlist/shows-service/db"
	"github.com/OleksandrBob/nextseasonlist/shows-service/db/migrations"
	"github.com/OleksandrBob/nextseasonlist/shows-service/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	sharedMiddlewares "github.com/OleksandrBob/nextseasonlist/shared/middlewares"
	sharedUtils "github.com/OleksandrBob/nextseasonlist/shared/utils"
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

	accessTokenSecret := []byte(os.Getenv("ACCESS_TOKEN_SECRET"))

	serialsCollection := db.GetCollection(db.SerialsCollection)
	episodesCollection := db.GetCollection(db.EpisodesCollection)
	categoriesCollection := db.GetCollection(db.CategoriesCollection)

	router := gin.Default()

	serialHandler := handlers.NewSerialHandler(serialsCollection, categoriesCollection)
	serialRoutes := router.Group("/serial", sharedMiddlewares.AuthMiddleware(accessTokenSecret))
	{
		serialRoutes.POST("/", sharedMiddlewares.AllowRoleMiddleware(sharedUtils.AdminRole), serialHandler.AddSerial)
		serialRoutes.POST("/search", serialHandler.SearchSerials)
	}

	episodesHandler := handlers.NewEpisodeHandler(episodesCollection, serialsCollection)
	episodesRoutes := router.Group("/episode", sharedMiddlewares.AuthMiddleware(accessTokenSecret))
	{
		episodesRoutes.GET("/:id", episodesHandler.GetEpisodeById)
		episodesRoutes.POST("/", episodesHandler.AddEpisode)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Println("Shows-Server running on port: ", port)
	router.Run(":" + port)
}
