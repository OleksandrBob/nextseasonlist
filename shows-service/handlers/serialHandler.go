package handlers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type SerialHandler struct {
	serialsCollection *mongo.Collection
}

func NewSerialHandler(serialsCollection *mongo.Collection) *SerialHandler {
	return &SerialHandler{serialsCollection: serialsCollection}
}

func (h *SerialHandler) SearchSerials(c *gin.Context) {

}
