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

// Service for stock quote queries
type StockService struct {
	apiURL string
}

// Creates the stock quote service
func NewStockService() *StockService {
	return &StockService{
		apiURL: os.Getenv("STOCK_API_URL"),
	}
}

// Fetches the current quote from stooq.com service
func (s *StockService) GetStockQuote(stockCode string) (*models.StockResponse, error) {
	// Builds the URL with the stock code
	url := fmt.Sprintf(s.apiURL, stockCode)

	// Makes the HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return &models.StockResponse{
			Symbol: stockCode,
			Error:  "failed to connect to stock API",
		}, nil
	}
	defer resp.Body.Close()

	// Verifies if the response was OK
	if resp.StatusCode != http.StatusOK {
		return &models.StockResponse{
			Symbol: stockCode,
			Error:  fmt.Sprintf("API returned status code %d", resp.StatusCode),
		}, nil
	}

	// Response is in CSV format, so we use the parser
	reader := csv.NewReader(resp.Body)

	// Reads the header to identify the price column
	header, err := reader.Read()
	if err != nil {
		return &models.StockResponse{
			Symbol: stockCode,
			Error:  "failed to parse CSV header",
		}, nil
	}

	// Looks for the "Close" column which has the closing price
	closeIndex := -1
	for i, column := range header {
		if strings.ToLower(column) == "close" {
			closeIndex = i
			break
		}
	}

	// If column not found, returns error
	if closeIndex == -1 {
		return &models.StockResponse{
			Symbol: stockCode,
			Error:  "close price not found in CSV",
		}, nil
	}

	// Reads the data line (first line after header)
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

	// Converts price from string to float
	closePrice, err := strconv.ParseFloat(row[closeIndex], 64)
	if err != nil {
		return &models.StockResponse{
			Symbol: stockCode,
			Error:  "invalid close price",
		}, nil
	}

	// Some APIs return 0 when the code is invalid
	if closePrice == 0 {
		return &models.StockResponse{
			Symbol: stockCode,
			Error:  "stock not found or invalid code",
		}, nil
	}

	// Returns the result with success
	return &models.StockResponse{
		Symbol: stockCode,
		Price:  closePrice,
	}, nil
}
