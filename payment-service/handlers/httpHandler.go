package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type HttpHandler struct {
	PaymentCustomersCollection *mongo.Collection
}

func NewHttpHandler(paymentCustomersCollection *mongo.Collection) *HttpHandler {
	return &HttpHandler{PaymentCustomersCollection: paymentCustomersCollection}
}

func (h *HttpHandler) GetPaymentSession(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{"error": "baba"})
}
