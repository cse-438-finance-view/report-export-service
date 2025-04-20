# Report Export Microservice

A microservice that connects to RabbitMQ for processing portfolio data, generating PDF reports, and storing report metadata in PostgreSQL.

## Architecture Overview

The Report Export Service is designed with a modular, event-driven architecture that processes portfolio data received via RabbitMQ and generates PDF reports. The architecture follows these key principles:

- **Event-Driven**: Uses an event-based system for processing messages
- **Modular**: Organized in packages with clear responsibilities
- **Extensible**: Easily add new event types and handlers
- **Resilient**: Includes reconnection mechanisms and error handling
- **Persistent**: Stores report metadata in a PostgreSQL database

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
├── examples/        # Example scripts (e.g. test publisher)
└── main.go          # Application entry point
```

## Message Processing Flow

1. **Connection**: Connect to RabbitMQ and PostgreSQL using credentials from environment variables
2. **Exchange/Queue Setup**: Declare exchange and queues with appropriate bindings
3. **Message Consumption**: Listen for messages on configured queues
4. **Event Routing**: Parse incoming messages and route to appropriate handlers
5. **Report Generation**: Process portfolio data and generate PDF reports
6. **Database Write**: Store report metadata in the `reports` table
7. **Acknowledgment**: Acknowledge processed messages to RabbitMQ

### RabbitMQ Integration

The service establishes a connection to RabbitMQ and configures the following:

- **Exchange**: `investment_exchange` (type: topic)
- **Queue**: `portfolio_report_queue` 
- **Binding**: Routes messages with the routing key `portfolio.report` to the queue

The service implements a robust connection handling mechanism:
- Automatic reconnection with exponential backoff
- Connection monitoring
- Graceful shutdown

### PostgreSQL Integration

- Stores each generated report's metadata in a `reports` table with columns: `id`, `created_at`, `user_id`, `type`.
- Table is created automatically if it does not exist.

## Event Processing

### Event Types

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
| `RABBITMQ_HOST` | RabbitMQ host | `host.docker.internal` (for Docker) |
| `RABBITMQ_PORT` | RabbitMQ port | `5672` |
| `RABBITMQ_USER` | RabbitMQ username | `guest` |
| `RABBITMQ_PASSWORD` | RabbitMQ password | `guest` |
| `RABBITMQ_VHOST` | RabbitMQ virtual host | `/` |
| `DB_HOST` | PostgreSQL host | `db` |
| `DB_PORT` | PostgreSQL port | `4450` |
| `DB_USER` | PostgreSQL username | `postgres` |
| `DB_PASSWORD` | PostgreSQL password | `postgres` |
| `DB_NAME` | PostgreSQL database name | `reportdb` |

## Running the Service

### With Docker Compose (Recommended)

1. Create a `.env` file in the project root with your RabbitMQ and PostgreSQL credentials (see example below).
2. Start the services:

```bash
docker-compose up --build
```

This will start the service and a PostgreSQL database. **RabbitMQ must be running separately and accessible to the service.**

### Directly with Go

```bash
go run main.go
```

## Example .env File

```
# RabbitMQ
RABBITMQ_HOST=host.docker.internal
RABBITMQ_PORT=5672
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest
RABBITMQ_VHOST=/

# Postgres
DB_HOST=db
DB_PORT=4450
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=reportdb
```

## Test Utilities

### Publish a Test Report Event

You can publish a sample report event to RabbitMQ using the provided example script:

```bash
# If running on your host, override RABBITMQ_HOST as needed:
RABBITMQ_HOST=localhost make publish-test-report
```

This will send a sample event to RabbitMQ. The service will process it, generate a PDF, and write a record to the database.

### Check the Database for Reports

To see the latest reports in the database:

```bash
docker-compose exec -T db psql -U $DB_USER -d $DB_NAME -c "SELECT * FROM reports ORDER BY created_at DESC LIMIT 10;" | cat
```

### Download Generated PDF Reports

To copy the latest generated PDF from the container to your host:

```bash
LATEST_PDF=$(docker exec report-export-service sh -c 'ls -t /app/reports/portfolio_report_*.pdf | head -n 1')
if [ -z "$LATEST_PDF" ]; then \
  echo "No PDF report found in container."; \
else \
  echo "Copying $LATEST_PDF..."; \
  docker cp report-export-service:"$LATEST_PDF" .; \
  echo "Copied to $(basename $LATEST_PDF)"; \
fi
```

Or to copy all reports:

```bash
mkdir -p downloaded_reports
docker cp report-export-service:/app/reports ./downloaded_reports/
```

## Technical Implementation Details

### Handler Registry

The service uses a registry pattern to map event types to their handlers:

```go
// Register handlers
registry.RegisterHandler(handler.NewPortfolioReportHandler(db, pdfGenerator))
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
