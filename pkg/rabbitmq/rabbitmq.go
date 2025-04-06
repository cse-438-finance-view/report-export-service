package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/burakmike/report-export-service/pkg/config"
	"github.com/burakmike/report-export-service/pkg/event"
	"github.com/burakmike/report-export-service/pkg/handler"
	"github.com/streadway/amqp"
)

// RabbitMQClient RabbitMQ bağlantısını ve kanalını yöneten yapı
type RabbitMQClient struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Config     config.Config
	Registry   *handler.HandlerRegistry
}

// NewRabbitMQClient yeni bir RabbitMQ client oluşturur
func NewRabbitMQClient(cfg config.Config, registry *handler.HandlerRegistry) *RabbitMQClient {
	return &RabbitMQClient{
		Config:   cfg,
		Registry: registry,
	}
}

// Connect RabbitMQ'ya bağlanır
func (r *RabbitMQClient) Connect() error {
	var err error

	// Bağlantı URL'ini oluştur
	connStr := fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		r.Config.RabbitMQUser,
		r.Config.RabbitMQPassword,
		r.Config.RabbitMQHost,
		r.Config.RabbitMQPort,
		r.Config.RabbitMQVHost)

	// RabbitMQ sunucusuna bağlan
	r.Connection, err = amqp.Dial(connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	// Kanal aç
	r.Channel, err = r.Connection.Channel()
	if err != nil {
		r.Connection.Close()
		return fmt.Errorf("failed to open a channel: %w", err)
	}

	log.Println("Successfully connected to RabbitMQ")
	return nil
}

// SetupExchangeAndQueues exchange ve kuyrukları oluşturur
func (r *RabbitMQClient) SetupExchangeAndQueues() error {
	// Exchange oluştur
	err := r.Channel.ExchangeDeclare(
		"investment_exchange", // name
		"topic",               // type
		true,                  // durable
		false,                 // auto-deleted
		false,                 // internal
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Kuyruk tanımlamaları ve bağlamaları
	queueBindings := map[string][]string{
		"portfolio_report_queue": {"portfolio.report"},
		// İleride yeni kuyruklar ve routing key'ler buraya eklenebilir
	}

	// Her bir kuyruğu oluştur ve bağla
	for queueName, routingKeys := range queueBindings {
		// Kuyruk oluştur
		queue, err := r.Channel.QueueDeclare(
			queueName, // name
			true,      // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // no-wait
			nil,       // arguments
		)
		if err != nil {
			return fmt.Errorf("failed to declare queue %s: %w", queueName, err)
		}

		// Routing key'leri bağla
		for _, routingKey := range routingKeys {
			err = r.Channel.QueueBind(
				queue.Name,            // queue name
				routingKey,            // routing key
				"investment_exchange", // exchange
				false,                 // no-wait
				nil,                   // arguments
			)
			if err != nil {
				return fmt.Errorf("failed to bind queue %s with routing key %s: %w", 
					queueName, routingKey, err)
			}
			log.Printf("Queue %s bound to exchange with routing key %s", queueName, routingKey)
		}
	}

	log.Println("Exchange and queues setup complete")
	return nil
}

// ConsumeMessages mesajları dinlemeye başlar ve işler
func (r *RabbitMQClient) ConsumeMessages(ctx context.Context) error {
	// Tüm kuyrukları dinlemek için bir map
	queueConsumers := []string{
		"portfolio_report_queue",
		// İleride yeni kuyruklar eklenebilir
	}

	// Her bir kuyruk için tüketici oluştur
	for _, queueName := range queueConsumers {
		msgs, err := r.Channel.Consume(
			queueName, // queue
			"",        // consumer
			false,     // auto-ack
			false,     // exclusive
			false,     // no-local
			false,     // no-wait
			nil,       // args
		)
		if err != nil {
			return fmt.Errorf("failed to register a consumer for queue %s: %w", queueName, err)
		}

		// Her bir kuyruk için ayrı bir goroutine başlat
		go func(qName string, deliveries <-chan amqp.Delivery) {
			log.Printf("Started consuming messages from queue: %s", qName)
			
			for msg := range deliveries {
				// İşlem başarılı olmazsa mesajı tekrar kuyruğa koy
				if err := r.processMessage(ctx, msg); err != nil {
					log.Printf("Error processing message from queue %s: %v", qName, err)
					msg.Nack(false, true) // mesajı tekrar kuyruğa koy
				} else {
					msg.Ack(false) // başarılı işleme
				}
			}
		}(queueName, msgs)
	}

	log.Println("Message consumers registered for all queues")
	return nil
}

// processMessage gelen mesajı uygun handler'a yönlendirir
func (r *RabbitMQClient) processMessage(ctx context.Context, msg amqp.Delivery) error {
	log.Printf("Received a message: %s", msg.Body)

	// Mesajı BaseEvent yapısına dönüştür
	baseEvent, err := event.ParseEvent(msg.Body)
	if err != nil {
		return fmt.Errorf("failed to parse event: %w", err)
	}

	// Uygun handler'ı bul ve mesajı işle
	return r.Registry.HandleEvent(ctx, baseEvent)
}

// PublishEvent bir event'i RabbitMQ'ya gönderir
func (r *RabbitMQClient) PublishEvent(evt event.BaseEvent, routingKey string) error {
	// Event'i JSON'a dönüştür
	eventData, err := json.Marshal(evt)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Mesajı yayınla
	err = r.Channel.Publish(
		"investment_exchange", // exchange
		routingKey,            // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        eventData,
			Timestamp:   time.Now(),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	log.Printf("Published event with type %s to exchange with routing key %s", evt.EventType, routingKey)
	return nil
}

// PublishPortfolioReport portföy rapor event'ini gönderir
func (r *RabbitMQClient) PublishPortfolioReport(portfolios []event.Portfolio) error {
	// Portföy rapor event'i oluştur
	evt, err := event.NewPortfolioReportEvent(portfolios)
	if err != nil {
		return fmt.Errorf("failed to create portfolio report event: %w", err)
	}

	// Event'i yayınla
	return r.PublishEvent(evt, "portfolio.report")
}

// Close bağlantıyı ve kanalı kapatır
func (r *RabbitMQClient) Close() {
	if r.Channel != nil {
		r.Channel.Close()
	}
	if r.Connection != nil {
		r.Connection.Close()
	}
	log.Println("RabbitMQ connection closed")
}

// Reconnect RabbitMQ'ya yeniden bağlanmayı ve kurulumu tekrarlamayı dener
func (r *RabbitMQClient) Reconnect(ctx context.Context, maxRetries int) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		log.Printf("Attempting to reconnect to RabbitMQ (attempt %d/%d)", i+1, maxRetries)
		err = r.Connect()
		if err == nil {
			err = r.SetupExchangeAndQueues()
			if err == nil {
				err = r.ConsumeMessages(ctx)
				if err == nil {
					log.Println("Successfully reconnected to RabbitMQ")
					return nil
				}
			}
			r.Close()
		}

		// Exponential backoff with max of 30 seconds
		backoff := time.Duration(min(int64(1<<uint(i)), 30)) * time.Second
		log.Printf("Reconnect failed, retrying in %v...", backoff)
		time.Sleep(backoff)
	}
	return fmt.Errorf("failed to reconnect after %d attempts: %w", maxRetries, err)
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
} 