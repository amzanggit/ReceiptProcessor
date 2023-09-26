package ReceiptProcessor

import (
	"testing"
)

func TestProcessReceiptsAndRetrievePoints(t *testing.T) {
	tests := []struct {
		name           string
		receiptJSON    string
		expectedPoints int
		wantStatus     int
	}{
		{
			name: "Example Test case 1",
			receiptJSON: `{
			  	"retailer": "Target",
			  	"purchaseDate": "2022-01-01",
			  	"purchaseTime": "13:01",
			  	"items": [
			    	{
			      		"shortDescription": "Mountain Dew 12PK",
			      		"price": "6.49"
			    	},{
			      		"shortDescription": "Emils Cheese Pizza",
			      		"price": "12.25"
			    	},{
			      		"shortDescription": "Knorr Creamy Chicken",
			      		"price": "1.26"
			    	},{
			      		"shortDescription": "Doritos Nacho Cheese",
			      		"price": "3.35"
			    	},{
			      		"shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
			      		"price": "12.00"
			    	}
			  	],
			  	"total": "35.35"
			}`,
			expectedPoints: 28,
		}, {
			name: "Example Test case 2",
			receiptJSON: `{
			  	"retailer": "M&M Corner Market",
			 	"purchaseDate": "2022-03-20",
			  	"purchaseTime": "14:33",
			  	"items": [
			    	{
				      	"shortDescription": "Gatorade",
				      	"price": "2.25"
				    },{
				      	"shortDescription": "Gatorade",
				      	"price": "2.25"
				    },{
				      	"shortDescription": "Gatorade",
				      	"price": "2.25"
				    },{
				      	"shortDescription": "Gatorade",
				      	"price": "2.25"
				    }
			  	],
			  	"total": "9.00"
			}`,
			expectedPoints: 109,
		}, {
			name: "Test time between 14 and 16",
			receiptJSON: `{
			  	"retailer": "Target",
			  	"purchaseDate": "2022-01-01",
			  	"purchaseTime": "14:00",
			  	"items": [
			    	{
			      		"shortDescription": "Mountain Dew 12PK  ",
			      		"price": "6.49"
			    	},{
			      		"shortDescription": " Emils Cheese Pizza ",
			      		"price": "12.25"
			    	},{
			      		"shortDescription": "   Knorr Creamy Chicken ",
			      		"price": "4.73"
			    	}
			  	],
			  	"total": "23.47"
			}`,
			expectedPoints: 30,
		}, {
			name: "Test invaild date",
			receiptJSON: `{
			  	"retailer": "Target",
			  	"purchaseDate": "2022-14-01",
			  	"purchaseTime": "22:01",
			  	"items": [
			    	{
			      		"shortDescription": "Mountain Dew 12PK  ",
			      		"price": "6.49"
			    	},{
			      		"shortDescription": " Emils Cheese Pizza ",
			      		"price": "12.25"
			    	},{
			      		"shortDescription": "   Knorr Creamy Chicken ",
			      		"price": "1.26"
			    	}
			  	],
			  	"total": "20.00"
			}`,
			// expectedPoints: 89,
			wantStatus: 
		}, {
			name: "Test invaild date/time 2",
			receiptJSON: `{
			  	"retailer": "Target",
			  	"purchaseDate": "2022/03/20",
			  	"purchaseTime": "13",
			  	"items": [
			    	{
			      		"shortDescription": "Mountain Dew 12PK  ",
			      		"price": "6.50"
			    	},{
			      		"shortDescription": " Emils Cheese Pizza ",
			      		"price": "12.25"
			    	},{
			      		"shortDescription": "   Knorr Creamy Chicken ",
			      		"price": "15.26"
			    	}
			  	],
			  	"total": "34.01"
			}`,
			expectedPoints: 14,
		},
		// Add more test cases here

	}
}

func TestProcessReceipts_InvalidJSON(t *testing.T) {
	invalidJSON := `{ "invalid": "json" }`
	resp, err := http.Post("http://localhost:8080/receipts/process", "application/json", strings.NewReader(invalidJSON))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Check the status code should be 400 Bad Request
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d, but got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestGetPoints_InvalidID(t *testing.T) {
	// Create a request with invalid ID
	invalidId := "id" // Corrected the variable name
	resp, err := http.Get("http://localhost:8080/receipts/" + invalidId + "/points")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Check the status code should be 404 Not Found
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, but got %d", http.StatusNotFound, resp.StatusCode)
	}
}
