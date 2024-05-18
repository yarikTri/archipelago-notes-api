# syntax=docker/dockerfile:1

FROM golang:1.19

# Set destination for COPY
WORKDIR /app

ENV LISTEN_PORT=8080

COPY . ./

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o archipelago-api ./cmd/api/main.go

EXPOSE ${LISTEN_PORT}

CMD ["./archipelago-api"]
