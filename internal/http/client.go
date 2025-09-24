package http

import (
	"context"
	"net/http"
	"sync"
	"time"
)

// ClientPool управляет пулом HTTP клиентов
type ClientPool struct {
	clients map[string]*http.Client
	mutex   sync.RWMutex
	config  ClientConfig
}

// ClientConfig конфигурация HTTP клиента
type ClientConfig struct {
	Timeout         time.Duration
	MaxIdleConns    int
	MaxConnsPerHost int
	IdleConnTimeout time.Duration
}

// DefaultClientConfig возвращает конфигурацию по умолчанию
func DefaultClientConfig() ClientConfig {
	return ClientConfig{
		Timeout:         30 * time.Second,
		MaxIdleConns:    100,
		MaxConnsPerHost: 10,
		IdleConnTimeout: 90 * time.Second,
	}
}

// NewClientPool создает новый пул HTTP клиентов
func NewClientPool(config ClientConfig) *ClientPool {
	return &ClientPool{
		clients: make(map[string]*http.Client),
		config:  config,
	}
}

// GetClient возвращает HTTP клиент для указанного ключа
func (cp *ClientPool) GetClient(key string) *http.Client {
	cp.mutex.RLock()
	client, exists := cp.clients[key]
	cp.mutex.RUnlock()

	if exists {
		return client
	}

	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	// Двойная проверка
	if client, exists := cp.clients[key]; exists {
		return client
	}

	// Создаем новый клиент
	client = &http.Client{
		Timeout: cp.config.Timeout,
		Transport: &http.Transport{
			MaxIdleConns:        cp.config.MaxIdleConns,
			MaxConnsPerHost:     cp.config.MaxConnsPerHost,
			IdleConnTimeout:     cp.config.IdleConnTimeout,
			DisableKeepAlives:   false,
			DisableCompression:  false,
			MaxIdleConnsPerHost: cp.config.MaxConnsPerHost,
		},
	}

	cp.clients[key] = client
	return client
}

// GetDefaultClient возвращает клиент по умолчанию
func (cp *ClientPool) GetDefaultClient() *http.Client {
	return cp.GetClient("default")
}

// Close закрывает все клиенты в пуле
func (cp *ClientPool) Close() {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	for key, client := range cp.clients {
		// Закрываем transport для освобождения соединений
		if transport, ok := client.Transport.(*http.Transport); ok {
			transport.CloseIdleConnections()
		}
		delete(cp.clients, key)
	}
}

// GetStats возвращает статистику пула клиентов
func (cp *ClientPool) GetStats() map[string]interface{} {
	cp.mutex.RLock()
	defer cp.mutex.RUnlock()

	return map[string]interface{}{
		"total_clients":      len(cp.clients),
		"timeout":            cp.config.Timeout.String(),
		"max_idle_conns":     cp.config.MaxIdleConns,
		"max_conns_per_host": cp.config.MaxConnsPerHost,
	}
}

// RequestOptions опции для HTTP запроса
type RequestOptions struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    []byte
	Timeout time.Duration
}

// DoRequest выполняет HTTP запрос с использованием пула клиентов
func (cp *ClientPool) DoRequest(ctx context.Context, opts RequestOptions) (*http.Response, error) {
	client := cp.GetDefaultClient()

	// Создаем запрос
	req, err := http.NewRequestWithContext(ctx, opts.Method, opts.URL, nil)
	if err != nil {
		return nil, err
	}

	// Добавляем заголовки
	for key, value := range opts.Headers {
		req.Header.Set(key, value)
	}

	// Выполняем запрос
	return client.Do(req)
}

// Глобальный пул клиентов
var globalClientPool *ClientPool
var once sync.Once

// GetGlobalClientPool возвращает глобальный пул HTTP клиентов
func GetGlobalClientPool() *ClientPool {
	once.Do(func() {
		globalClientPool = NewClientPool(DefaultClientConfig())
	})
	return globalClientPool
}

// CloseGlobalClientPool закрывает глобальный пул клиентов
func CloseGlobalClientPool() {
	if globalClientPool != nil {
		globalClientPool.Close()
	}
}
