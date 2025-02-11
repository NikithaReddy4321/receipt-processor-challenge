package main

import (
"fmt"
"log"
"net/http"

"github.com/gorilla/mux"
)

func main() {
r := mux.NewRouter()

r.HandleFunc("/receipts/process", processReceiptHandler).Methods("POST")
r.HandleFunc("/receipts/{id}/points", getPointsHandler).Methods("GET")

port := "8080"
fmt.Printf("Server is running on port %s...\n", port)
log.Fatal(http.ListenAndServe(":"+port, r))
}