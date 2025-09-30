# Build stage
FROM golang:1.24 AS builder

WORKDIR /app

# Copy go.mod and go.sum, install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary for Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app main.go

# Final image
FROM alpine:latest

WORKDIR /app

# Copy binary and frontend
COPY --from=builder /app/app .
COPY --from=builder /app/web ./web

# Установка зависимостей для запуска Go-приложения
RUN apk add --no-cache ca-certificates

# Default environment variables (can be overridden at runtime)
ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/scheduler.db
ENV TODO_PASSWORD=changeme
ENV TODO_SECRET=changeme_secret

# Start application
CMD ["./app"]
