package database

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Database представляет подключение к базе данных
type Database struct {
	db *gorm.DB
}

// NewDatabase создает новое подключение к базе данных
func NewDatabase(dsn string, dbType string) (*Database, error) {
	config := &gorm.Config{
		// Отключаем логирование для продакшена (можно включить для отладки)
		// Logger: logger.Default.LogMode(logger.Silent),
	}

	var db *gorm.DB
	var err error

	if dbType == "sqlite" {
		db, err = gorm.Open(sqlite.Open(dsn), config)
	} else {
		db, err = gorm.Open(postgres.Open(dsn), config)
	}

	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(0)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ошибка тестирования соединения: %w", err)
	}

	database := &Database{db: db}

	// Выполняем миграции
	if err := database.migrate(); err != nil {
		return nil, fmt.Errorf("ошибка миграции: %w", err)
	}

	return database, nil
}

// migrate выполняет миграции базы данных
func (d *Database) migrate() error {
	models := []interface{}{
		&Track{},
		&Trainer{},
		&Admin{},
		&Training{},
		&User{},
		&TrainingRegistration{},
		&TrainingRequest{},
	}

	for _, model := range models {
		if err := d.db.AutoMigrate(model); err != nil {
			return fmt.Errorf("ошибка миграции %T: %w", model, err)
		}
	}

	return nil
}

// GetDB возвращает экземпляр GORM DB
func (d *Database) GetDB() *gorm.DB {
	return d.db
}

// Close закрывает подключение к базе данных
func (d *Database) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Ping проверяет подключение к базе данных
func (d *Database) Ping(ctx context.Context) error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}
