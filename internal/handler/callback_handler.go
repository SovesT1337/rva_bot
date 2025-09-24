package handler

import (
	"strconv"
	"strings"

	"x.localhost/rvabot/internal/commands"
	"x.localhost/rvabot/internal/database"
	"x.localhost/rvabot/internal/logger"
	"x.localhost/rvabot/internal/states"
	"x.localhost/rvabot/internal/telegram"
)

// CallbackHandler обрабатывает callback запросы
type CallbackHandler struct {
	botUrl string
	repo   database.ContentRepositoryInterface
}

// NewCallbackHandler создает новый обработчик callback'ов
func NewCallbackHandler(botUrl string, repo database.ContentRepositoryInterface) *CallbackHandler {
	return &CallbackHandler{
		botUrl: botUrl,
		repo:   repo,
	}
}

// HandleCallback обрабатывает callback запрос
func (ch *CallbackHandler) HandleCallback(query *telegram.CallbackQuery, state states.State) states.State {
	if query == nil {
		return states.SetError()
	}

	chatId := query.Message.Chat.ChatId
	messageId := query.Message.MessageId
	data := query.Data

	telegram.AnswerCallbackQuery(ch.botUrl, query.ID)

	logger.UserInfo(chatId, "Callback %s", data)

	prefix, id := ch.parseCallbackData(data, chatId)
	if prefix == "" && id == -1 {
		return states.SetError()
	}

	// Обрабатываем специальные случаи
	switch data {
	case "confirm":
		return ch.handleConfirmAction(chatId, messageId, state)
	case "cancel":
		return ch.handleCancelAction(chatId, messageId, state)
	}

	// Обрабатываем callback'и с префиксами
	if prefix != "" {
		return ch.handlePrefixedCallback(prefix, id, chatId, messageId, state)
	}

	// Обрабатываем простые callback'и
	return ch.handleSimpleCallback(data, chatId, messageId, state)
}

// parseCallbackData парсит данные callback'а
func (ch *CallbackHandler) parseCallbackData(data string, chatId int) (string, int) {
	prefix := ""
	id_str := ""
	if idx := strings.Index(data, "_"); idx != -1 {
		prefix = data[:idx]
		id_str = data[idx+1:]
	}

	if id_str == "" {
		return prefix, 0
	}

	parsedId, err := strconv.ParseUint(id_str, 10, 32)
	if err != nil {
		logger.UserError(chatId, "Ошибка парсинга ID: %s", err)
		return "", -1
	}

	logger.UserInfo(chatId, "Префикс %s", prefix)
	return prefix, int(parsedId)
}

// handlePrefixedCallback обрабатывает callback'и с префиксами
func (ch *CallbackHandler) handlePrefixedCallback(prefix string, id int, chatId, messageId int, state states.State) states.State {
	callbackHandlers := map[string]func() states.State{
		"editTrainerName": func() states.State { return commands.EditTrainerName(ch.botUrl, chatId, messageId, uint(id)) },
		"editTrainerTgId": func() states.State { return commands.EditTrainerTgId(ch.botUrl, chatId, messageId, uint(id)) },
		"editTrainerInfo": func() states.State { return commands.EditTrainerInfo(ch.botUrl, chatId, messageId, uint(id)) },
		"deleteTrainer": func() states.State {
			return commands.ConfirmTrainerDeletion(ch.botUrl, chatId, messageId, uint(id), ch.repo)
		},
		"confirmDelete": func() states.State {
			return commands.ExecuteTrainerDeletion(ch.botUrl, chatId, messageId, uint(id), ch.repo)
		},
		"editTrackName": func() states.State { return commands.EditTrackName(ch.botUrl, chatId, messageId, uint(id)) },
		"editTrackInfo": func() states.State { return commands.EditTrackInfo(ch.botUrl, chatId, messageId, uint(id)) },
		"deleteTrack": func() states.State {
			return commands.ConfirmTrackDeletion(ch.botUrl, chatId, messageId, uint(id), ch.repo)
		},
		"confirmDeleteTrack": func() states.State {
			return commands.ExecuteTrackDeletion(ch.botUrl, chatId, messageId, uint(id), ch.repo)
		},
		"selectTraining": func() states.State {
			return commands.ConfirmTrainingRegistration(ch.botUrl, chatId, messageId, uint(id), ch.repo)
		},
		"confirmTrainingRegistration": func() states.State {
			return commands.ExecuteTrainingRegistration(ch.botUrl, chatId, messageId, uint(id), ch.repo)
		},
		"approveRegistration": func() states.State {
			return commands.ApproveTrainingRegistration(ch.botUrl, chatId, messageId, uint(id), ch.repo)
		},
		"rejectRegistration": func() states.State {
			return commands.RejectTrainingRegistration(ch.botUrl, chatId, messageId, uint(id), ch.repo)
		},
		"selectTrainerForTraining": func() states.State {
			return commands.SetTrainingTrainer(ch.botUrl, chatId, messageId, uint(id), ch.repo, state)
		},
		"selectTrackForTraining": func() states.State {
			return commands.SetTrainingTrack(ch.botUrl, chatId, messageId, uint(id), ch.repo, state)
		},
		"editTraining": func() states.State { return commands.EditTraining(ch.botUrl, chatId, messageId, uint(id), ch.repo) },
		"toggleTrainingStatus": func() states.State {
			return commands.ToggleTrainingStatus(ch.botUrl, chatId, messageId, uint(id), ch.repo)
		},
		"deleteTraining": func() states.State {
			return commands.ConfirmTrainingDeletion(ch.botUrl, chatId, messageId, uint(id), ch.repo)
		},
		"confirmDeleteTraining": func() states.State {
			return commands.ExecuteTrainingDeletion(ch.botUrl, chatId, messageId, uint(id), ch.repo)
		},
		"selectTrackForRegistration": func() states.State {
			return commands.SelectTrackForRegistration(ch.botUrl, chatId, messageId, uint(id), ch.repo, state)
		},
		"selectTrainerForRegistration": func() states.State {
			return commands.SelectTrainerForRegistration(ch.botUrl, chatId, messageId, uint(id), ch.repo, state)
		},
		"selectTrainingTimeForRegistration": func() states.State {
			return commands.SelectTrainingTimeForRegistration(ch.botUrl, chatId, messageId, uint(id), ch.repo, state)
		},
		"markRequestReviewed": func() states.State {
			return commands.MarkTrainingRequestAsReviewed(ch.botUrl, chatId, messageId, uint(id), ch.repo)
		},
	}

	if handler, ok := callbackHandlers[prefix]; ok {
		return handler()
	}

	return states.SetStart()
}

// handleSimpleCallback обрабатывает простые callback'и
func (ch *CallbackHandler) handleSimpleCallback(data string, chatId, messageId int, state states.State) states.State {
	simpleCallbackHandlers := map[string]func() states.State{
		"start": func() states.State { return commands.ReturnToStart(ch.botUrl, chatId, messageId, ch.repo) },
		"help":  func() states.State { return commands.SendHelpMessage(ch.botUrl, chatId, messageId) },
		"admin": func() states.State {
			if !database.IsAdmin(chatId, ch.repo) {
				return commands.SendAccessDeniedMessage(ch.botUrl, chatId, messageId)
			}
			return commands.SendAdminPanelMessage(ch.botUrl, chatId, messageId)
		},
		"trainersMenu":     func() states.State { return commands.SendTrainersMenuMessage(ch.botUrl, chatId, messageId, ch.repo) },
		"tracksMenu":       func() states.State { return commands.SendTracksMenuMessage(ch.botUrl, chatId, messageId, ch.repo) },
		"scheduleMenu":     func() states.State { return commands.SendScheduleMenuMessage(ch.botUrl, chatId, messageId, ch.repo) },
		"createTrainer":    func() states.State { return commands.CreateTrainer(ch.botUrl, chatId, messageId) },
		"viewTrainers":     func() states.State { return commands.ViewTrainers(ch.botUrl, chatId, messageId, ch.repo) },
		"createTrack":      func() states.State { return commands.CreateTrack(ch.botUrl, chatId, messageId) },
		"viewTracks":       func() states.State { return commands.ViewTracks(ch.botUrl, chatId, messageId, ch.repo) },
		"createSchedule":   func() states.State { return commands.CreateTraining(ch.botUrl, chatId, messageId, ch.repo) },
		"viewSchedule":     func() states.State { return commands.ViewSchedule(ch.botUrl, chatId, messageId, ch.repo) },
		"editSchedule":     func() states.State { return commands.EditSchedule(ch.botUrl, chatId, messageId, ch.repo) },
		"BookTraining":     func() states.State { return commands.StartTrainingRegistration(ch.botUrl, chatId, messageId, ch.repo) },
		"Info":             func() states.State { return commands.Info(ch.botUrl, chatId, messageId) },
		"infoTrainer":      func() states.State { return commands.InfoTrainer(ch.botUrl, chatId, messageId, ch.repo) },
		"infoTrack":        func() states.State { return commands.InfoTrack(ch.botUrl, chatId, messageId, ch.repo) },
		"viewScheduleUser": func() states.State { return commands.ViewScheduleUser(ch.botUrl, chatId, messageId, ch.repo) },
		"infoFormat":       func() states.State { return commands.InfoFormat(ch.botUrl, chatId, messageId) },
		"suggestTraining":  func() states.State { return commands.SuggestTraining(ch.botUrl, chatId, messageId, ch.repo) },
		"trainingRequests": func() states.State { return commands.ViewTrainingRequests(ch.botUrl, chatId, messageId, ch.repo) },
		"backToTrackSelection": func() states.State {
			return commands.BackToTrackSelection(ch.botUrl, chatId, messageId, ch.repo, state)
		},
		"backToTrainerSelection": func() states.State {
			return commands.BackToTrainerSelection(ch.botUrl, chatId, messageId, ch.repo, state)
		},
	}

	if handler, ok := simpleCallbackHandlers[data]; ok {
		return handler()
	}

	return states.SetStart()
}

// handleConfirmAction обрабатывает подтверждение действий
func (ch *CallbackHandler) handleConfirmAction(chatId, messageId int, state states.State) states.State {
	switch state.Type {
	case states.StateConfirmTrainerCreation:
		tempData := state.GetTempTrainerData()
		if tempData.Name != "" && tempData.TgId != "" && tempData.Info != "" {
			return commands.ConfirmTrainerCreation(ch.botUrl, chatId, messageId, ch.repo, tempData)
		}
	case states.StateConfirmTrackCreation:
		tempData := state.GetTempTrackData()
		if tempData.Name != "" && tempData.Info != "" {
			return commands.ConfirmTrackCreation(ch.botUrl, chatId, messageId, ch.repo, tempData)
		}
	case states.StateConfirmUserRegistration:
		tempData := state.GetTempUserData()
		if tempData.Name != "" {
			return commands.ConfirmUserRegistration(ch.botUrl, chatId, messageId, ch.repo, tempData)
		}
	case states.StateConfirmTrainingCreation:
		tempData := state.GetTempTrainingData()
		if tempData.TrainerID != 0 && tempData.TrackID != 0 && tempData.StartTime != "" && tempData.EndTime != "" {
			return commands.ConfirmTrainingCreation(ch.botUrl, chatId, messageId, ch.repo, tempData)
		}
	case states.StateConfirmTrainingRegistration:
		if trainingId, ok := state.Data["trainingId"].(uint); ok {
			return commands.ExecuteTrainingRegistration(ch.botUrl, chatId, messageId, uint(trainingId), ch.repo)
		}
		logger.UserError(chatId, "Неверный тип trainingId в состоянии")
		return states.SetError()
	}
	return states.SetError()
}

// handleCancelAction обрабатывает отмену действий
func (ch *CallbackHandler) handleCancelAction(chatId, messageId int, state states.State) states.State {
	cancelHandlers := map[states.StateType]func() states.State{
		states.StateConfirmTrainerCreation: func() states.State {
			return commands.CancelTrainerCreation(ch.botUrl, chatId, messageId)
		},
		states.StateEditTrainerName: func() states.State {
			return commands.SendOperationCancelledWithTrainersMenu(ch.botUrl, chatId, messageId)
		},
		states.StateEditTrainerTgId: func() states.State {
			return commands.SendOperationCancelledWithTrainersMenu(ch.botUrl, chatId, messageId)
		},
		states.StateEditTrainerInfo: func() states.State {
			return commands.SendOperationCancelledWithTrainersMenu(ch.botUrl, chatId, messageId)
		},
		states.StateConfirmTrackCreation: func() states.State {
			return commands.CancelTrackCreation(ch.botUrl, chatId, messageId)
		},
		states.StateEditTrackName: func() states.State {
			return commands.SendOperationCancelledWithTracksMenu(ch.botUrl, chatId, messageId)
		},
		states.StateEditTrackInfo: func() states.State {
			return commands.SendOperationCancelledWithTracksMenu(ch.botUrl, chatId, messageId)
		},
		states.StateConfirmTrainingCreation: func() states.State {
			return commands.SendOperationCancelledWithScheduleMenu(ch.botUrl, chatId, messageId)
		},
		states.StateConfirmUserRegistration: func() states.State {
			return commands.SendOperationCancelledMessage(ch.botUrl, chatId, messageId)
		},
		states.StateConfirmTrainingRegistration: func() states.State {
			return commands.SendOperationCancelledMessage(ch.botUrl, chatId, messageId)
		},
		states.StateConfirmTrainingDelete: func() states.State {
			return commands.SendOperationCancelledWithScheduleMenu(ch.botUrl, chatId, messageId)
		},
	}

	if handler, ok := cancelHandlers[state.Type]; ok {
		return handler()
	}

	return commands.SendOperationCancelledMessage(ch.botUrl, chatId, messageId)
}
