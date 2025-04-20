package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/burakmike/report-export-service/pkg/config"
	"github.com/burakmike/report-export-service/pkg/event"
	"github.com/streadway/amqp"
)

func main() {
	log.Println("Publisher starting...")
	cfg := config.LoadConfigFromEnv()

	if hostOverride := os.Getenv("RABBITMQ_HOST"); hostOverride != "" {
		log.Printf("Overriding RabbitMQ host from env: %s", hostOverride)
		cfg.RabbitMQHost = hostOverride
	}

	connStr := fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		cfg.RabbitMQUser,
		cfg.RabbitMQPassword,
		cfg.RabbitMQHost,
		cfg.RabbitMQPort,
		cfg.RabbitMQVHost)

	log.Printf("Connecting to RabbitMQ at %s:%s...", cfg.RabbitMQHost, cfg.RabbitMQPort)

	conn, err := amqp.Dial(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	log.Println("RabbitMQ connection successful.")

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()
	log.Println("RabbitMQ channel opened.")

	portfolios := event.CreateSamplePortfolios()
	log.Printf("Created %d sample portfolios.", len(portfolios))

	reportEvent, err := event.NewPortfolioReportEvent(portfolios)
	if err != nil {
		log.Fatalf("Failed to create report event: %v", err)
	}

	eventData, err := json.Marshal(reportEvent)
	if err != nil {
		log.Fatalf("Failed to marshal event: %v", err)
	}

	routingKey := string(event.PortfolioReport)
	exchangeName := "investment_exchange"
	log.Printf("Publishing event to exchange '%s' with routing key '%s'...", exchangeName, routingKey)

	err = ch.Publish(
		exchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        eventData,
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	}

	log.Printf("Successfully published sample '%s' event.", reportEvent.EventType)
}