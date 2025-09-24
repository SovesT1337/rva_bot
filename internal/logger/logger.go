package logger

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"
)

// LogLevel определяет уровень логирования
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// Logger структура для логирования
type Logger struct {
	level  LogLevel
	logger *log.Logger
}

var (
	// Цвета для терминала
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorGray   = "\033[37m"
)

// New создает новый логгер
func New(level LogLevel) *Logger {
	return &Logger{
		level:  level,
		logger: log.New(os.Stdout, "", 0),
	}
}

// getColor возвращает цвет для уровня логирования
func (l *Logger) getColor(level LogLevel) string {
	switch level {
	case DEBUG:
		return colorGray
	case INFO:
		return colorBlue
	case WARN:
		return colorYellow
	case ERROR:
		return colorRed
	default:
		return colorReset
	}
}

// getLevelName возвращает название уровня
func (l *Logger) getLevelName(level LogLevel) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// getCallerInfo получает информацию о файле и строке вызова
func getCallerInfo() (string, int, bool) {
	// Пробуем разные уровни, чтобы найти реальный вызов
	for i := 3; i <= 8; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok {
			// Оставляем только имя файла
			parts := strings.Split(file, "/")
			fileName := parts[len(parts)-1]
			// Пропускаем файлы логгера
			if fileName != "logger.go" {
				return fileName, line, true
			}
		}
	}
	return "", 0, false
}

// maskSensitiveData маскирует чувствительные данные в сообщениях
func maskSensitiveData(message string) string {
	// Маскируем токены бота
	tokenRegex := regexp.MustCompile(`bot\d+:[A-Za-z0-9_-]{35}`)
	message = tokenRegex.ReplaceAllString(message, "bot***:***")

	// Маскируем API ключи
	apiKeyRegex := regexp.MustCompile(`[A-Za-z0-9]{32,}`)
	message = apiKeyRegex.ReplaceAllStringFunc(message, func(match string) string {
		if len(match) > 8 {
			return match[:4] + "***" + match[len(match)-4:]
		}
		return "***"
	})

	return message
}

// formatMessage форматирует сообщение лога
func (l *Logger) formatMessage(level LogLevel, context, message string) string {
	// Маскируем чувствительные данные
	message = maskSensitiveData(message)

	timestamp := time.Now().Format("15:04:05")
	color := l.getColor(level)
	levelName := l.getLevelName(level)

	// Получаем информацию о файле и строке
	file, line, ok := getCallerInfo()

	var contextStr string
	if context != "" {
		contextStr = fmt.Sprintf("[%s] ", context)
	}

	if ok {
		return fmt.Sprintf("%s%s %s[%s]%s %s:%d %s%s%s",
			colorGray, timestamp,
			color, levelName, colorReset,
			file, line,
			contextStr, message, colorReset)
	}

	return fmt.Sprintf("%s%s %s[%s]%s %s%s%s",
		colorGray, timestamp,
		color, levelName, colorReset,
		contextStr, message, colorReset)
}

// log выводит сообщение с указанным уровнем
func (l *Logger) log(level LogLevel, context, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	message := fmt.Sprintf(format, args...)
	formatted := l.formatMessage(level, context, message)
	l.logger.Println(formatted)
}

// Debug выводит отладочное сообщение
func (l *Logger) Debug(context, format string, args ...interface{}) {
	l.log(DEBUG, context, format, args...)
}

// Info выводит информационное сообщение
func (l *Logger) Info(context, format string, args ...interface{}) {
	l.log(INFO, context, format, args...)
}

// Warn выводит предупреждение
func (l *Logger) Warn(context, format string, args ...interface{}) {
	l.log(WARN, context, format, args...)
}

// Error выводит ошибку
func (l *Logger) Error(context, format string, args ...interface{}) {
	l.log(ERROR, context, format, args...)
}

// SetLevel устанавливает уровень логирования
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// Глобальный логгер
var defaultLogger = New(INFO)

// Глобальные функции для удобства использования
func Debug(context, format string, args ...interface{}) {
	defaultLogger.Debug(context, format, args...)
}

func Info(context, format string, args ...interface{}) {
	defaultLogger.Info(context, format, args...)
}

func Warn(context, format string, args ...interface{}) {
	defaultLogger.Warn(context, format, args...)
}

func Error(context, format string, args ...interface{}) {
	defaultLogger.Error(context, format, args...)
}

func SetLevel(level LogLevel) {
	defaultLogger.SetLevel(level)
}

// Удобные функции для логирования с контекстом

// BotInfo логирует информацию о работе бота
func BotInfo(format string, args ...interface{}) {
	Info("BOT", format, args...)
}

// BotError логирует ошибки бота
func BotError(format string, args ...interface{}) {
	Error("BOT", format, args...)
}

// TelegramInfo логирует информацию о Telegram API
func TelegramInfo(format string, args ...interface{}) {
	Info("TELEGRAM", format, args...)
}

// TelegramError логирует ошибки Telegram API
func TelegramError(format string, args ...interface{}) {
	Error("TELEGRAM", format, args...)
}

// TelegramWarn логирует предупреждения Telegram API
func TelegramWarn(format string, args ...interface{}) {
	Warn("TELEGRAM", format, args...)
}

// DatabaseInfo логирует информацию о работе с БД
func DatabaseInfo(format string, args ...interface{}) {
	Info("DATABASE", format, args...)
}

// DatabaseError логирует ошибки БД
func DatabaseError(format string, args ...interface{}) {
	Error("DATABASE", format, args...)
}

// UserInfo логирует информацию о действиях пользователей
func UserInfo(userID int, format string, args ...interface{}) {
	context := fmt.Sprintf("USER_%d", userID)
	Info(context, format, args...)
}

// UserError логирует ошибки пользователей
func UserError(userID int, format string, args ...interface{}) {
	context := fmt.Sprintf("USER_%d", userID)
	Error(context, format, args...)
}

// AdminInfo логирует информацию о действиях админов
func AdminInfo(adminID int, format string, args ...interface{}) {
	context := fmt.Sprintf("ADMIN_%d", adminID)
	Info(context, format, args...)
}

// AdminError логирует ошибки админов
func AdminError(adminID int, format string, args ...interface{}) {
	context := fmt.Sprintf("ADMIN_%d", adminID)
	Error(context, format, args...)
}

// HTTPInfo логирует информацию о HTTP запросах
func HTTPInfo(method, endpoint string, statusCode int, format string, args ...interface{}) {
	context := fmt.Sprintf("HTTP_%s_%s_%d", method, endpoint, statusCode)
	Info(context, format, args...)
}

// HTTPError логирует ошибки HTTP запросов
func HTTPError(method, endpoint string, format string, args ...interface{}) {
	context := fmt.Sprintf("HTTP_%s_%s", method, endpoint)
	Error(context, format, args...)
}
