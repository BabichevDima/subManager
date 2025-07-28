# Subscription Management API (API для управления подписками)

## Описание проекта
RESTful API для управления подписками пользователей. Позволяет:
- Создавать/просматривать/обновлять/удалять подписки
- Получать списки подписок с пагинацией
- Управлять статусами подписок

## Технологии
- **Язык**: Go 1.21+
- **Фреймворк**: Чистая архитектура (Clean Architecture)
- **База данных**: PostgreSQL 15+
- **Документация**: Swagger (OpenAPI 3.0)

## Установка и запуск

### Требования
- Установленный Go 1.21+
- PostgreSQL 15+
- Make (опционально)

### 1. Клонирование репозитория
```bash
git clone https://github.com/your-username/subscription-service.git
cd subscription-service

### 2. Запуск

```bash
# Установка зависимостей
go mod download

# Запуск
go build -o ./out ./cmd/submanager && ./out

# или
make run

### 3. Документация

http://localhost:8080/swagger/index.html
