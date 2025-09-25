# RVA Bot - Минимальный Makefile

.PHONY: help build run dev prod clean docker-build docker-run docker-stop

# Переменные
APP_NAME := rva_bot

help: ## Показать справку
	@echo "RVA Bot - Команды:"
	@echo "  make dev     - Запуск в режиме разработки"
	@echo "  make prod    - Запуск в продакшене"
	@echo "  make build   - Собрать приложение"
	@echo "  make run     - Запустить приложение"
	@echo "  make clean   - Очистить временные файлы"
	@echo "  make docker-build - Собрать Docker образ"
	@echo "  make docker-run   - Запустить в Docker"
	@echo "  make docker-stop  - Остановить Docker"

dev: ## Запуск в режиме разработки
	@if [ ! -f .env ]; then \
		echo "Создание .env файла..."; \
		cp env.production.example .env; \
		sed -i 's/DB_HOST=postgres/DB_HOST=localhost/' .env; \
		sed -i 's/DB_SSLMODE=require/DB_SSLMODE=disable/' .env; \
		echo "Отредактируйте TELEGRAM_TOKEN в .env файле"; \
	fi
	docker compose up -d postgres
	sleep 5
	go run main.go

prod: ## Запуск в продакшене
	@if [ ! -f .env ]; then \
		echo "Ошибка: файл .env не найден!"; \
		echo "Создайте .env файл на основе env.production.example"; \
		exit 1; \
	fi
	docker compose -f docker-compose.prod.yml up -d

build: ## Собрать приложение
	go build -o $(APP_NAME) main.go

run: ## Запустить приложение
	@if [ ! -f .env ]; then \
		echo "Ошибка: файл .env не найден!"; \
		exit 1; \
	fi
	go run main.go

clean: ## Очистить временные файлы
	rm -f $(APP_NAME) $(APP_NAME).exe *.log

docker-build: ## Собрать Docker образ
	docker build -t rva_bot:latest .

docker-run: ## Запустить в Docker
	docker compose up -d

docker-stop: ## Остановить Docker
	docker compose down
	docker compose -f docker-compose.prod.yml down