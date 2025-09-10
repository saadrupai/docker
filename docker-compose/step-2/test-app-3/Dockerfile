FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/server .

EXPOSE 8090

CMD ["./server"]