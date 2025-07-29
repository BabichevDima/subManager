.PHONY: swagger run build

# Генерация Swagger документации
swagger:
	@echo "Generating Swagger docs..."
	@swag init -g cmd/submanager/main.go --output internal/docs
	@echo "Docs generated at internal/docs/"

# Запуск приложения
run: swagger
	@echo "Starting server..."
	@go run cmd/submanager/main.go