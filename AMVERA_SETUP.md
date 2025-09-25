# Настройка RVA Bot в Amvera

## Переменные окружения для Amvera

В панели управления Amvera необходимо установить следующие переменные окружения:

### Telegram Bot Configuration
```
TELEGRAM_TOKEN=your_production_bot_token_here
TELEGRAM_API=https://api.telegram.org/bot
```

### Database Configuration (PostgreSQL)
```
DB_HOST=amvera-sovest-cnpg-rva-telegram-bot-rw
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_secure_password_here
DB_NAME=rva_bot
DB_SSLMODE=require
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

1. **DB_HOST**: Используйте доменное имя в формате `amvera-<username>-cnpg-<project_name>-rw`
   - Замените `sovest` на ваше имя пользователя в Amvera
   - Замените `rva-telegram-bot` на название вашего проекта

2. **DB_SSLMODE**: Используйте `require` для безопасного подключения в продакшене

3. **TELEGRAM_TOKEN**: Получите токен от @BotFather в Telegram

4. **DB_PASSWORD**: Установите надежный пароль для базы данных

## Проверка подключения

После настройки переменных окружения приложение должно успешно подключиться к базе данных PostgreSQL в Amvera.

Если возникают проблемы с подключением, проверьте:
- Правильность доменного имени базы данных
- Корректность пароля и имени пользователя
- Доступность базы данных в панели Amvera
