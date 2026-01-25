FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o webhook-inbox ./cmd/webhook-inbox

FROM debian:bullseye-slim

WORKDIR /app

COPY --from=builder /app/webhook-inbox .

EXPOSE 8080

CMD ["./webhook-inbox"]
