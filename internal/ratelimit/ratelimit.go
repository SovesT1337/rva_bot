package ratelimit

import (
	"context"
	"sync"
	"time"

	"x.localhost/rvabot/internal/logger"
)

// Limiter интерфейс для rate limiting
type Limiter interface {
	Allow(ctx context.Context, key string) bool
	Wait(ctx context.Context, key string) error
}

// TokenBucket реализует алгоритм token bucket
type TokenBucket struct {
	capacity   int
	tokens     int
	refillRate time.Duration
	lastRefill time.Time
	mu         sync.Mutex
}

// NewTokenBucket создает новый token bucket
func NewTokenBucket(capacity int, refillRate time.Duration) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     capacity,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// Allow проверяет, можно ли выполнить запрос
func (tb *TokenBucket) Allow(ctx context.Context, key string) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)

	// Пополняем токены
	tokensToAdd := int(elapsed / tb.refillRate)
	if tokensToAdd > 0 {
		tb.tokens += tokensToAdd
		if tb.tokens > tb.capacity {
			tb.tokens = tb.capacity
		}
		tb.lastRefill = now
	}

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}

// Wait ждет, пока не появится доступный токен
func (tb *TokenBucket) Wait(ctx context.Context, key string) error {
	for {
		if tb.Allow(ctx, key) {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(tb.refillRate):
			// Продолжаем попытки
		}
	}
}

// UserRateLimiter управляет rate limiting для пользователей
type UserRateLimiter struct {
	buckets map[string]*TokenBucket
	mu      sync.RWMutex
	config  Config
}

// Config конфигурация rate limiter
type Config struct {
	Capacity   int           // Максимальное количество запросов
	RefillRate time.Duration // Интервал пополнения токенов
	CleanupAge time.Duration // Время жизни неактивных bucket'ов
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() Config {
	return Config{
		Capacity:   10,              // 10 запросов
		RefillRate: time.Second,     // 1 токен в секунду
		CleanupAge: 5 * time.Minute, // Очистка через 5 минут
	}
}

// NewUserRateLimiter создает новый rate limiter для пользователей
func NewUserRateLimiter(config Config) *UserRateLimiter {
	limiter := &UserRateLimiter{
		buckets: make(map[string]*TokenBucket),
		config:  config,
	}

	// Запускаем очистку неактивных bucket'ов
	go limiter.cleanup()

	return limiter
}

// Allow проверяет, может ли пользователь выполнить запрос
func (url *UserRateLimiter) Allow(ctx context.Context, userID string) bool {
	bucket := url.getOrCreateBucket(userID)
	allowed := bucket.Allow(ctx, userID)

	if !allowed {
		logger.Warn("RATE_LIMIT", "Rate limit exceeded for user %s", userID)
	}

	return allowed
}

// Wait ждет, пока пользователь сможет выполнить запрос
func (url *UserRateLimiter) Wait(ctx context.Context, userID string) error {
	bucket := url.getOrCreateBucket(userID)
	return bucket.Wait(ctx, userID)
}

// getOrCreateBucket получает или создает bucket для пользователя
func (url *UserRateLimiter) getOrCreateBucket(userID string) *TokenBucket {
	url.mu.RLock()
	bucket, exists := url.buckets[userID]
	url.mu.RUnlock()

	if exists {
		return bucket
	}

	url.mu.Lock()
	defer url.mu.Unlock()

	// Двойная проверка
	if bucket, exists := url.buckets[userID]; exists {
		return bucket
	}

	bucket = NewTokenBucket(url.config.Capacity, url.config.RefillRate)
	url.buckets[userID] = bucket

	return bucket
}

// cleanup периодически очищает неактивные bucket'ы
func (url *UserRateLimiter) cleanup() {
	ticker := time.NewTicker(url.config.CleanupAge)
	defer ticker.Stop()

	for range ticker.C {
		url.mu.Lock()
		now := time.Now()

		for userID, bucket := range url.buckets {
			bucket.mu.Lock()
			lastActivity := bucket.lastRefill
			bucket.mu.Unlock()

			if now.Sub(lastActivity) > url.config.CleanupAge {
				delete(url.buckets, userID)
				logger.BotInfo("Очищен rate limit bucket для пользователя %s", userID)
			}
		}

		url.mu.Unlock()
	}
}

// GetStats возвращает статистику rate limiter
func (url *UserRateLimiter) GetStats() map[string]interface{} {
	url.mu.RLock()
	defer url.mu.RUnlock()

	return map[string]interface{}{
		"active_buckets": len(url.buckets),
		"capacity":       url.config.Capacity,
		"refill_rate":    url.config.RefillRate.String(),
	}
}
