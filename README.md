# RVA Academy Bot

Telegram-бот для управления тренировками в академии бега.

## Возможности

- 👨‍🏫 Управление тренерами
- 🏁 Управление трассами  
- 📅 Управление расписанием тренировок
- 📝 Регистрация пользователей на тренировки
- ⚙️ Админ-панель
- 🛡️ Rate limiting для защиты от спама
- 🚀 Graceful shutdown
- 🔒 Маскировка чувствительных данных в логах
- 💾 SQLite база данных для простого развертывания

## Быстрый старт

### Разработка

1. **Клонирование:**
```bash
git clone https://github.com/SovesT1337/rva_bot.git
cd rva_bot
```

2. **Запуск в режиме разработки:**
```bash
make dev
# Автоматически создаст .env файл
# Отредактируйте TELEGRAM_TOKEN в .env файле
```

### Продакшен

1. **Настройка:**
```bash
cp env.production.example .env
# Отредактируйте .env файл с вашими настройками
```

2. **Запуск:**
```bash
make prod
```

## Ручная установка

### Локально

1. **Клонирование:**
```bash
git clone https://github.com/SovesT1337/rva_bot.git
cd rva_bot
```

2. **Настройка:**
```bash
# Создайте .env файл
cp env.production.example .env
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
cp env.production.example .env
# Добавьте TELEGRAM_TOKEN=your_bot_token_here
```

2. **Запуск:**
```bash
go mod tidy
go run main.go
```

## Конфигурация

Создайте файл `.env`:
```env
TELEGRAM_TOKEN=your_bot_token_here
TELEGRAM_API=https://api.telegram.org/bot

# SQLite Configuration
DB_FILE_PATH=./rva_bot.db

LOG_LEVEL=INFO
BOT_TIMEOUT=30
MAX_RETRIES=3
SERVER_PORT=8080
SERVER_READ_TIMEOUT=5
SERVER_WRITE_TIMEOUT=5
```

## Мониторинг

Бот предоставляет простой HTTP endpoint для проверки готовности:

- `GET /ready` - Простая проверка готовности (возвращает "OK")

## Развертывание на сервере

1. **Установка Go:**
```bash
# Установите Go 1.21+ на ваш сервер
```

2. **Клонирование и настройка:**
```bash
git clone https://github.com/SovesT1337/rva_bot.git /opt/rva_bot
cd /opt/rva_bot
cp env.production.example .env
# Отредактируйте .env файл с вашими настройками
```

3. **Запуск:**
```bash
go mod tidy
go run main.go
```

### Мониторинг

```bash
# Проверка готовности
curl http://localhost:8080/ready

# Проверка логов (если используете systemd)
journalctl -u rva-bot -f
```

### Бэкапы

Ручной бэкап SQLite базы данных:
```bash
# Создание бэкапа
cp rva_bot.db backup_$(date +%Y%m%d_%H%M%S).db

# Восстановление из бэкапа
cp backup_file.db rva_bot.db
```

### Безопасность

- ✅ Panic recovery во всех горутинах
- ✅ Rate limiting для защиты от спама
- ✅ Валидация всех входных данных
- ✅ Graceful shutdown
- ✅ Маскировка токенов в логах

## Команды Makefile

- `make help` - Показать справку
- `make dev` - Запуск в режиме разработки
- `make build` - Собрать приложение
- `make run` - Запустить приложение
- `make clean` - Очистить временные файлы

## Команды бота

- `/start` - Главное меню
- `/help` - Справка
- `/admin` - Админ-панель (только для админов)

## Структура проекта

```
rva_bot/
├── config/                    # Конфигурация
├── internal/
│   ├── backoff/              # Retry механизм
│   ├── commands/             # Команды бота
│   ├── database/             # База данных (SQLite)
│   ├── errors/               # Обработка ошибок
│   ├── handler/              # Обработчик сообщений
│   ├── http/                 # HTTP клиент
│   ├── logger/               # Логирование
│   ├── metrics/              # Метрики
│   ├── ratelimit/            # Rate limiting
│   ├── recovery/             # Panic recovery
│   ├── shutdown/             # Graceful shutdown
│   ├── state/                # Менеджер состояний
│   ├── states/               # Состояния диалогов
│   ├── telegram/             # Telegram API
│   └── validation/           # Валидация данных
├── main.go                   # Точка входа
├── Makefile                  # Команды для управления
├── Dockerfile                # Docker для развертывания
└── env.production.example    # Пример конфигурации
```

## Требования

- Go 1.21+
- Telegram Bot Token от @BotFather

## Особенности

- **SQLite** - Простая и надежная база данных
- **Docker** - Простое развертывание и изоляция
- **Makefile** - Удобные команды для управления
- **Graceful Shutdown** - Корректное завершение работы
- **Rate Limiting** - Защита от спама
- **Structured Logging** - Детальное логирование

## Лицензия

MIT License