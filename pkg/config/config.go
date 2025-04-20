package config

import (
	"os"
)

// Config struct to hold RabbitMQ connection parameters
type Config struct {
	RabbitMQHost     string
	RabbitMQPort     string
	RabbitMQUser     string
	RabbitMQPassword string
	RabbitMQVHost    string

	// Database configuration
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
}

// LoadConfigFromEnv loads configuration from environment variables
func LoadConfigFromEnv() Config {
	return Config{
		RabbitMQHost:     getEnv("RABBITMQ_HOST", "host.docker.internal"),
		RabbitMQPort:     getEnv("RABBITMQ_PORT", "5672"),
		RabbitMQUser:     getEnv("RABBITMQ_USER", "guest"),
		RabbitMQPassword: getEnv("RABBITMQ_PASSWORD", "guest"),
		RabbitMQVHost:    getEnv("RABBITMQ_VHOST", "/"),

		// Load database configuration
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           getEnv("DB_PORT", "5432"),
		DBUser:           getEnv("DB_USER", "postgres"),
		DBPassword:       getEnv("DB_PASSWORD", "postgres"),
		DBName:           getEnv("DB_NAME", "reportdb"),
	}
}

// getEnv retrieves an environment variable or returns the default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
} 