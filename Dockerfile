FROM golang:1.22-alpine AS build

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /app/main .

# Install ca-certificates for HTTPS connections
RUN apk --no-cache add ca-certificates

# Create non-root user for security
RUN adduser -D -g '' appuser

# Set the ownership of the app directory
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# The environment variables will be provided at runtime
CMD ["./main"]
