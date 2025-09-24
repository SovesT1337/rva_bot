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

## Установка

### Локально

1. **Клонирование:**
```bash
git clone <your-repo-url>
cd rva_bot
```

2. **Настройка:**
```bash
cp .env.example .env
# Отредактируйте .env файл, добавив токен бота
```

3. **Запуск:**
```bash
go mod tidy
go run main.go
```

### Docker

1. **Создание .env файла:**
```bash
cp .env.example .env
# Добавьте TELEGRAM_TOKEN=your_bot_token_here
```

2. **Запуск:**
```bash
docker-compose up -d
```

3. **Остановка:**
```bash
docker-compose down
```

## Конфигурация

Создайте файл `.env`:
```env
TELEGRAM_TOKEN=your_bot_token_here
TELEGRAM_API=https://api.telegram.org/bot
DB_PATH=rva_bot.db
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

## Production Deployment

### Docker Production

1. **Создайте production конфигурацию:**
```bash
cp env.production.example .env.production
# Отредактируйте .env.production, добавив реальный токен
```

2. **Запуск в production:**
```bash
docker-compose -f docker-compose.prod.yml up -d
```

3. **Мониторинг:**
```bash
# Проверка статуса
docker-compose -f docker-compose.prod.yml ps

# Логи
docker-compose -f docker-compose.prod.yml logs -f

# Health check
curl http://localhost:8080/health
```

### Бэкапы

Автоматический бэкап базы данных:
```bash
# Ручной бэкап
./scripts/backup.sh

# Автоматический бэкап (cron)
0 2 * * * /path/to/rva_bot/scripts/backup.sh
```

### Безопасность

- ✅ Panic recovery во всех горутинах
- ✅ Rate limiting для защиты от спама
- ✅ Валидация всех входных данных
- ✅ Graceful shutdown
- ✅ Health checks и мониторинг
- ✅ Маскировка токенов в логах

## Команды

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
- SQLite3
- Telegram Bot Token от @BotFather

## Лицензия

MIT License