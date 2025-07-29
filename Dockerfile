FROM golang:1.24.4 AS builder
WORKDIR /app

# Кэширование зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем бинарник в текущую директорию
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o ./submanager ./cmd/submanager

# Финальный образ
FROM alpine:3.19
WORKDIR /app

# Копируем из текущей директории builder'а
COPY --from=builder /app/submanager .
COPY --from=builder /app/app ./app/
COPY --from=builder /app/config ./config/
COPY --from=builder /app/config/config.docker.yaml ./config/config.dev.yaml

RUN apk add --no-cache ca-certificates && \
    chmod +x submanager

EXPOSE 8080
CMD ["./submanager"]