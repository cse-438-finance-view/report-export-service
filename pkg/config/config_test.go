package config

import (
	"os"
	"testing"
)

func TestLoadConfigFromEnv_Defaults(t *testing.T) {
	// Clear relevant environment variables
	keys := []string{
		"RABBITMQ_HOST", "RABBITMQ_PORT", "RABBITMQ_USER", "RABBITMQ_PASSWORD", "RABBITMQ_VHOST",
		"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME",
	}
	for _, k := range keys {
		os.Unsetenv(k)
	}

	cfg := LoadConfigFromEnv()
	// Verify some default values
	if got, want := cfg.RabbitMQHost, "host.docker.internal"; got != want {
		t.Errorf("RabbitMQHost default = %q; want %q", got, want)
	}
	if got, want := cfg.RabbitMQPort, "5672"; got != want {
		t.Errorf("RabbitMQPort default = %q; want %q", got, want)
	}
	if got, want := cfg.DBName, "reportdb"; got != want {
		t.Errorf("DBName default = %q; want %q", got, want)
	}
}

func TestLoadConfigFromEnv_Custom(t *testing.T) {
	// Set custom environment variables
	domain := "testhost"
	dbname := "mytestdb"
	os.Setenv("RABBITMQ_HOST", domain)
	defer os.Unsetenv("RABBITMQ_HOST")
	os.Setenv("DB_NAME", dbname)
	defer os.Unsetenv("DB_NAME")

	cfg := LoadConfigFromEnv()
	if got, want := cfg.RabbitMQHost, domain; got != want {
		t.Errorf("RabbitMQHost = %q; want %q", got, want)
	}
	if got, want := cfg.DBName, dbname; got != want {
		t.Errorf("DBName = %q; want %q", got, want)
	}
} 