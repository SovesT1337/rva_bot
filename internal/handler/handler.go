package handler

import (
	"time"

	"x.localhost/rvabot/internal/backoff"
	"x.localhost/rvabot/internal/database"
	"x.localhost/rvabot/internal/errors"
	"x.localhost/rvabot/internal/logger"
	"x.localhost/rvabot/internal/ratelimit"
	"x.localhost/rvabot/internal/state"
	"x.localhost/rvabot/internal/telegram"
)

func BotLoop(botUrl string, repo database.ContentRepositoryInterface) {
	BotLoopWithRateLimit(botUrl, repo, nil)
}

func BotLoopWithRateLimit(botUrl string, repo database.ContentRepositoryInterface, rateLimiter *ratelimit.UserRateLimiter) {
	BotLoopWithComponents(botUrl, repo, rateLimiter, nil, nil)
}

func BotLoopWithComponents(botUrl string, repo database.ContentRepositoryInterface, rateLimiter *ratelimit.UserRateLimiter, stateManager *state.Manager, backoffStrategy backoff.BackoffStrategy) {
	offSet := 0

	// Создаем компоненты если они не переданы
	if stateManager == nil {
		stateManager = state.NewManager(30*time.Minute, 5*time.Minute) // 30 мин TTL, очистка каждые 5 мин
	}

	if backoffStrategy == nil {
		backoffStrategy = backoff.NewExponentialBackoff(
			1*time.Second,  // базовая задержка
			60*time.Second, // максимальная задержка
			2.0,            // множитель
			10,             // максимальное количество попыток
		)
	}

	// Создаем процессор обновлений
	updateProcessor := NewUpdateProcessor(botUrl, repo, rateLimiter, stateManager)
	updateProcessor.Start()

	for {
		updates, err := telegram.GetUpdates(botUrl, offSet)
		if err != nil {
			appErr := errors.WrapError(err, errors.ErrorTypeTelegram, "Ошибка получения обновлений")
			logger.BotError("Ошибка при получении обновлений: %v", appErr)

			// Используем правильный exponential backoff
			backoffDuration := backoffStrategy.Next()
			logger.BotInfo("Повторная попытка через %v...", backoffDuration)

			time.Sleep(backoffDuration)

			// Сбрасываем backoff при успешном получении обновлений
			if len(updates) > 0 {
				backoffStrategy.Reset()
			}
			continue
		}

		// Сбрасываем backoff при успешном получении обновлений
		backoffStrategy.Reset()

		for _, update := range updates {
			updateProcessor.ProcessUpdate(update)
			offSet = update.UpdateId + 1
		}
	}
}
