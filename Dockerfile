FROM golang:1.22-alpine AS build

WORKDIR /app

# Install Swagger
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate Swagger docs
RUN /go/bin/swag init

# Build the application
RUN go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /app/main .
# Copy the Swagger docs
COPY --from=build /app/docs ./docs

EXPOSE 4444

CMD ["./main"]
