# Build Stage
FROM golang:1.25-alpine AS builder

# Install git for fetching dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
# CGO_ENABLED=0 for static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Runtime Stage
FROM gcr.io/distroless/static-debian12

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Copy config files if needed (e.g., .env is usually injected via secrets, but config.yaml might be needed)
# For now, we assume env vars are sufficient or mounted secrets.

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./server"]
