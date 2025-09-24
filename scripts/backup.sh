#!/bin/bash

# Скрипт для бэкапа базы данных RVA Bot

set -e

# Конфигурация
DB_PATH="${DB_PATH:-/data/rva_bot.db}"
BACKUP_DIR="${BACKUP_DIR:-/backups}"
RETENTION_DAYS="${RETENTION_DAYS:-30}"

# Создаем директорию для бэкапов если не существует
mkdir -p "$BACKUP_DIR"

# Генерируем имя файла бэкапа
BACKUP_FILE="$BACKUP_DIR/rva_bot_$(date +%Y%m%d_%H%M%S).db"

echo "Создание бэкапа базы данных..."
echo "Источник: $DB_PATH"
echo "Назначение: $BACKUP_FILE"

# Проверяем существование файла БД
if [ ! -f "$DB_PATH" ]; then
    echo "ОШИБКА: Файл базы данных не найден: $DB_PATH"
    exit 1
fi

# Создаем бэкап
cp "$DB_PATH" "$BACKUP_FILE"

# Проверяем размер бэкапа
BACKUP_SIZE=$(stat -f%z "$BACKUP_FILE" 2>/dev/null || stat -c%s "$BACKUP_FILE" 2>/dev/null)
echo "Бэкап создан успешно. Размер: $BACKUP_SIZE байт"

# Сжимаем бэкап
echo "Сжатие бэкапа..."
gzip "$BACKUP_FILE"
BACKUP_FILE="${BACKUP_FILE}.gz"

COMPRESSED_SIZE=$(stat -f%z "$BACKUP_FILE" 2>/dev/null || stat -c%s "$BACKUP_FILE" 2>/dev/null)
echo "Сжатый размер: $COMPRESSED_SIZE байт"

# Удаляем старые бэкапы
echo "Удаление старых бэкапов (старше $RETENTION_DAYS дней)..."
find "$BACKUP_DIR" -name "rva_bot_*.db.gz" -type f -mtime +$RETENTION_DAYS -delete

echo "Бэкап завершен успешно: $BACKUP_FILE"
