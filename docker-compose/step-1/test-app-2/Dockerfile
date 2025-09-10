# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o server .

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8090
CMD ["./server"]
