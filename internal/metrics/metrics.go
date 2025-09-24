package metrics

import (
	"sync"
	"sync/atomic"
	"time"

	"x.localhost/rvabot/internal/logger"
)

// Metrics собирает метрики приложения
type Metrics struct {
	// Счетчики
	TotalUpdates     int64
	ProcessedUpdates int64
	FailedUpdates    int64
	RateLimitedUsers int64
	DatabaseQueries  int64
	DatabaseErrors   int64
	TelegramRequests int64
	TelegramErrors   int64

	// Время
	LastUpdateTime time.Time
	Uptime         time.Time

	// Состояния
	ActiveUsers  int64
	ActiveStates int64

	mutex sync.RWMutex
}

// GlobalMetrics глобальные метрики приложения
var GlobalMetrics = &Metrics{
	Uptime: time.Now(),
}

// IncrementTotalUpdates увеличивает счетчик общих обновлений
func (m *Metrics) IncrementTotalUpdates() {
	atomic.AddInt64(&m.TotalUpdates, 1)
	m.mutex.Lock()
	m.LastUpdateTime = time.Now()
	m.mutex.Unlock()
}

// IncrementProcessedUpdates увеличивает счетчик обработанных обновлений
func (m *Metrics) IncrementProcessedUpdates() {
	atomic.AddInt64(&m.ProcessedUpdates, 1)
}

// IncrementFailedUpdates увеличивает счетчик неудачных обновлений
func (m *Metrics) IncrementFailedUpdates() {
	atomic.AddInt64(&m.FailedUpdates, 1)
}

// IncrementRateLimitedUsers увеличивает счетчик заблокированных пользователей
func (m *Metrics) IncrementRateLimitedUsers() {
	atomic.AddInt64(&m.RateLimitedUsers, 1)
}

// IncrementDatabaseQueries увеличивает счетчик запросов к БД
func (m *Metrics) IncrementDatabaseQueries() {
	atomic.AddInt64(&m.DatabaseQueries, 1)
}

// IncrementDatabaseErrors увеличивает счетчик ошибок БД
func (m *Metrics) IncrementDatabaseErrors() {
	atomic.AddInt64(&m.DatabaseErrors, 1)
}

// IncrementTelegramRequests увеличивает счетчик запросов к Telegram
func (m *Metrics) IncrementTelegramRequests() {
	atomic.AddInt64(&m.TelegramRequests, 1)
}

// IncrementTelegramErrors увеличивает счетчик ошибок Telegram
func (m *Metrics) IncrementTelegramErrors() {
	atomic.AddInt64(&m.TelegramErrors, 1)
}

// SetActiveUsers устанавливает количество активных пользователей
func (m *Metrics) SetActiveUsers(count int64) {
	atomic.StoreInt64(&m.ActiveUsers, count)
}

// SetActiveStates устанавливает количество активных состояний
func (m *Metrics) SetActiveStates(count int64) {
	atomic.StoreInt64(&m.ActiveStates, count)
}

// GetStats возвращает текущие метрики
func (m *Metrics) GetStats() map[string]interface{} {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	uptime := time.Since(m.Uptime)

	return map[string]interface{}{
		"uptime_seconds":     int64(uptime.Seconds()),
		"total_updates":      atomic.LoadInt64(&m.TotalUpdates),
		"processed_updates":  atomic.LoadInt64(&m.ProcessedUpdates),
		"failed_updates":     atomic.LoadInt64(&m.FailedUpdates),
		"rate_limited_users": atomic.LoadInt64(&m.RateLimitedUsers),
		"database_queries":   atomic.LoadInt64(&m.DatabaseQueries),
		"database_errors":    atomic.LoadInt64(&m.DatabaseErrors),
		"telegram_requests":  atomic.LoadInt64(&m.TelegramRequests),
		"telegram_errors":    atomic.LoadInt64(&m.TelegramErrors),
		"active_users":       atomic.LoadInt64(&m.ActiveUsers),
		"active_states":      atomic.LoadInt64(&m.ActiveStates),
		"last_update_time":   m.LastUpdateTime,
		"success_rate":       m.getSuccessRate(),
		"error_rate":         m.getErrorRate(),
	}
}

// getSuccessRate возвращает процент успешных операций
func (m *Metrics) getSuccessRate() float64 {
	total := atomic.LoadInt64(&m.TotalUpdates)
	if total == 0 {
		return 100.0
	}

	processed := atomic.LoadInt64(&m.ProcessedUpdates)
	return float64(processed) / float64(total) * 100.0
}

// getErrorRate возвращает процент ошибок
func (m *Metrics) getErrorRate() float64 {
	total := atomic.LoadInt64(&m.TotalUpdates)
	if total == 0 {
		return 0.0
	}

	failed := atomic.LoadInt64(&m.FailedUpdates)
	return float64(failed) / float64(total) * 100.0
}

// LogMetrics логирует текущие метрики
func (m *Metrics) LogMetrics() {
	stats := m.GetStats()
	logger.BotInfo("=== METRICS ===")
	logger.BotInfo("Uptime: %ds", stats["uptime_seconds"])
	logger.BotInfo("Total Updates: %d", stats["total_updates"])
	logger.BotInfo("Processed: %d", stats["processed_updates"])
	logger.BotInfo("Failed: %d", stats["failed_updates"])
	logger.BotInfo("Success Rate: %.2f%%", stats["success_rate"])
	logger.BotInfo("Error Rate: %.2f%%", stats["error_rate"])
	logger.BotInfo("Active Users: %d", stats["active_users"])
	logger.BotInfo("Active States: %d", stats["active_states"])
	logger.BotInfo("DB Queries: %d", stats["database_queries"])
	logger.BotInfo("DB Errors: %d", stats["database_errors"])
	logger.BotInfo("Telegram Requests: %d", stats["telegram_requests"])
	logger.BotInfo("Telegram Errors: %d", stats["telegram_errors"])
	logger.BotInfo("Rate Limited: %d", stats["rate_limited_users"])
}

// StartMetricsLogger запускает периодическое логирование метрик
func (m *Metrics) StartMetricsLogger(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			m.LogMetrics()
		}
	}()
}

// StartMetricsLogger запускает периодическое логирование глобальных метрик
func StartMetricsLogger(interval time.Duration) {
	GlobalMetrics.StartMetricsLogger(interval)
}

// Глобальные функции для удобства
func IncrementTotalUpdates() {
	GlobalMetrics.IncrementTotalUpdates()
}

func IncrementProcessedUpdates() {
	GlobalMetrics.IncrementProcessedUpdates()
}

func IncrementFailedUpdates() {
	GlobalMetrics.IncrementFailedUpdates()
}

func IncrementRateLimitedUsers() {
	GlobalMetrics.IncrementRateLimitedUsers()
}

func IncrementDatabaseQueries() {
	GlobalMetrics.IncrementDatabaseQueries()
}

func IncrementDatabaseErrors() {
	GlobalMetrics.IncrementDatabaseErrors()
}

func IncrementTelegramRequests() {
	GlobalMetrics.IncrementTelegramRequests()
}

func IncrementTelegramErrors() {
	GlobalMetrics.IncrementTelegramErrors()
}

func SetActiveUsers(count int64) {
	GlobalMetrics.SetActiveUsers(count)
}

func SetActiveStates(count int64) {
	GlobalMetrics.SetActiveStates(count)
}

func GetStats() map[string]interface{} {
	return GlobalMetrics.GetStats()
}
