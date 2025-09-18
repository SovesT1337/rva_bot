package database

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB(filepath string) error {
	var err error

	// Настройки для лучшей производительности SQLite
	config := &gorm.Config{
		// Отключаем логирование для продакшена (можно включить для отладки)
		// Logger: logger.Default.LogMode(logger.Silent),
	}

	// Настройки SQLite для конкурентного доступа
	sqliteConfig := sqlite.Open(filepath + "?_journal_mode=WAL&_synchronous=NORMAL&_cache_size=1000&_foreign_keys=ON")

	db, err = gorm.Open(sqliteConfig, config)
	if err != nil {
		return fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	// Настройка пула соединений
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("ошибка получения DB: %w", err)
	}

	// Настройки пула соединений для конкурентного доступа
	sqlDB.SetMaxOpenConns(25)   // Максимум открытых соединений
	sqlDB.SetMaxIdleConns(10)   // Максимум неактивных соединений
	sqlDB.SetConnMaxLifetime(0) // Время жизни соединения (0 = без ограничений)

	// Тестируем соединение с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("ошибка тестирования соединения: %w", err)
	}

	if err := db.AutoMigrate(&Track{}); err != nil {
		return fmt.Errorf("ошибка миграции: %w", err)
	}
	if err := db.AutoMigrate(&Trainer{}); err != nil {
		return fmt.Errorf("ошибка миграции: %w", err)
	}
	if err := db.AutoMigrate(&Admin{}); err != nil {
		return fmt.Errorf("ошибка миграции: %w", err)
	}
	if err := db.AutoMigrate(&Training{}); err != nil {
		return fmt.Errorf("ошибка миграции: %w", err)
	}
	if err := db.AutoMigrate(&User{}); err != nil {
		return fmt.Errorf("ошибка миграции: %w", err)
	}
	if err := db.AutoMigrate(&TrainingRegistration{}); err != nil {
		return fmt.Errorf("ошибка миграции: %w", err)
	}

	return nil
}

// GetDBWithTimeout возвращает контекст с таймаутом для операций с БД
func GetDBWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}
