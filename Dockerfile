FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Install swag and generate docs
RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    swag init -g cmd/submanager/main.go --output ./internal/docs

# Build app
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o out ./cmd/submanager

FROM alpine:3.18
WORKDIR /app

COPY --from=builder /app/out .
COPY --from=builder /app/internal/docs ./internal/docs

EXPOSE 8080
CMD ["./out"]