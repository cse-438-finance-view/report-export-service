.PHONY: build run docker-build docker-run clean deps test help publish-test-report

# Build the application
build:
	go build -o bin/report-export-service .

# Run the application
run:
	go run main.go

# Build Docker image
docker-build:
	docker build -t report-export-service .

# Run Docker container
docker-run:
	docker run -d --name report-export-service \
	-e RABBITMQ_HOST=localhost \
	-e RABBITMQ_PORT=5672 \
	-e RABBITMQ_USER=guest \
	-e RABBITMQ_PASSWORD=guest \
	report-export-service

# Clean build artifacts
clean:
	rm -rf bin/

# Get dependencies
deps:
	go mod tidy

# Test the application
test:
	go test -v ./...

# Publish a sample report event using the example script
# Requires .env file to be present for RabbitMQ connection details
publish-test-report:
	@echo "Publishing sample portfolio report event..."
	go run examples/publish_sample_report.go

# Help message
help:
	@echo "Available commands:"
	@echo "  make build         - Build the application"
	@echo "  make run           - Run the application"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-run    - Run Docker container"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make deps          - Get dependencies"
	@echo "  make test          - Run tests"
	@echo "  make publish-test-report - Publish a sample portfolio report event to RabbitMQ"
	@echo "  make help          - Show this help message"

default: help 