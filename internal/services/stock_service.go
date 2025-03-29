package services

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/dbvitor/chat-go/internal/models"
)

// StockService handles stock-related business logic
type StockService struct {
	apiURL string
}

// NewStockService creates a new stock service
func NewStockService() *StockService {
	return &StockService{
		apiURL: os.Getenv("STOCK_API_URL"),
	}
}

// GetStockQuote retrieves a stock quote from the API
func (s *StockService) GetStockQuote(stockCode string) (*models.StockResponse, error) {
	// Create API URL
	url := fmt.Sprintf(s.apiURL, stockCode)

	// Make HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return &models.StockResponse{
			Symbol: stockCode,
			Error:  "failed to connect to stock API",
		}, nil
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return &models.StockResponse{
			Symbol: stockCode,
			Error:  fmt.Sprintf("API returned status code %d", resp.StatusCode),
		}, nil
	}

	// Parse CSV response
	reader := csv.NewReader(resp.Body)

	// Read header
	header, err := reader.Read()
	if err != nil {
		return &models.StockResponse{
			Symbol: stockCode,
			Error:  "failed to parse CSV header",
		}, nil
	}

	// Find close price index
	closeIndex := -1
	for i, column := range header {
		if strings.ToLower(column) == "close" {
			closeIndex = i
			break
		}
	}

	if closeIndex == -1 {
		return &models.StockResponse{
			Symbol: stockCode,
			Error:  "close price not found in CSV",
		}, nil
	}

	// Read data row
	row, err := reader.Read()
	if err != nil {
		if err == io.EOF {
			return &models.StockResponse{
				Symbol: stockCode,
				Error:  "no data found for stock code",
			}, nil
		}
		return &models.StockResponse{
			Symbol: stockCode,
			Error:  "failed to parse CSV data",
		}, nil
	}

	// Parse close price
	closePrice, err := strconv.ParseFloat(row[closeIndex], 64)
	if err != nil {
		return &models.StockResponse{
			Symbol: stockCode,
			Error:  "invalid close price",
		}, nil
	}

	// Check if close price is valid
	if closePrice == 0 {
		return &models.StockResponse{
			Symbol: stockCode,
			Error:  "stock not found or invalid code",
		}, nil
	}

	// Return stock response
	return &models.StockResponse{
		Symbol: stockCode,
		Price:  closePrice,
	}, nil
}
