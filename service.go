package main

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Function to calculate points for a given receipt
func calculatePoints(receipt Receipt) int {
	points := 0

	// 1. One point for every alphanumeric character in the retailer name.
	regex := regexp.MustCompile("[a-zA-Z0-9]")
	points += len(regex.FindAllString(receipt.Retailer, -1))

	// 2. 50 points if the total is a round dollar amount (e.g., "35.00").
	total, _ := strconv.ParseFloat(receipt.Total, 64)
	if total == float64(int(total)) {
		points += 50
	}

	// 3. 25 points if the total is a multiple of 0.25.
	if math.Mod(total, 0.25) == 0 {
		points += 25
	}

	// 4. 5 points for every two items on the receipt.
	points += (len(receipt.Items) / 2) * 5

	// 5. If the item description length is a multiple of 3, add (price * 0.2) rounded up.
	for _, item := range receipt.Items {
		descLen := len(strings.TrimSpace(item.ShortDescription))
		price, _ := strconv.ParseFloat(item.Price, 64)
		if descLen%3 == 0 {
			points += int(math.Ceil(price * 0.2))
		}
	}

	// 6. 6 points if the day in the purchase date is odd.
	parsedDate, _ := time.Parse("2006-01-02", receipt.PurchaseDate)
	if parsedDate.Day()%2 == 1 {
		points += 6
	}

	// 7. 10 points if purchase time is between 2:00pm and 4:00pm.
	parsedTime, _ := time.Parse("15:04", receipt.PurchaseTime)
	if parsedTime.Hour() >= 14 && parsedTime.Hour() < 16 {
		points += 10
	}

	return points
}
