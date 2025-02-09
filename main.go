package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"sync"
)

var receiptStore = make(map[string]Receipt)
var mutex = &sync.Mutex{}

func processReceipt(c *gin.Context) {
	var receipt Receipt
	if err := c.ShouldBindJSON(&receipt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid receipt format"})
		return
	}

	// Generate a UUID for the receipt
	receiptID := uuid.New().String()

	// Store receipt in memory
	mutex.Lock()
	receiptStore[receiptID] = receipt
	mutex.Unlock()

	// Return the generated ID
	c.JSON(http.StatusOK, ProcessReceiptResponse{ID: receiptID})
}

func getPoints(c *gin.Context) {
	id := c.Param("id")

	// Retrieve the receipt from memory
	mutex.Lock()
	receipt, exists := receiptStore[id]
	mutex.Unlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Receipt not found"})
		return
	}

	points := calculatePoints(receipt)
	c.JSON(http.StatusOK, PointsResponse{Points: points})
}

func main() {
	router := gin.Default()

	router.POST("/receipts/process", processReceipt)
	router.GET("/receipts/:id/points", getPoints)

	router.Run(":8080")
}
