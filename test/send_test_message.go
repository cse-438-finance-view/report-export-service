//go:build ignore
// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/burakmike/report-export-service/pkg/event"
	"github.com/streadway/amqp"
)

// Portfolio represents a single portfolio
type Portfolio struct {
	PortID     int    `json:"portID"`
	Name       string `json:"name"`
	UserID     string `json:"userID"`
	CreatedAt  string `json:"createdAt"`
	LastUpdate string `json:"lastUpdate"`
}

// PortfolioPayload represents the payload part of the message
type PortfolioPayload struct {
	Portfolios []Portfolio `json:"portfolios"`
}

// PortfolioReportMessage represents the complete message structure
type PortfolioReportMessage struct {
	EventType string          `json:"event_type"`
	Timestamp string          `json:"timestamp"`
	Payload   PortfolioPayload `json:"payload"`
}

func main() {
	// RabbitMQ connection parameters
	host := getEnv("RABBITMQ_HOST", "localhost")
	port := getEnv("RABBITMQ_PORT", "5672")
	user := getEnv("RABBITMQ_USER", "guest")
	password := getEnv("RABBITMQ_PASSWORD", "guest")
	vhost := getEnv("RABBITMQ_VHOST", "/")

	// Create the connection string
	connStr := fmt.Sprintf("amqp://%s:%s@%s:%s%s", user, password, host, port, vhost)
	
	fmt.Printf("Connecting to RabbitMQ at: %s\n", host)

	// Connect to RabbitMQ
	conn, err := amqp.Dial(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	
	fmt.Println("Successfully connected to RabbitMQ")

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare the exchange
	err = ch.ExchangeDeclare(
		"investment_exchange", // name
		"topic",               // type
		true,                  // durable
		false,                 // auto-deleted
		false,                 // internal
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare exchange: %v", err)
	}
	
	fmt.Println("Exchange declared successfully")

	// Create sample test message
	portfolios := []event.Portfolio{
		{
			PortID:     1,
			Name:       "Teknoloji Portföyüm",
			UserID:     "user123",
			CreatedAt:  time.Now().AddDate(0, -6, 0).Format("2006-01-02 15:04:05"),
			LastUpdate: time.Now().AddDate(0, 0, -10).Format("2006-01-02 15:04:05"),
		},
		{
			PortID:     2,
			Name:       "Emeklilik Fonu",
			UserID:     "user456",
			CreatedAt:  time.Now().AddDate(0, -4, 0).Format("2006-01-02 15:04:05"),
			LastUpdate: time.Now().AddDate(0, 0, -5).Format("2006-01-02 15:04:05"),
		},
		{
			PortID:     3,
			Name:       "Büyüme Portföyü",
			UserID:     "user123",
			CreatedAt:  time.Now().AddDate(0, -2, 0).Format("2006-01-02 15:04:05"),
			LastUpdate: time.Now().AddDate(0, 0, -1).Format("2006-01-02 15:04:05"),
		},
	}

	// Event oluştur
	evt, err := event.NewPortfolioReportEvent(portfolios)
	if err != nil {
		log.Fatalf("Olay oluşturulamadı: %v", err)
	}

	// JSON'a dönüştür
	jsonData, err := json.MarshalIndent(evt, "", "  ")
	if err != nil {
		log.Fatalf("JSON dönüşümü başarısız: %v", err)
	}

	// Publish the message
	err = ch.Publish(
		"investment_exchange", // exchange
		"portfolio.report",    // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonData,
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	}

	log.Println("Test message sent successfully!")
	fmt.Println("Message content:")
	fmt.Println(string(jsonData))
}

// getEnv retrieves an environment variable or returns the default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
} 