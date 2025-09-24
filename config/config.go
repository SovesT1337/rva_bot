package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"x.localhost/rvabot/internal/errors"
)

// Config содержит конфигурацию приложения
type Config struct {
	Telegram TelegramConfig
	Database DatabaseConfig
	Bot      BotConfig
	Logging  LoggingConfig
	Server   ServerConfig
}

// TelegramConfig содержит настройки Telegram API
type TelegramConfig struct {
	API   string
	Token string
}

// DatabaseConfig содержит настройки базы данных
type DatabaseConfig struct {
	Path string
}

// BotConfig содержит настройки бота
type BotConfig struct {
	Timeout    time.Duration
	MaxRetries int
}

// LoggingConfig содержит настройки логирования
type LoggingConfig struct {
	Level string
}

// ServerConfig содержит настройки HTTP сервера
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// Load загружает конфигурацию из переменных окружения
func Load() (*Config, error) {
	config := &Config{}

	// Telegram конфигурация
	config.Telegram.API = getEnv("TELEGRAM_API", "https://api.telegram.org/bot")
	config.Telegram.Token = getEnv("TELEGRAM_TOKEN", "")
	if config.Telegram.Token == "" {
		return nil, errors.NewValidationError("Отсутствует TELEGRAM_TOKEN", "Токен бота обязателен")
	}

	// Database конфигурация
	config.Database.Path = getEnv("DB_PATH", "rva_bot.db")

	// Bot конфигурация
	timeoutStr := getEnv("BOT_TIMEOUT", "30")
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		return nil, errors.NewValidationError("Неверный BOT_TIMEOUT", "Таймаут должен быть числом")
	}
	config.Bot.Timeout = time.Duration(timeout) * time.Second

	maxRetriesStr := getEnv("MAX_RETRIES", "3")
	maxRetries, err := strconv.Atoi(maxRetriesStr)
	if err != nil {
		return nil, errors.NewValidationError("Неверный MAX_RETRIES", "Количество попыток должно быть числом")
	}
	config.Bot.MaxRetries = maxRetries

	// Logging конфигурация
	config.Logging.Level = getEnv("LOG_LEVEL", "INFO")

	// Server конфигурация
	config.Server.Port = getEnv("SERVER_PORT", "8080")

	readTimeoutStr := getEnv("SERVER_READ_TIMEOUT", "5")
	readTimeout, err := strconv.Atoi(readTimeoutStr)
	if err != nil {
		return nil, errors.NewValidationError("Неверный SERVER_READ_TIMEOUT", "Таймаут должен быть числом")
	}
	config.Server.ReadTimeout = time.Duration(readTimeout) * time.Second

	writeTimeoutStr := getEnv("SERVER_WRITE_TIMEOUT", "5")
	writeTimeout, err := strconv.Atoi(writeTimeoutStr)
	if err != nil {
		return nil, errors.NewValidationError("Неверный SERVER_WRITE_TIMEOUT", "Таймаут должен быть числом")
	}
	config.Server.WriteTimeout = time.Duration(writeTimeout) * time.Second

	return config, nil
}

// getEnv получает переменную окружения с значением по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Validate проверяет корректность конфигурации
func (c *Config) Validate() error {
	// Telegram конфигурация
	if c.Telegram.Token == "" {
		return errors.NewValidationError("Отсутствует токен бота", "TELEGRAM_TOKEN обязателен")
	}

	// Проверяем формат токена бота
	if len(c.Telegram.Token) < 10 {
		return errors.NewValidationError("Неверный формат токена", "TELEGRAM_TOKEN должен быть валидным токеном бота")
	}

	if c.Telegram.API == "" {
		return errors.NewValidationError("Отсутствует API URL", "TELEGRAM_API обязателен")
	}

	// Database конфигурация
	if c.Database.Path == "" {
		return errors.NewValidationError("Отсутствует путь к БД", "DB_PATH обязателен")
	}

	// Проверяем расширение файла БД
	if !strings.HasSuffix(c.Database.Path, ".db") {
		return errors.NewValidationError("Неверное расширение БД", "DB_PATH должен заканчиваться на .db")
	}

	// Bot конфигурация
	if c.Bot.Timeout <= 0 {
		return errors.NewValidationError("Неверный таймаут", "BOT_TIMEOUT должен быть положительным")
	}

	if c.Bot.Timeout > 300*time.Second {
		return errors.NewValidationError("Слишком большой таймаут", "BOT_TIMEOUT не должен превышать 300 секунд")
	}

	if c.Bot.MaxRetries <= 0 {
		return errors.NewValidationError("Неверное количество попыток", "MAX_RETRIES должен быть положительным")
	}

	if c.Bot.MaxRetries > 10 {
		return errors.NewValidationError("Слишком много попыток", "MAX_RETRIES не должен превышать 10")
	}

	// Logging конфигурация
	validLogLevels := map[string]bool{
		"DEBUG": true,
		"INFO":  true,
		"WARN":  true,
		"ERROR": true,
	}

	if !validLogLevels[c.Logging.Level] {
		return errors.NewValidationError("Неверный уровень логирования",
			"LOG_LEVEL должен быть одним из: DEBUG, INFO, WARN, ERROR")
	}

	// Server конфигурация
	if c.Server.Port == "" {
		return errors.NewValidationError("Отсутствует порт сервера", "SERVER_PORT обязателен")
	}

	// Проверяем, что порт - это число
	if port, err := strconv.Atoi(c.Server.Port); err != nil || port <= 0 || port > 65535 {
		return errors.NewValidationError("Неверный порт", "SERVER_PORT должен быть числом от 1 до 65535")
	}

	if c.Server.ReadTimeout <= 0 {
		return errors.NewValidationError("Неверный таймаут чтения", "SERVER_READ_TIMEOUT должен быть положительным")
	}

	if c.Server.ReadTimeout > 60*time.Second {
		return errors.NewValidationError("Слишком большой таймаут чтения", "SERVER_READ_TIMEOUT не должен превышать 60 секунд")
	}

	if c.Server.WriteTimeout <= 0 {
		return errors.NewValidationError("Неверный таймаут записи", "SERVER_WRITE_TIMEOUT должен быть положительным")
	}

	if c.Server.WriteTimeout > 60*time.Second {
		return errors.NewValidationError("Слишком большой таймаут записи", "SERVER_WRITE_TIMEOUT не должен превышать 60 секунд")
	}

	return nil
}

// GetBotURL возвращает полный URL бота
func (c *Config) GetBotURL() string {
	return c.Telegram.API + c.Telegram.Token
}
