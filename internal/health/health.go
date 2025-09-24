package health

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"

	"x.localhost/rvabot/internal/logger"
)

// Status статус проверки
type Status string

const (
	StatusHealthy   Status = "healthy"
	StatusUnhealthy Status = "unhealthy"
	StatusDegraded  Status = "degraded"
)

// Check интерфейс для health check
type Check interface {
	Name() string
	Check(ctx context.Context) CheckResult
}

// CheckResult результат проверки
type CheckResult struct {
	Name      string                 `json:"name"`
	Status    Status                 `json:"status"`
	Message   string                 `json:"message,omitempty"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Duration  time.Duration          `json:"duration"`
	Timestamp time.Time              `json:"timestamp"`
}

// Manager управляет health checks
type Manager struct {
	checks  map[string]Check
	mu      sync.RWMutex
	timeout time.Duration
}

// NewManager создает новый менеджер health checks
func NewManager(timeout time.Duration) *Manager {
	return &Manager{
		checks:  make(map[string]Check),
		timeout: timeout,
	}
}

// RegisterCheck регистрирует новый health check
func (m *Manager) RegisterCheck(check Check) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.checks[check.Name()] = check
	logger.BotInfo("Зарегистрирован health check: %s", check.Name())
}

// CheckAll выполняет все зарегистрированные проверки
func (m *Manager) CheckAll(ctx context.Context) map[string]CheckResult {
	m.mu.RLock()
	checks := make(map[string]Check, len(m.checks))
	for name, check := range m.checks {
		checks[name] = check
	}
	m.mu.RUnlock()

	results := make(map[string]CheckResult, len(checks))
	var wg sync.WaitGroup
	var mu sync.Mutex

	for name, check := range checks {
		wg.Add(1)
		go func(name string, check Check) {
			defer wg.Done()

			checkCtx, cancel := context.WithTimeout(ctx, m.timeout)
			defer cancel()

			start := time.Now()
			result := check.Check(checkCtx)
			result.Duration = time.Since(start)
			result.Timestamp = time.Now()

			mu.Lock()
			results[name] = result
			mu.Unlock()
		}(name, check)
	}

	wg.Wait()
	return results
}

// GetOverallStatus возвращает общий статус системы
func (m *Manager) GetOverallStatus(ctx context.Context) (Status, map[string]CheckResult) {
	results := m.CheckAll(ctx)

	hasUnhealthy := false
	hasDegraded := false

	for _, result := range results {
		switch result.Status {
		case StatusUnhealthy:
			hasUnhealthy = true
		case StatusDegraded:
			hasDegraded = true
		}
	}

	var overallStatus Status
	if hasUnhealthy {
		overallStatus = StatusUnhealthy
	} else if hasDegraded {
		overallStatus = StatusDegraded
	} else {
		overallStatus = StatusHealthy
	}

	return overallStatus, results
}

// HTTPHandler возвращает HTTP handler для health checks
func (m *Manager) HTTPHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		overallStatus, results := m.GetOverallStatus(ctx)

		response := map[string]interface{}{
			"status":    overallStatus,
			"timestamp": time.Now(),
			"checks":    results,
		}

		w.Header().Set("Content-Type", "application/json")

		statusCode := http.StatusOK
		if overallStatus == StatusUnhealthy {
			statusCode = http.StatusServiceUnavailable
		} else if overallStatus == StatusDegraded {
			statusCode = http.StatusOK // 200, но с предупреждением
		}

		w.WriteHeader(statusCode)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			logger.BotError("Ошибка кодирования health check response: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

// DatabaseCheck проверяет состояние базы данных
type DatabaseCheck struct {
	checkFunc func(ctx context.Context) error
}

// NewDatabaseCheck создает новый check для базы данных
func NewDatabaseCheck(checkFunc func(ctx context.Context) error) *DatabaseCheck {
	return &DatabaseCheck{checkFunc: checkFunc}
}

// Name возвращает имя проверки
func (dc *DatabaseCheck) Name() string {
	return "database"
}

// Check выполняет проверку базы данных
func (dc *DatabaseCheck) Check(ctx context.Context) CheckResult {
	if err := dc.checkFunc(ctx); err != nil {
		return CheckResult{
			Name:    dc.Name(),
			Status:  StatusUnhealthy,
			Message: fmt.Sprintf("Database connection failed: %v", err),
		}
	}

	return CheckResult{
		Name:    dc.Name(),
		Status:  StatusHealthy,
		Message: "Database connection is healthy",
	}
}

// TelegramCheck проверяет состояние Telegram API
type TelegramCheck struct {
	checkFunc func(ctx context.Context) error
}

// NewTelegramCheck создает новый check для Telegram API
func NewTelegramCheck(checkFunc func(ctx context.Context) error) *TelegramCheck {
	return &TelegramCheck{checkFunc: checkFunc}
}

// Name возвращает имя проверки
func (tc *TelegramCheck) Name() string {
	return "telegram"
}

// Check выполняет проверку Telegram API
func (tc *TelegramCheck) Check(ctx context.Context) CheckResult {
	if err := tc.checkFunc(ctx); err != nil {
		return CheckResult{
			Name:    tc.Name(),
			Status:  StatusUnhealthy,
			Message: fmt.Sprintf("Telegram API check failed: %v", err),
		}
	}

	return CheckResult{
		Name:    tc.Name(),
		Status:  StatusHealthy,
		Message: "Telegram API is accessible",
	}
}

// MemoryCheck проверяет использование памяти
type MemoryCheck struct {
	maxMemoryMB int64
}

// NewMemoryCheck создает новый check для памяти
func NewMemoryCheck(maxMemoryMB int64) *MemoryCheck {
	return &MemoryCheck{maxMemoryMB: maxMemoryMB}
}

// Name возвращает имя проверки
func (mc *MemoryCheck) Name() string {
	return "memory"
}

// Check выполняет проверку памяти
func (mc *MemoryCheck) Check(ctx context.Context) CheckResult {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	memoryMB := int64(m.Alloc / 1024 / 1024)

	details := map[string]interface{}{
		"alloc_mb":      memoryMB,
		"max_memory_mb": mc.maxMemoryMB,
		"heap_objects":  m.HeapObjects,
		"gc_cycles":     m.NumGC,
	}

	if memoryMB > mc.maxMemoryMB {
		return CheckResult{
			Name:    mc.Name(),
			Status:  StatusDegraded,
			Message: fmt.Sprintf("High memory usage: %dMB", memoryMB),
			Details: details,
		}
	}

	return CheckResult{
		Name:    mc.Name(),
		Status:  StatusHealthy,
		Message: fmt.Sprintf("Memory usage is normal: %dMB", memoryMB),
		Details: details,
	}
}
