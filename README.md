# Report Export Microservice

A microservice that connects to RabbitMQ for processing portfolio data and generating PDF reports.

## Architecture Overview

The Report Export Service is designed with a modular, event-driven architecture that processes portfolio data received via RabbitMQ and generates PDF reports. The architecture follows these key principles:

- **Event-Driven**: Uses an event-based system for processing messages
- **Modular**: Organized in packages with clear responsibilities
- **Extensible**: Easily add new event types and handlers
- **Resilient**: Includes reconnection mechanisms and error handling

### Project Structure

```
report-export-service/
├── pkg/
│   ├── config/      # Configuration management
│   ├── event/       # Event definitions and payload structures
│   ├── handler/     # Event handlers for processing messages
│   ├── rabbitmq/    # RabbitMQ connection and messaging
│   ├── report/      # PDF report generation
│   └── service/     # Main service coordination
├── test/            # Test utilities
├── reports/         # Generated PDF reports
└── main.go          # Application entry point
```

## Message Processing Flow

1. **Connection**: Connect to RabbitMQ server using credentials from environment variables
2. **Exchange/Queue Setup**: Declare exchange and queues with appropriate bindings
3. **Message Consumption**: Listen for messages on configured queues
4. **Event Routing**: Parse incoming messages and route to appropriate handlers
5. **Report Generation**: Process portfolio data and generate PDF reports
6. **Acknowledgment**: Acknowledge processed messages to RabbitMQ

### RabbitMQ Integration

The service establishes a connection to RabbitMQ and configures the following:

- **Exchange**: `investment_exchange` (type: topic)
- **Queue**: `portfolio_report_queue` 
- **Binding**: Routes messages with the routing key `portfolio.report` to the queue

The service implements a robust connection handling mechanism:
- Automatic reconnection with exponential backoff
- Connection monitoring
- Graceful shutdown

## Event Processing

### Event Types

The service processes the following event types:

| Event Type | Description | Routing Key |
|------------|-------------|-------------|
| `portfolio.report` | Request to generate portfolio reports | `portfolio.report` |

### Message Format

Messages follow this JSON structure:

```json
{
  "event_type": "portfolio.report",
  "timestamp": "2023-08-10T12:00:00Z",
  "payload": {
    "portfolios": [
      {
        "portID": 1,
        "name": "My Tech Portfolio",
        "userID": "user123",
        "createdAt": "2023-01-01 10:00:00",
        "lastUpdate": "2023-06-15 14:30:00"
      }
      // Additional portfolios...
    ]
  }
}
```

## Aggregation and PDF Report Generation

The service aggregates portfolio data from incoming messages and generates professional PDF reports:

### PDF Report Features

- **Layout**: Landscape A4 format with proper margins and page numbering
- **Content Structure**:
  - Professional header with title and generation date
  - Data table showing portfolio information
  - Summary section with totals
  - Footer with generation timestamp and copyright information
- **Styling**:
  - Colored table headers
  - Alternating row colors for better readability
  - Proper typography with font variations

### PDF Report Content

Each report contains the following portfolio information:
- Portfolio ID
- Portfolio Name
- User ID
- Created At
- Last Update

## Configuration

The service can be configured using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `RABBITMQ_HOST` | RabbitMQ host | `localhost` |
| `RABBITMQ_PORT` | RabbitMQ port | `5672` |
| `RABBITMQ_USER` | RabbitMQ username | `guest` |
| `RABBITMQ_PASSWORD` | RabbitMQ password | `guest` |
| `RABBITMQ_VHOST` | RabbitMQ virtual host | `/` |

## Running the Service

### Directly with Go

```bash
go run main.go
```

### With Docker

Build the Docker image:

```bash
docker build -t report-export-service .
```

Run the container:

```bash
docker run -d --name report-export-service \
  -e RABBITMQ_HOST=rabbitmq-host \
  -e RABBITMQ_PORT=5672 \
  -e RABBITMQ_USER=username \
  -e RABBITMQ_PASSWORD=password \
  -v $(pwd)/reports:/app/reports \
  report-export-service
```

## Test Utilities

The project includes test utilities to verify functionality:

### Send Test Message

Simulates sending a portfolio report event to RabbitMQ:

```bash
./bin/send_test_message
```

### Generate Test Report

Directly generates a PDF report with sample data:

```bash
./bin/generate_report [output_directory]
```

## Technical Implementation Details

### Handler Registry

The service uses a registry pattern to map event types to their handlers:

```go
// Register handlers
registry.RegisterHandler(handler.NewPortfolioReportHandler())
```

### Connection Resilience

The service implements connection monitoring and automatic reconnection:

```go
// Monitor connection for closure
closeChan := make(chan *amqp.Error)
connection.NotifyClose(closeChan)

// Reconnect with exponential backoff
reconnect(maxRetries)
```

### PDF Report Generation

PDF reports are generated using the `gofpdf` library with a component-based approach:

```go
// Generate report
filePath, err := pdfGenerator.GeneratePortfolioReport(portfolios)
```

## Future Enhancements

Potential future enhancements include:
- Additional report formats (Excel, CSV)
- Email delivery of generated reports
- More sophisticated report templates
- Enhanced data visualization
- Additional event types for different report aggregations
