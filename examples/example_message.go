//go:build ignore
// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/burakmike/report-export-service/pkg/model"
)

func main() {
	// Create sample portfolios
	portfolios := []model.Portfolio{
		{
			PortID:     1,
			Name:       "My Tech Portfolio",
			UserID:     "user123",
			CreatedAt:  "2023-01-01 10:00:00",
			LastUpdate: "2023-06-15 14:30:00",
		},
		{
			PortID:     2,
			Name:       "Retirement Fund",
			UserID:     "user456",
			CreatedAt:  "2023-02-10 09:00:00",
			LastUpdate: "2023-07-20 11:00:00",
		},
		{
			PortID:     3,
			Name:       "Growth Portfolio",
			UserID:     "user123",
			CreatedAt:  "2023-03-05 12:30:00",
			LastUpdate: "2023-08-01 09:15:00",
		},
	}

	// Create a portfolio report message
	message := model.NewPortfolioReportMessage(portfolios)

	// Convert to JSON
	jsonData, err := json.MarshalIndent(message, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal message: %v", err)
	}

	// Print the JSON message
	fmt.Println("Portfolio Report Message JSON:")
	fmt.Println(string(jsonData))

	// Example of how to parse a JSON message back into structs
	var parsedMessage model.PortfolioReportMessage
	err = json.Unmarshal(jsonData, &parsedMessage)
	if err != nil {
		log.Fatalf("Failed to unmarshal message: %v", err)
	}

	// Access the parsed data
	fmt.Printf("\nParsed Message Details:\n")
	fmt.Printf("Event Type: %s\n", parsedMessage.EventType)
	fmt.Printf("Timestamp: %s\n", parsedMessage.Timestamp)
	fmt.Printf("Number of Portfolios: %d\n", len(parsedMessage.Payload.Portfolios))

	// Access individual portfolios
	for i, p := range parsedMessage.Payload.Portfolios {
		fmt.Printf("\nPortfolio %d:\n", i+1)
		fmt.Printf("  ID: %d\n", p.PortID)
		fmt.Printf("  Name: %s\n", p.Name)
		fmt.Printf("  User ID: %s\n", p.UserID)
		fmt.Printf("  Created At: %s\n", p.CreatedAt)
		fmt.Printf("  Last Update: %s\n", p.LastUpdate)
	}
} 