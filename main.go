package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Item struct {
	ShortDescription string
	Price            string
}

type Receipt struct {
	Retailer     string
	PurchaseDate string
	PurchaseTime string
	Items        []Item
	Total        string
}

var pointsMap sync.Map

func ProcessReceipts(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt
	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the date
	_, errDate := time.Parse("2006-01-02", receipt.PurchaseDate)
	if errDate != nil {
		_, errDateSlash := time.Parse("2006/01/02", receipt.PurchaseDate)
		if errDateSlash != nil {
			http.Error(w, "Invalid date", http.StatusBadRequest)
			return
		}
	}

	// Validate the time
	_, errTime := time.Parse("15:04", receipt.PurchaseTime)
	if errTime != nil {
		http.Error(w, "Invalid time", http.StatusBadRequest)
		return
	}

	// Calculate points based on the rules
	points := CalculatePoints(&receipt)

	// Generate an ID for the receipt
	receiptID := uuid.New().String()

	pointsMap.Store(receiptID, points)

	// Return the ID as JSON response
	response := map[string]string{"id": receiptID}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetPoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	receiptID := vars["id"]
	points, isExist := pointsMap.Load(receiptID)
	if !isExist {
		http.Error(w, "Receipt ID not found", http.StatusNotFound)
		return
	}

	// Return the points as JSON response
	response := map[string]int{"points": points.(int)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func CalculatePoints(receipt *Receipt) int {
	points := 0

	// Rule 1: One point for every alphanumeric character in the retailer name.
	re := regexp.MustCompile("[^a-zA-Z0-9]+")
	points += len(re.ReplaceAllString(receipt.Retailer, ""))

	// Rule 2: 50 points if the total is a round dollar amount with no cents.
	totalFloat, _ := strconv.ParseFloat(receipt.Total, 64)
	fractionalPart := totalFloat - float64(int(totalFloat))
	if fractionalPart == 0.0 {
		points += 50
	}

	// Rule 3: 25 points if the total is a multiple of 0.25.
	if math.Mod(totalFloat, 0.25) == 0 {
		points += 25
	}

	// Rule 4: 5 points for every two items on the receipt.
	points += (len(receipt.Items) / 2) * 5

	// Rule 5: Points based on trimmed item description length
	for _, item := range receipt.Items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			priceFloat, _ := strconv.ParseFloat(item.Price, 64)
			points += int(math.Ceil(priceFloat * 0.2))
		}
	}

	// Rule 6: 6 points if the day in the purchase date is odd.
	purchaseDate, err := time.Parse("2006-01-02", receipt.PurchaseDate)
	if err == nil && purchaseDate.Day()%2 == 1 {
		points += 6
	}

	// Rule 7: 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	purchaseTime, err := time.Parse("15:04", receipt.PurchaseTime)
	start, _ := time.Parse("15:04", "13:59")
	end, _ := time.Parse("15:04", "16:01")
	if err == nil && purchaseTime.After(start) && purchaseTime.Before(end) {
		points += 10
	}

	return points
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/receipts/process", ProcessReceipts).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", GetPoints).Methods("GET")

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
