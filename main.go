package ReceiptProcessor

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
	"sync"
	"time"
)

type Item struct {
	ShortDescription string  `json:"shortDescription"`
	Price            float64 `json:"price,string"`
}

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

// var pointsMap sync.Map
var receiptPoint sync.Map

func ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the date
	_, errDate := time.Parse("2006-01-02", receipt.PurchaseDate)
	if errDate != nil {
		http.Error(w, "Invalid date", http.StatusBadRequest)
		return
	}

	// Validate the time
	_, errTime := time.Parse("15:04", receipt.PurchaseTime)
	if errTime != nil {
		http.Error(w, "Invalid time", http.StatusBadRequest)
		return
	}

	// generate unique ID for receipt
	receiptID := uuid.New().String()

	// calculate reward points
	points := CalculatePoints(&receipt)

	// store receipt's ID and points into map
	receiptPoint.Store(receiptID, points)

	// locate and return the response by ID and encode in json format
	response := map[string]string{"ID": receiptID}
	// w.Header().set("ContentType", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetPoints(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	receipt, exists := receipts[id]
	if !exists {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	points := CalculatePoints(receipt)
	json.NewEncoder(w).Encode(receiptPoint{"points": points})
}

func CalculatePoints(receipt *Receipt) int {
	points := 0

	// One point for every alphanumeric character in the retailer name
	re := regexp.MustCompile("[^a-zA-Z0-9]+")
	points += len(re.ReplaceAllString(receipt.Retailer, ""))

	// 50 points if the total is a round dollar amount with no cents
	total, _ := strconv.ParseFloat(receipt.Total, 64)
	fracPart := total - float64(int(total))
	if fracPart == 0.0 {
		points += 50
	}

	// 25 points if the total is a multiple of 0.25
	if math.Mod(total, 0.25) == 0 {
		points += 25
	}

	// 5 points for every two items on the receipt
	points += 5 * (len(receipt.Items) / 2)

	//
	for _, item := range receipt.Items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			priceFloat, _ := strconv.ParseFloat(item.Price, 64)
			points += int(math.Ceil(item.Price * 0.2))
		}
	}

	purchaseDate, err := time.Parse("2006-01-02", receipt.PurchaseDate)
	if err == nil && purchaseDate.Day()%2 == 1 {
		points += 6
	}

	purchaseTime, err = time.Parse("15:04", receipt.PurchaseTime)
	if err == nil && purchaseTime.Hour() >= 14 && purchaseTime.Hour() < 16 {
		points += 10
	}

	return points
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/receipts/process", ProcessReceipt).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", GetPoints).Methods("GET")

	http.Handle("/", r)
	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

//func isPurchaseDateOdd(dateStr string) (bool, error) {
//	purchaseDate, err := time.Parse("2006-01-02", dateStr)
//	if err != nil {
//		// Return an error indicating there was an issue with parsing the date
//		return false, fmt.Errorf("failed to parse the date: %w", err)
//	}
//	return purchaseDate.Day()%2 == 1, nil
//}
//
//func isPurchaseTimeBetweenTwoAndFour(timeStr string) (bool, error) {
//	purchaseTime, err := time.Parse("15:04", timeStr)
//	if err != nil {
//		// Return an error indicating there was an issue with parsing the time
//		return false, fmt.Errorf("failed to parse the time: %w", err)
//	}
//
//	start, _ := time.Parse("15:04", "13:59")
//	end, _ := time.Parse("15:04", "16:01")
//
//	return purchaseTime.After(start) && purchaseTime.Before(end), nil
//}
