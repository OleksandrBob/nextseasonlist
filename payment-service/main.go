package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/OleksandrBob/nextseasonlist/payment-service/db"
	"github.com/OleksandrBob/nextseasonlist/payment-service/db/migrations"
	"github.com/OleksandrBob/nextseasonlist/payment-service/handlers"
	paymentpb "github.com/OleksandrBob/nextseasonlist/payment-service/proto/payment"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
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

	// Start GRPC server
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "8083"
	}
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	paymentHandler := handlers.NewPaymentHandler(db.GetCollection(db.PaymentCustomersCollection))
	paymentpb.RegisterPaymentServiceServer(grpcServer, paymentHandler)
	log.Printf("Payment GRPC server listening on port %s", grpcPort)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Optionally keep HTTP server for health checks
	router := gin.Default()
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Hello world")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Println("Payment-Server running on port: ", port)
	router.Run(":" + port)
}
