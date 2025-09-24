package shutdown

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"x.localhost/rvabot/internal/logger"
)

// Manager управляет graceful shutdown
type Manager struct {
	shutdownChan chan os.Signal
	handlers     []ShutdownHandler
	mu           sync.RWMutex
	timeout      time.Duration
}

// ShutdownHandler интерфейс для обработчиков shutdown
type ShutdownHandler interface {
	Shutdown(ctx context.Context) error
	Name() string
}

// NewManager создает новый менеджер shutdown
func NewManager(timeout time.Duration) *Manager {
	return &Manager{
		shutdownChan: make(chan os.Signal, 1),
		handlers:     make([]ShutdownHandler, 0),
		timeout:      timeout,
	}
}

// RegisterHandler регистрирует обработчик shutdown
func (m *Manager) RegisterHandler(handler ShutdownHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handlers = append(m.handlers, handler)
	logger.BotInfo("Зарегистрирован shutdown handler: %s", handler.Name())
}

// Start начинает прослушивание сигналов
func (m *Manager) Start() {
	signal.Notify(m.shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-m.shutdownChan
		logger.BotInfo("Получен сигнал %v, начинаем graceful shutdown...", sig)
		m.shutdown()
	}()
}

// Wait блокирует до получения сигнала shutdown
func (m *Manager) Wait() {
	<-m.shutdownChan
}

// shutdown выполняет graceful shutdown всех зарегистрированных обработчиков
func (m *Manager) shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	m.mu.RLock()
	handlers := make([]ShutdownHandler, len(m.handlers))
	copy(handlers, m.handlers)
	m.mu.RUnlock()

	var wg sync.WaitGroup
	errorChan := make(chan error, len(handlers))

	for _, handler := range handlers {
		wg.Add(1)
		go func(h ShutdownHandler) {
			defer wg.Done()

			logger.BotInfo("Завершение работы %s...", h.Name())
			if err := h.Shutdown(ctx); err != nil {
				logger.BotError("Ошибка при завершении %s: %v", h.Name(), err)
				errorChan <- err
			} else {
				logger.BotInfo("%s успешно завершен", h.Name())
			}
		}(handler)
	}

	// Ждем завершения всех обработчиков или таймаута
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		logger.BotInfo("Все обработчики успешно завершены")
	case <-ctx.Done():
		logger.BotError("Таймаут при graceful shutdown: %v", ctx.Err())
	}

	close(errorChan)

	// Подсчитываем ошибки
	errorCount := 0
	for range errorChan {
		errorCount++
	}

	if errorCount > 0 {
		logger.BotError("Завершение с %d ошибками", errorCount)
		os.Exit(1)
	}

	logger.BotInfo("Graceful shutdown завершен успешно")
	os.Exit(0)
}
