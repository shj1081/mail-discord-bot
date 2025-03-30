# Build stage
FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bot ./cmd/main.go

# Final stage
FROM alpine:3.19
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/bot .
CMD ["./bot"]