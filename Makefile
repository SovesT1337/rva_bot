# RVA Bot - Минимальный Makefile

.PHONY: help build run dev clean docker-build

# Переменные
APP_NAME := rva_bot

help: ## Показать справку
	@echo "RVA Bot - Команды:"
	@echo "  make dev     - Запуск в режиме разработки"
	@echo "  make build   - Собрать приложение"
	@echo "  make run     - Запустить приложение"
	@echo "  make clean   - Очистить временные файлы"
	@echo "  make docker-build - Собрать Docker образ"

dev: ## Запуск в режиме разработки
	@if [ ! -f .env ]; then \
		echo "Создание .env файла..."; \
		cp env.production.example .env; \
		echo "Отредактируйте TELEGRAM_TOKEN в .env файле"; \
	fi
	go run main.go

build: ## Собрать приложение
	go build -o $(APP_NAME) main.go

run: ## Запустить приложение
	@if [ ! -f .env ]; then \
		echo "Ошибка: файл .env не найден!"; \
		exit 1; \
	fi
	go run main.go

clean: ## Очистить временные файлы
	rm -f $(APP_NAME) $(APP_NAME).exe *.log *.db *.db-shm *.db-wal

docker-build: ## Собрать Docker образ
	docker build -t rva_bot:latest .