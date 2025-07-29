# Subscription Management API (API для управления подписками)

[![Typing SVG](https://readme-typing-svg.herokuapp.com?color=%2336BCF7&lines=Subscription+Management+API)](https://git.io/typing-svg)

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

## 🚀 Быстрый старт

### 1. Клонирование репозитория

```bash
git clone https://github.com/BabichevDima/subManager.git

cd subManager/
```

### 2.1. Запуск (Вариант A)

```
# Запуск
docker compose up -d --build
```

### 2.2. Запуск (Вариант B)

```
# Сборка и запуск
go build -o ./out ./cmd/submanager && ./out

# Запуск
./out
```

### 2.3. Запуск (Вариант C)

```
# Запуск
make run
```

### 3. 📚 Документация

🔗 http://localhost:8080/swagger/index.html
