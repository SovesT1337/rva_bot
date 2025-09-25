# RVA Academy Bot

Telegram-бот для управления тренировками в академии бега.

## Возможности

- 👨‍🏫 Управление тренерами
- 🏁 Управление трассами  
- 📅 Управление расписанием тренировок
- 📝 Регистрация пользователей на тренировки
- ⚙️ Админ-панель
- 🛡️ Rate limiting для защиты от спама
- 🔍 Health checks и мониторинг
- 🚀 Graceful shutdown
- 🔒 Маскировка чувствительных данных в логах

## Быстрый старт

### С помощью Makefile (рекомендуется)

1. **Клонирование и настройка:**
```bash
git clone <your-repo-url>
cd rva_bot
make install-deps
```

2. **Запуск в режиме разработки:**
```bash
make dev
# Автоматически создаст .env файл и запустит PostgreSQL
# Отредактируйте TELEGRAM_TOKEN в .env файле
```

3. **Запуск в продакшене:**
```bash
make prod
```

### Ручная установка

#### Локально

1. **Клонирование:**
```bash
git clone <your-repo-url>
cd rva_bot
```

2. **Настройка базы данных:**
```bash
# Запустите PostgreSQL через Docker
docker compose up -d postgres

# Создайте .env файл
cp env.production.example .env
# Отредактируйте .env файл, добавив токен бота и настройки БД
```

3. **Запуск:**
```bash
go mod tidy
go run main.go
```

#### Docker

1. **Создание .env файла:**
```bash
cp env.production.example .env
# Добавьте TELEGRAM_TOKEN=your_bot_token_here
# Настройте параметры PostgreSQL если нужно
```

2. **Запуск:**
```bash
docker compose up -d
```

3. **Остановка:**
```bash
docker compose down
```

## Конфигурация

Создайте файл `.env`:
```env
TELEGRAM_TOKEN=your_bot_token_here
TELEGRAM_API=https://api.telegram.org/bot

# PostgreSQL Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=rva_bot
DB_SSLMODE=disable

LOG_LEVEL=INFO
BOT_TIMEOUT=30
MAX_RETRIES=3
SERVER_PORT=8080
SERVER_READ_TIMEOUT=5
SERVER_WRITE_TIMEOUT=5
```

## Мониторинг

Бот предоставляет HTTP endpoints для мониторинга:

- `GET /health` - Детальная информация о состоянии системы
- `GET /ready` - Простая проверка готовности

Пример ответа `/health`:
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "checks": {
    "database": {
      "name": "database",
      "status": "healthy",
      "message": "Database connection is healthy",
      "duration": "2ms"
    },
    "telegram": {
      "name": "telegram", 
      "status": "healthy",
      "message": "Telegram API is accessible",
      "duration": "150ms"
    }
  }
}
```

## Развертывание на сервере

1. **Установка Docker:**
```bash
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

2. **Клонирование и настройка:**
```bash
git clone <your-repo-url> /opt/rva_bot
cd /opt/rva_bot
cp env.production.example .env
# Отредактируйте .env файл с вашими настройками
```

3. **Запуск:**
```bash
make prod
```

### Мониторинг

```bash
# Проверка статуса
docker compose -f docker-compose.prod.yml ps

# Логи
docker compose -f docker-compose.prod.yml logs -f

# Health check
curl http://localhost:8080/health
```

### Бэкапы

Ручной бэкап базы данных:
```bash
# Создание бэкапа
docker exec rva_bot_postgres_prod pg_dump -U postgres rva_bot > backup_$(date +%Y%m%d_%H%M%S).sql

# Восстановление из бэкапа
docker exec -i rva_bot_postgres_prod psql -U postgres rva_bot < backup_file.sql
```

### Безопасность

- ✅ Panic recovery во всех горутинах
- ✅ Rate limiting для защиты от спама
- ✅ Валидация всех входных данных
- ✅ Graceful shutdown
- ✅ Health checks и мониторинг
- ✅ Маскировка токенов в логах

## Команды Makefile

- `make help` - Показать справку
- `make dev` - Запуск в режиме разработки
- `make prod` - Запуск в продакшене
- `make build` - Собрать приложение
- `make run` - Запустить приложение
- `make clean` - Очистить временные файлы
- `make docker-build` - Собрать Docker образ
- `make docker-run` - Запустить в Docker
- `make docker-stop` - Остановить Docker

## Команды бота

- `/start` - Главное меню
- `/help` - Справка
- `/admin` - Админ-панель (только для админов)

## Структура

```
internal/
├── commands/     # Команды бота
├── database/     # База данных
├── errors/       # Обработка ошибок
├── handler/      # Обработчик сообщений
├── health/       # Health checks
├── logger/       # Логирование
├── ratelimit/    # Rate limiting
├── shutdown/     # Graceful shutdown
├── states/       # Состояния диалогов
├── telegram/     # Telegram API
└── validation/   # Валидация данных
```

## Требования

- Go 1.21+
- PostgreSQL 12+ (или Docker для запуска PostgreSQL)
- Telegram Bot Token от @BotFather

## Лицензия

MIT License