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
		log.Println("Warning: .env file not found, using system enviromant variables")
	}

	// stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	mongoUri := os.Getenv("MONGO_URI")
	if mongoUri == "" {
		fmt.Println("mongo uri is unset")
		return
	}

	// err = db.ConnectDb(mongoUri)
	// defer db.DisconnectDb()
	// if err != nil {
	// 	log.Fatalf("failed to connect to mongo: %v", err)
	// 	return
	// }

	// if err = migrations.Migrate_v1(); err != nil {
	// 	log.Println(err.Error())
	// 	return
	// }

	// pcc := db.GetCollection(db.PaymentCustomersCollection)

	// grpcPort := os.Getenv("GRPC_PORT")
	// if grpcPort == "" {
	// 	grpcPort = "8083"
	// }
	// lis, err := net.Listen("tcp", ":"+grpcPort)
	// if err != nil {
	// 	log.Fatalf("failed to listen: %v", err)
	// }
	// grpcServer := grpc.NewServer()
	// grpcHandler := handlers.NewGrpcHandler(pcc)
	// paymentpb.RegisterPaymentServiceServer(grpcServer, grpcHandler)
	// log.Printf("Payment GRPC server listening on port %s", grpcPort)
	// go func() {
	// 	if err := grpcServer.Serve(lis); err != nil {
	// 		log.Fatalf("failed to serve: %v", err)
	// 	}
	// }()

	router := gin.Default()
	// httpHandler := handlers.NewHttpHandler(pcc)
	// webhookHandler := handlers.NewWebhookHandler(pcc)

	// clientRoutes := router.Group("/client", sharedMiddlewares.AuthMiddleware([]byte(os.Getenv("ACCESS_TOKEN_SECRET"))))
	// {
	// 	clientRoutes.GET("/payment-session", httpHandler.GetPaymentSession)
	// 	clientRoutes.GET("/subscription-status/:customerId", httpHandler.GetCustomerSubscriptionStatus)
	// }

	// webhookRoutes := router.Group("/webhook")
	// {
	// 	webhookRoutes.POST("/stripe", webhookHandler.HandleStripeWebhook)
	// }

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Hello world for payment (2)")
	})

	router.GET("/google-check", func(c *gin.Context) {
		resp, err := http.Get("https://www.google.com")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"accessed": false, "error": err.Error()})
			return
		}
		defer resp.Body.Close()
		c.JSON(http.StatusOK, gin.H{"accessed": resp.StatusCode == http.StatusOK})
	})

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8082"
	}

	log.Println("Payment-Server running on port: ", httpPort)
	router.Run(":" + httpPort)
}
