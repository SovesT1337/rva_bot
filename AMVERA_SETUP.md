# Настройка RVA Bot в Amvera

## Переменные окружения для Amvera

В панели управления Amvera необходимо установить следующие переменные окружения:

### Telegram Bot Configuration
```
TELEGRAM_TOKEN=your_production_bot_token_here
TELEGRAM_API=https://api.telegram.org/bot
```

### Database Configuration (SQLite)
```
DB_FILE_PATH=/data/rva_bot.db
```

### Logging Configuration
```
LOG_LEVEL=INFO
```

### Bot Configuration
```
BOT_TIMEOUT=30
MAX_RETRIES=3
```

### Server Configuration
```
SERVER_PORT=8080
SERVER_READ_TIMEOUT=10
SERVER_WRITE_TIMEOUT=10
```

## Важные замечания

1. **DB_FILE_PATH**: Путь к файлу SQLite базы данных (будет создан автоматически)
3. **TELEGRAM_TOKEN**: Получите токен от @BotFather в Telegram
4. **Persistent Storage**: Amvera автоматически монтирует `/data` для постоянного хранения

## Преимущества SQLite

- ✅ Простое развертывание - не нужна отдельная база данных
- ✅ Автоматическое создание файла базы данных
- ✅ Данные сохраняются в persistent storage Amvera
- ✅ Нет проблем с подключением к внешним сервисам
- ✅ Быстрая работа для небольших проектов

## Проверка подключения

После настройки переменных окружения приложение должно успешно создать SQLite базу данных в `/data/rva_bot.db`.

Если возникают проблемы, проверьте:
- Правильность токена Telegram бота
- Доступность persistent storage в Amvera
