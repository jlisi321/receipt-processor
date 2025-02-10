package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type Receipt struct {
	ID     string `json:"id"`
	Points int    `json:"points"`
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type IncomingReceipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

var receiptStore = make(map[string]Receipt)

// Error handling function
func sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func GetPoints(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	receiptID := params["id"]

	// Find the points related to the receipt ID in the receiptStore, provide error response if not found
	receipt, exists := receiptStore[receiptID]
	if !exists {
		sendErrorResponse(w, http.StatusNotFound, "No receipt found for that ID.")
		return
	}

	// Provide back a response with the points related to the receipt ID
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]int{"points": receipt.Points})
	if err != nil {
		sendErrorResponse(w, http.StatusNotFound, "No receipt found for that ID.")
		return
	}
}

func pointsForRetailer(receipt IncomingReceipt) int {
	// return count of alphanumeric characters
	return len(strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}
		return -1
	}, receipt.Retailer))
}

func pointsForTotal(receipt IncomingReceipt) int {
	totalAmount, err := strconv.ParseFloat(receipt.Total, 64)
	if err != nil {
		return 0
	}

	points := 0
	totalAmountInCents := int(totalAmount * 100)

	// 50 points if the total is a round dollar amount with no cents.
	if totalAmountInCents%100 == 0 {
		points += 50
	}

	// 25 points if the total is a multiple of 0.25
	if totalAmountInCents%25 == 0 {
		points += 25
	}

	return points
}

func pointsForItemCountAndDescription(receipt IncomingReceipt) int {
	points := 0
	// 5 points for every two items on the receipt.
	points += (len(receipt.Items) / 2) * 5

	// If the trimmed length of the item description is a multiple of 3, calculate points.
	for _, item := range receipt.Items {
		trimmedDescription := strings.TrimSpace(item.ShortDescription)
		if len(trimmedDescription)%3 == 0 {
			itemPrice, err := strconv.ParseFloat(item.Price, 64)
			if err == nil {
				points += int(math.Ceil(itemPrice * 0.2))
			}
		}
	}
	return points
}

func pointsForDate(receipt IncomingReceipt) int {
	date, err := time.Parse("2006-01-02", receipt.PurchaseDate)

	// 6 points if the date is odd
	if err == nil && date.Day()%2 != 0 {
		return 6
	}
	return 0
}

func pointsForTime(receipt IncomingReceipt) int {
	purchaseTime, err := time.Parse("15:04", receipt.PurchaseTime)

	// 10 points if the purchase is between 2:00pm and before 4:00pm non-inclusive
	if err == nil && purchaseTime.Hour() > 14 && purchaseTime.Hour() < 16 {
		return 10
	}
	return 0
}

func CalculatePoints(receipt IncomingReceipt) int {
	points := 0
	points += pointsForRetailer(receipt)
	points += pointsForTotal(receipt)
	points += pointsForItemCountAndDescription(receipt)
	points += pointsForDate(receipt)
	points += pointsForTime(receipt)
	return points
}

func validateReceipt(receipt IncomingReceipt) bool {
	// Validate retailer
	regExpRetailer := regexp.MustCompile("^[\\w\\s\\-&]+$")
	if !regExpRetailer.MatchString(receipt.Retailer) {
		return false
	}

	// Validate purchase date
	_, err := time.Parse("2006-01-02", receipt.PurchaseDate)
	if err != nil {
		return false
	}

	// Validate purchase time
	_, err = time.Parse("15:04", receipt.PurchaseTime)
	if err != nil {
		return false
	}

	// Validate items (must have at least 1 item)
	if len(receipt.Items) == 0 {
		return false
	}

	// Validate each item
	for _, item := range receipt.Items {
		regExpItemDesc := regexp.MustCompile("^[\\w\\s\\-]+$")
		if !regExpItemDesc.MatchString(item.ShortDescription) {
			return false
		}

		regExpPrice := regexp.MustCompile("^\\d+\\.\\d{2}$")
		if !regExpPrice.MatchString(item.Price) {
			return false
		}
	}

	// Validate total amount
	regExpTotal := regexp.MustCompile("^\\d+\\.\\d{2}$")
	if !regExpTotal.MatchString(receipt.Total) {
		return false
	}

	return true
}

func ProcessReceipts(w http.ResponseWriter, r *http.Request) {
	var incomingReceipt IncomingReceipt

	// Decode the incoming JSON request body into a Receipt struct
	if err := json.NewDecoder(r.Body).Decode(&incomingReceipt); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "The receipt is invalid.")
		return
	}

	// Validate the incoming receipt
	if !validateReceipt(incomingReceipt) {
		sendErrorResponse(w, http.StatusBadRequest, "The receipt is invalid.")
		return
	}

	// Provide unique ID for the stored receipt
	newID := uuid.New().String()

	receipt := Receipt{
		ID:     newID,
		Points: CalculatePoints(incomingReceipt),
	}
	receiptStore[newID] = receipt

	// Provide back a response with the unique ID created for the receipt
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(map[string]string{"status": "success", "id": newID})
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "The receipt is invalid.")
		return
	}
}

func main() {
	// Create router
	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/receipts/{id}/points", GetPoints).Methods("GET")
	r.HandleFunc("/receipts/process", ProcessReceipts).Methods("POST")

	// Start server
	fmt.Println("API is running on http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
