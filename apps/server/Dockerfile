# --- Stage 1: Build the Go binary ---
FROM golang:1.25-alpine AS builder

# Install required packages
RUN apk add --no-cache git

WORKDIR /app

# Cache dependencies separately
COPY go.mod .
COPY go.sum .
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy source code
COPY . .

# Build binary with cache
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./main.go

# --- Stage 2: Minimal runtime image ---
FROM alpine:3.20

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/server .

# Fiber runs on port 8082
EXPOSE 8082

CMD ["./server"]
