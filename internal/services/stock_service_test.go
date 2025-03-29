package services

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestStockService_GetStockQuote(t *testing.T) {
	// Mock HTTP server for API responses
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return mock CSV data
		w.Header().Set("Content-Type", "text/csv")
		w.WriteHeader(http.StatusOK)

		// Write CSV header and data
		w.Write([]byte("Symbol,Date,Time,Open,High,Low,Close,Volume\n"))
		w.Write([]byte("AAPL.US,2023-01-01,00:00:00,150.0,155.0,145.0,152.5,10000000\n"))
	}))
	defer server.Close()

	// Save current API URL and restore after test
	originalURL := os.Getenv("STOCK_API_URL")
	defer os.Setenv("STOCK_API_URL", originalURL)

	// Set API URL to mock server
	os.Setenv("STOCK_API_URL", server.URL+"?s=%s")

	// Create stock service
	service := NewStockService()

	// Test valid stock quote
	response, err := service.GetStockQuote("AAPL.US")

	// Check for errors
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Check response data
	if response.Symbol != "AAPL.US" {
		t.Errorf("Expected symbol 'AAPL.US', got: %s", response.Symbol)
	}

	if response.Price != 152.5 {
		t.Errorf("Expected price 152.5, got: %f", response.Price)
	}

	if response.Error != "" {
		t.Errorf("Expected no error message, got: %s", response.Error)
	}
}

func TestStockService_GetStockQuote_Error(t *testing.T) {
	// Mock HTTP server for API errors
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return error status
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	// Save current API URL and restore after test
	originalURL := os.Getenv("STOCK_API_URL")
	defer os.Setenv("STOCK_API_URL", originalURL)

	// Set API URL to mock server
	os.Setenv("STOCK_API_URL", server.URL+"?s=%s")

	// Create stock service
	service := NewStockService()

	// Test stock quote with API error
	response, err := service.GetStockQuote("INVALID")

	// Should not return an error but include error in response
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Check response has error message
	if response.Error == "" {
		t.Errorf("Expected error message, got empty string")
	}
}

func TestStockService_GetStockQuote_InvalidCSV(t *testing.T) {
	// Mock HTTP server for invalid CSV
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/csv")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Invalid CSV data"))
	}))
	defer server.Close()

	// Save current API URL and restore after test
	originalURL := os.Getenv("STOCK_API_URL")
	defer os.Setenv("STOCK_API_URL", originalURL)

	// Set API URL to mock server
	os.Setenv("STOCK_API_URL", server.URL+"?s=%s")

	// Create stock service
	service := NewStockService()

	// Test stock quote with invalid CSV
	response, err := service.GetStockQuote("TEST")

	// Should not return an error but include error in response
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Check response has error message
	if response.Error == "" {
		t.Errorf("Expected error message, got empty string")
	}
}
