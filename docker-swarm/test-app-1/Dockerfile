# -------------------------------
# 1. Build stage
# -------------------------------
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git (needed for some go modules)
RUN apk add --no-cache git

# Copy go.mod and go.sum first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the project
COPY . .

# Build binary
RUN go build -o server .

# -------------------------------
# 2. Run stage
# -------------------------------
FROM alpine:latest

WORKDIR /root/

# Copy the binary only
COPY --from=builder /app/server .

# Expose the default app port (optional)
EXPOSE 8080

# Run the server
CMD ["./server"]
