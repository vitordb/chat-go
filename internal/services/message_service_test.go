package services

import (
	"regexp"
	"testing"
)

// Create a regex pattern just for testing
var testStockCommandPattern = regexp.MustCompile(`^/stock=([A-Za-z0-9.]+)$`)

func TestStockCommandRegex(t *testing.T) {
	testCases := []struct {
		name      string
		input     string
		isMatch   bool
		stockCode string
	}{
		{
			name:      "Valid stock command",
			input:     "/stock=AAPL.US",
			isMatch:   true,
			stockCode: "AAPL.US",
		},
		{
			name:      "Valid stock command with lowercase",
			input:     "/stock=aapl.us",
			isMatch:   true,
			stockCode: "aapl.us",
		},
		{
			name:      "Valid stock command with numbers",
			input:     "/stock=AAPL123",
			isMatch:   true,
			stockCode: "AAPL123",
		},
		{
			name:      "Valid stock command with dots",
			input:     "/stock=AAPL.US.XYZ",
			isMatch:   true,
			stockCode: "AAPL.US.XYZ",
		},
		{
			name:      "Invalid stock command - missing equals",
			input:     "/stockAAPL.US",
			isMatch:   false,
			stockCode: "",
		},
		{
			name:      "Invalid stock command - extra space",
			input:     "/stock= AAPL.US",
			isMatch:   false,
			stockCode: "",
		},
		{
			name:      "Invalid stock command - extra text",
			input:     "/stock=AAPL.US extra",
			isMatch:   false,
			stockCode: "",
		},
		{
			name:      "Regular message",
			input:     "This is a regular message",
			isMatch:   false,
			stockCode: "",
		},
		{
			name:      "Empty message",
			input:     "",
			isMatch:   false,
			stockCode: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			match := testStockCommandPattern.FindStringSubmatch(tc.input)

			isMatch := match != nil
			if isMatch != tc.isMatch {
				t.Errorf("Expected isMatch=%v, got %v", tc.isMatch, isMatch)
			}

			if isMatch && match[1] != tc.stockCode {
				t.Errorf("Expected stockCode=%s, got %s", tc.stockCode, match[1])
			}
		})
	}
}
