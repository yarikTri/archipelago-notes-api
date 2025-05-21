# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.19 AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o archipelago-auth ./cmd/auth/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

ENV LISTEN_PORT=8888

# Copy binary from builder
COPY --from=builder /app/archipelago-auth .

EXPOSE ${LISTEN_PORT}

CMD ["./archipelago-auth"]
