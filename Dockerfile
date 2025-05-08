# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/img-generator ./cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/img-generator .
# Copy configs directory
COPY --from=builder /app/configs ./configs
# Copy assets directory if needed
COPY --from=builder /app/assets ./assets
# Copy font file
COPY --from=builder /app/wqy-zenhei.ttf .

# Create a non-root user
RUN adduser -D -g '' appuser
USER appuser

# Expose the port
EXPOSE 8080

# Set environment variables
ENV SERVER_PORT=8080

# Run the application
CMD ["./img-generator"] 