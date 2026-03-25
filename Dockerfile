# Build stage
FROM golang:1.26.1-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0 for a static binary (required for modernc.org/sqlite)
RUN CGO_ENABLED=0 GOOS=linux go build -o aura-pos-api ./cmd/api/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/aura-pos-api .
# Copy migrations if they are needed for runtime (yes, they are)
COPY --from=builder /app/tenant/migrations ./tenant/migrations

# Expose the API port
EXPOSE 8081

# Command to run the application
CMD ["./aura-pos-api"]
