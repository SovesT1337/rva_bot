package recovery

import (
	"context"
	"fmt"
	"runtime"
	"runtime/debug"

	"x.localhost/rvabot/internal/logger"
)

// Recoverer интерфейс для обработки паник
type Recoverer interface {
	Recover(ctx context.Context, panicValue interface{})
}

// DefaultRecoverer стандартный обработчик паник
type DefaultRecoverer struct {
	serviceName string
}

// NewDefaultRecoverer создает новый обработчик паник
func NewDefaultRecoverer(serviceName string) *DefaultRecoverer {
	return &DefaultRecoverer{
		serviceName: serviceName,
	}
}

// Recover обрабатывает панику
func (r *DefaultRecoverer) Recover(ctx context.Context, panicValue interface{}) {
	stack := debug.Stack()

	logger.Error("PANIC", "Service: %s, Panic: %v\nStack: %s",
		r.serviceName, panicValue, string(stack))

	// В продакшене можно отправить алерт в систему мониторинга
	// sendAlertToMonitoring(panicValue, stack)
}

// RecoverFunc выполняет функцию с обработкой паник
func RecoverFunc(ctx context.Context, serviceName string, fn func()) {
	defer func() {
		if r := recover(); r != nil {
			recoverer := NewDefaultRecoverer(serviceName)
			recoverer.Recover(ctx, r)
		}
	}()

	fn()
}

// RecoverFuncWithError выполняет функцию с обработкой паник и возвратом ошибки
func RecoverFuncWithError(ctx context.Context, serviceName string, fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			recoverer := NewDefaultRecoverer(serviceName)
			recoverer.Recover(ctx, r)
			err = fmt.Errorf("panic recovered: %v", r)
		}
	}()

	return fn()
}

// RecoverGoroutine запускает горутину с обработкой паник
func RecoverGoroutine(ctx context.Context, serviceName string, fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				recoverer := NewDefaultRecoverer(serviceName)
				recoverer.Recover(ctx, r)
			}
		}()

		fn()
	}()
}

// GetGoroutineStats возвращает статистику горутин
func GetGoroutineStats() map[string]interface{} {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return map[string]interface{}{
		"goroutines":        runtime.NumGoroutine(),
		"heap_alloc_mb":     m.HeapAlloc / 1024 / 1024,
		"heap_sys_mb":       m.HeapSys / 1024 / 1024,
		"gc_cycles":         m.NumGC,
		"gc_pause_total_ns": m.PauseTotalNs,
	}
}
