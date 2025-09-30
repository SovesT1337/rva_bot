package handler

import (
	"context"
	"strconv"
	"time"

	"x.localhost/rvabot/internal/commands"
	"x.localhost/rvabot/internal/database"
	"x.localhost/rvabot/internal/logger"
	"x.localhost/rvabot/internal/metrics"
	"x.localhost/rvabot/internal/ratelimit"
	"x.localhost/rvabot/internal/recovery"
	"x.localhost/rvabot/internal/state"
	"x.localhost/rvabot/internal/states"
	"x.localhost/rvabot/internal/telegram"
)

// UpdateProcessor обрабатывает обновления от Telegram
type UpdateProcessor struct {
	botUrl       string
	repo         database.ContentRepositoryInterface
	rateLimiter  *ratelimit.UserRateLimiter
	stateManager *state.Manager
	updateChan   chan telegram.Update
}

// NewUpdateProcessor создает новый процессор обновлений
func NewUpdateProcessor(botUrl string, repo database.ContentRepositoryInterface,
	rateLimiter *ratelimit.UserRateLimiter, stateManager *state.Manager) *UpdateProcessor {
	return &UpdateProcessor{
		botUrl:       botUrl,
		repo:         repo,
		rateLimiter:  rateLimiter,
		stateManager: stateManager,
		updateChan:   make(chan telegram.Update, 100),
	}
}

// Start запускает обработку обновлений
func (up *UpdateProcessor) Start() {
	go up.processUpdates()
}

// ProcessUpdate добавляет обновление в очередь обработки
func (up *UpdateProcessor) ProcessUpdate(update telegram.Update) {
	select {
	case up.updateChan <- update:
	default:
		logger.Warn("HANDLER", "Очередь обновлений переполнена! Пропускаем обновление %d", update.UpdateId)
	}
}

// processUpdates обрабатывает обновления из очереди
func (up *UpdateProcessor) processUpdates() {
	for update := range up.updateChan {
		// Обрабатываем каждое обновление с recovery
		recovery.RecoverFunc(context.Background(), "processUpdates", func() {
			up.handleUpdate(update)
		})
	}
}

// handleUpdate обрабатывает одно обновление
func (up *UpdateProcessor) handleUpdate(update telegram.Update) {
	chatId := up.getChatId(update)
	updateType := up.getUpdateType(update)

	// Увеличиваем счетчик общих обновлений
	metrics.IncrementTotalUpdates()

	logger.Info("HANDLER", "Обработка %s от чата %d", updateType, chatId)

	// Применяем rate limiting если доступен
	if up.rateLimiter != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if !up.rateLimiter.Allow(ctx, strconv.Itoa(chatId)) {
			logger.Warn("HANDLER", "Rate limit exceeded for user %d, skipping update", chatId)
			metrics.IncrementRateLimitedUsers()
			cancel()
			return
		}
		cancel()
	}

	// Получаем или создаем состояние пользователя
	currentState := up.stateManager.GetOrCreateState(chatId)

	// Обрабатываем обновление
	newState := up.respond(update, currentState)

	// Сохраняем новое состояние
	up.stateManager.SetState(chatId, newState)

	// Увеличиваем счетчик обработанных обновлений
	metrics.IncrementProcessedUpdates()

	logger.UserInfo(chatId, "Состояние: %s", newState.Type)
}

// getChatId извлекает chat ID из обновления
func (up *UpdateProcessor) getChatId(update telegram.Update) int {
	if update.Message.Chat.ChatId != 0 {
		return update.Message.Chat.ChatId
	}
	if update.CallbackQuery != nil {
		return update.CallbackQuery.Message.Chat.ChatId
	}
	return 0
}

// getUpdateType определяет тип обновления
func (up *UpdateProcessor) getUpdateType(update telegram.Update) string {
	if update.CallbackQuery != nil {
		return "callback"
	}
	return "message"
}

// isCommand проверяет, является ли текст командой
func isCommand(text string) bool {
	return text == "/help" || text == "/start" || text == "/admin"
}

// isTextInputState проверяет, является ли состояние состоянием ввода текста
func isTextInputState(stateType states.StateType) bool {
	textInputStates := map[states.StateType]bool{
		states.StateSetTrainerName:             true,
		states.StateSetTrainerTgId:             true,
		states.StateSetTrainerChatId:           true,
		states.StateSetTrainerInfo:             true,
		states.StateEditTrainerName:            true,
		states.StateEditTrainerTgId:            true,
		states.StateEditTrainerInfo:            true,
		states.StateSetTrackName:               true,
		states.StateSetTrackInfo:               true,
		states.StateEditTrackName:              true,
		states.StateEditTrackInfo:              true,
		states.StateSetUserName:                true,
		states.StateSetUserTgId:                true,
		states.StateSetTrainingStartTime:       true,
		states.StateSetTrainingEndTime:         true,
		states.StateSetTrainingMaxParticipants: true,
		states.StateSuggestTraining:            true,
	}
	return textInputStates[stateType]
}

// respond обрабатывает обновление и возвращает новое состояние
func (up *UpdateProcessor) respond(update telegram.Update, state states.State) states.State {
	var chatId int
	if update.Message.Chat.ChatId != 0 {
		chatId = update.Message.Chat.ChatId
	} else if update.CallbackQuery != nil {
		chatId = update.CallbackQuery.Message.Chat.ChatId
	}

	// Обрабатываем callback запросы
	if update.CallbackQuery != nil {
		callbackHandler := NewCallbackHandler(up.botUrl, up.repo)
		return callbackHandler.HandleCallback(update.CallbackQuery, state)
	}

	// Обрабатываем текстовые сообщения
	if update.Message.Text != "" {
		// Сначала проверяем, является ли это командой
		if isCommand(update.Message.Text) {
			return up.handleTextCommand(update, chatId)
		}

		// Если не команда, но есть состояние ввода текста - обрабатываем как ввод
		if isTextInputState(state.Type) {
			return up.handleTextInput(update, chatId, state)
		}

		// Если не команда и нет состояния ввода - показываем помощь
		return commands.Help(up.botUrl, chatId)
	}

	return states.SetStart()
}

// handleTextCommand обрабатывает текстовые команды
func (up *UpdateProcessor) handleTextCommand(update telegram.Update, chatId int) states.State {
	switch update.Message.Text {
	case "/help":
		return commands.Help(up.botUrl, chatId)
	case "/start":
		return commands.Start(up.botUrl, chatId, up.repo)
	case "/admin":
		return commands.Admin(up.botUrl, chatId, up.repo)
	default:
		// Неизвестная команда - показываем помощь
		return commands.Help(up.botUrl, chatId)
	}
}

// handleTextInput обрабатывает ввод текста в различных состояниях
func (up *UpdateProcessor) handleTextInput(update telegram.Update, chatId int, state states.State) states.State {
	textInputHandlers := map[states.StateType]func() states.State{
		states.StateSetTrainerName:   func() states.State { return commands.SetTrainerName(up.botUrl, chatId, update, up.repo, state) },
		states.StateSetTrainerTgId:   func() states.State { return commands.SetTrainerTgId(up.botUrl, chatId, update, up.repo, state) },
		states.StateSetTrainerChatId: func() states.State { return commands.SetTrainerChatId(up.botUrl, chatId, update, up.repo, state) },
		states.StateSetTrainerInfo:   func() states.State { return commands.SetTrainerInfo(up.botUrl, chatId, update, up.repo, state) },
		states.StateEditTrainerName: func() states.State {
			return commands.SetEditTrainerName(up.botUrl, chatId, update, up.repo, state.GetID())
		},
		states.StateEditTrainerTgId: func() states.State {
			return commands.SetEditTrainerTgId(up.botUrl, chatId, update, up.repo, state.GetID())
		},
		states.StateEditTrainerInfo: func() states.State {
			return commands.SetEditTrainerInfo(up.botUrl, chatId, update, up.repo, state.GetID())
		},
		states.StateSetTrackName: func() states.State { return commands.SetTrackName(up.botUrl, chatId, update, up.repo, state) },
		states.StateSetTrackInfo: func() states.State { return commands.SetTrackInfo(up.botUrl, chatId, update, up.repo, state) },
		states.StateEditTrackName: func() states.State {
			return commands.SetEditTrackName(up.botUrl, chatId, update, up.repo, state.GetID())
		},
		states.StateEditTrackInfo: func() states.State {
			return commands.SetEditTrackInfo(up.botUrl, chatId, update, up.repo, state.GetID())
		},
		states.StateSetUserName:          func() states.State { return commands.SetUserName(up.botUrl, chatId, update, up.repo, state) },
		states.StateSetUserTgId:          func() states.State { return commands.SetUserTgId(up.botUrl, chatId, update, up.repo, state) },
		states.StateSetTrainingStartTime: func() states.State { return commands.SetTrainingStartTime(up.botUrl, chatId, update, up.repo, state) },
		states.StateSetTrainingEndTime:   func() states.State { return commands.SetTrainingEndTime(up.botUrl, chatId, update, up.repo, state) },
		states.StateSetTrainingMaxParticipants: func() states.State {
			return commands.SetTrainingMaxParticipants(up.botUrl, chatId, update, up.repo, state)
		},
		states.StateSuggestTraining: func() states.State {
			return commands.ProcessTrainingSuggestion(up.botUrl, chatId, update, up.repo, state)
		},
		states.StateStart: func() states.State { return commands.Start(up.botUrl, chatId, up.repo) },
		states.StateError: func() states.State { return commands.Help(up.botUrl, chatId) },
	}

	if handler, ok := textInputHandlers[state.Type]; ok {
		logger.UserInfo(chatId, "Вызов обработчика текстового ввода для состояния: %s", state.Type)
		return handler()
	}

	logger.Warn("HANDLER", "Состояние не существует: %s", state)
	return states.SetError()
}
