package main

// To represents the structure of a receipt
type Receipt struct {
Retailer     string `json:"retailer"`
PurchaseDate string `json:"purchaseDate"`
PurchaseTime string `json:"purchaseTime"`
Total        string `json:"total"`
Items        []Item `json:"items"`
}

// To represents an individual item in the receipt
type Item struct {
ShortDescription string `json:"shortDescription"`
Price            string `json:"price"`
}

// ReceiptData holds receipt ID and calculated points
type ReceiptData struct {
Receipt   Receipt `json:"receipt"`
Points    int     `json:"points"`
Breakdown string  `json:"breakdown"`
}

// In-memory storage
var receipts = make(map[string]ReceiptData)