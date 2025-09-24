package backoff

import (
	"math"
	"time"
)

// BackoffStrategy интерфейс для стратегий backoff
type BackoffStrategy interface {
	Next() time.Duration
	Reset()
}

// ExponentialBackoff реализует exponential backoff с jitter
type ExponentialBackoff struct {
	baseDelay   time.Duration
	maxDelay    time.Duration
	multiplier  float64
	attempts    int
	maxAttempts int
}

// NewExponentialBackoff создает новый exponential backoff
func NewExponentialBackoff(baseDelay, maxDelay time.Duration, multiplier float64, maxAttempts int) *ExponentialBackoff {
	return &ExponentialBackoff{
		baseDelay:   baseDelay,
		maxDelay:    maxDelay,
		multiplier:  multiplier,
		attempts:    0,
		maxAttempts: maxAttempts,
	}
}

// Next возвращает следующую задержку
func (eb *ExponentialBackoff) Next() time.Duration {
	if eb.attempts >= eb.maxAttempts {
		return eb.maxDelay
	}

	// Вычисляем задержку: baseDelay * multiplier^attempts
	delay := float64(eb.baseDelay) * math.Pow(eb.multiplier, float64(eb.attempts))

	// Добавляем jitter (±25%)
	jitter := delay * 0.25
	delay = delay + (jitter * (2*math.Mod(float64(time.Now().UnixNano()), 1.0) - 1))

	// Ограничиваем максимальной задержкой
	if delay > float64(eb.maxDelay) {
		delay = float64(eb.maxDelay)
	}

	eb.attempts++
	return time.Duration(delay)
}

// Reset сбрасывает счетчик попыток
func (eb *ExponentialBackoff) Reset() {
	eb.attempts = 0
}

// GetAttempts возвращает количество попыток
func (eb *ExponentialBackoff) GetAttempts() int {
	return eb.attempts
}

// IsMaxAttemptsReached проверяет, достигнуто ли максимальное количество попыток
func (eb *ExponentialBackoff) IsMaxAttemptsReached() bool {
	return eb.attempts >= eb.maxAttempts
}
