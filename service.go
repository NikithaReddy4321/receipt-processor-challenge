package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// To calculatePoints
func calculatePoints(receipt Receipt) (int, string) {
	points := 0
	var breakdownLines []string

	// 1. One point for every alphanumeric character in the retailer name
	alnumRegex := regexp.MustCompile("[a-zA-Z0-9]")
	retailerPoints := len(alnumRegex.FindAllString(receipt.Retailer, -1))
	points += retailerPoints
	breakdownLines = append(breakdownLines, fmt.Sprintf("     %d points - retailer name has %d characters", retailerPoints, retailerPoints))

	total, _ := strconv.ParseFloat(receipt.Total, 64)

	// 2. 50 points if the total is a round dollar amount
	if math.Mod(total, 1.00) == 0 {
		points += 50
		breakdownLines = append(breakdownLines, "    50 points - total is a round dollar amount")
	}

	// 3. 25 points if the total is a multiple of 0.25
	if math.Mod(total, 0.25) == 0 {
		points += 25
		breakdownLines = append(breakdownLines, "    25 points - total is a multiple of 0.25")
	}

	// 4. 5 points for every two items on the receipt
	itemPairs := len(receipt.Items) / 2
	itemPoints := itemPairs * 5
	points += itemPoints
	breakdownLines = append(breakdownLines, fmt.Sprintf("    10 points - %d items (%d pairs @ 5 points each)", len(receipt.Items), itemPairs))

	// 5. Points based on item description length being a multiple of 3
	for _, item := range receipt.Items {
		trimmedDesc := strings.TrimSpace(item.ShortDescription)
		if len(trimmedDesc)%3 == 0 {
			itemPrice, _ := strconv.ParseFloat(item.Price, 64)
			itemPoints := int(math.Ceil(itemPrice * 0.2))
			points += itemPoints
			breakdownLines = append(breakdownLines, fmt.Sprintf("     %d Points - \"%s\" is %d characters (a multiple of 3)", itemPoints, trimmedDesc, len(trimmedDesc)))
			breakdownLines = append(breakdownLines, fmt.Sprintf("                item price of %s * 0.2 = %.2f, rounded up is %d points", item.Price, itemPrice*0.2, itemPoints))
		}
	}

	// 6. 6 points if the purchase day is odd
	purchaseDate, _ := time.Parse("2006-01-02", receipt.PurchaseDate)
	if purchaseDate.Day()%2 != 0 {
		points += 6
		breakdownLines = append(breakdownLines, "     6 points - purchase day is odd")
	}

	// 7. 10 points if the purchase time is between 2:00pm and 4:00pm
	purchaseTime, _ := time.Parse("15:04", receipt.PurchaseTime)
	if purchaseTime.Hour() >= 14 && purchaseTime.Hour() < 16 {
		points += 10
		breakdownLines = append(breakdownLines, "    10 points - 2:33pm is between 2:00pm and 4:00pm")
	}

	breakdownLines = append(breakdownLines, "  + ---------")
	breakdownLines = append(breakdownLines, fmt.Sprintf("  = %d points", points))

	return points, strings.Join(breakdownLines, "\n")
}

// To processReceiptHandler which handles the POST /receipts/process request
func processReceiptHandler(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt
	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// To Generate a unique receipt ID
	receiptID := uuid.New().String()

	points, breakdown := calculatePoints(receipt)

	receipts[receiptID] = ReceiptData{
		Receipt:   receipt,
		Points:    points,
		Breakdown: breakdown,
	}

	// Returning response
	response := map[string]string{"id": receiptID}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// getPointsHandler to handle the GET /receipts/{id}/points request
func getPointsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	receiptID := vars["id"]

	receiptData, exists := receipts[receiptID]
	if !exists {
		http.Error(w, "No receipt found for that ID", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"total_points": receiptData.Points,
		"breakdown":    receiptData.Breakdown,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}