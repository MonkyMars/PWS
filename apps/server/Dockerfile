# --- Stage 1: Build the Go binary ---
FROM golang:1.25-alpine AS builder

# Install required packages
RUN apk add --no-cache git

WORKDIR /app

# Cache modules first
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./main.go

# --- Stage 2: Minimal runtime image ---
FROM alpine:3.20

WORKDIR /app

# Copy binary
COPY --from=builder /app/server .

# Copy env file (or mount with docker-compose)
COPY .env .env

# Fiber runs on port 8082
EXPOSE 8082

# Run
CMD ["./server"]
