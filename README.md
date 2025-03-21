# Report Export Service

## Requirements

- Go 1.21 or higher

## Getting Started

1. Download dependencies:

```
go mod download
```

2. Generate Swagger documentation:

```
go install github.com/swaggo/swag/cmd/swag@latest
swag init
```

3. Run the service:

```
go run main.go
```

4. Access the API:
   - Swagger documentation: <http://localhost:8080/swagger/index.html>
