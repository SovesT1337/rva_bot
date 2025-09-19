package handler

import (
	"log"
	"strconv"
	"strings"
	"sync"

	"x.localhost/rvabot/internal/commands"
	"x.localhost/rvabot/internal/database"
	"x.localhost/rvabot/internal/states"
	"x.localhost/rvabot/internal/telegram"
)

func BotLoop(botUrl string, repo database.ContentRepositoryInterface) {
	offSet := 0

	userStates := make(map[int]states.State)
	var statesMutex sync.RWMutex

	updateChan := make(chan telegram.Update, 100)

	go processUpdates(botUrl, repo, userStates, &statesMutex, updateChan)

	for {
		updates, err := telegram.GetUpdates(botUrl, offSet)
		if err != nil {
			log.Panicln("telegram.GetUpdates error: ", err)
			continue
		}

		for _, update := range updates {
			select {
			case updateChan <- update:
				log.Printf("Update %d queued for processing", update.UpdateId)
			default:
				log.Printf("WARNING: Update channel full! Dropping update %d", update.UpdateId)
			}
			offSet = update.UpdateId + 1
		}
	}
}

func processUpdates(botUrl string, repo database.ContentRepositoryInterface,
	userStates map[int]states.State, statesMutex *sync.RWMutex, updateChan <-chan telegram.Update) {

	for update := range updateChan {
		chatId := update.Message.Chat.ChatId
		updateType := "message"
		if update.CallbackQuery != nil {
			chatId = update.CallbackQuery.Message.Chat.ChatId
			updateType = "callback"
		}

		log.Printf("Processing %s update %d from chat %d", updateType, update.UpdateId, chatId)

		statesMutex.Lock()
		if _, ok := userStates[chatId]; !ok {
			userStates[chatId] = states.SetStart()
			log.Printf("New user %d initialized with start state", chatId)
		}
		currentState := userStates[chatId]
		statesMutex.Unlock()

		newState := respond(botUrl, update, currentState, repo)

		statesMutex.Lock()
		userStates[chatId] = newState
		statesMutex.Unlock()

		log.Printf("User %d state updated: %s", chatId, newState.Type)
	}
}

func respond(botUrl string, update telegram.Update, state states.State, repo database.ContentRepositoryInterface) states.State {
	chatId := update.Message.Chat.ChatId

	switch update.Message.Text {
	case "/help":
		return commands.Help(botUrl, chatId)
	case "/start":
		return commands.Start(botUrl, chatId)
	case "/admin":
		return commands.Admin(botUrl, chatId, repo)
	}

	handlers := map[states.StateType]func() states.State{
		states.StateAdminKeyboard: func() states.State { return handleCallback(botUrl, update.CallbackQuery, repo, state) },
		states.StateStartKeyboard: func() states.State { return handleCallback(botUrl, update.CallbackQuery, repo, state) },

		states.StateEnterTrainerName:       func() states.State { return commands.SetTrainerName(botUrl, chatId, update, repo, state) },
		states.StateEnterTrainerTgId:       func() states.State { return commands.SetTrainerTgId(botUrl, chatId, update, repo, state) },
		states.StateEnterTrainerChatId:     func() states.State { return commands.SetTrainerChatId(botUrl, chatId, update, repo, state) },
		states.StateEnterTrainerInfo:       func() states.State { return commands.SetTrainerInfo(botUrl, chatId, update, repo, state) },
		states.StateConfirmTrainerCreation: func() states.State { return handleCallback(botUrl, update.CallbackQuery, repo, state) },
		states.StateSelectTrainerToEdit:    func() states.State { return handleCallback(botUrl, update.CallbackQuery, repo, state) },
		states.StateEditTrainerName:        func() states.State { return commands.SetEditTrainerName(botUrl, chatId, update, repo, state.GetID()) },
		states.StateEditTrainerTgId:        func() states.State { return commands.SetEditTrainerTgId(botUrl, chatId, update, repo, state.GetID()) },
		states.StateEditTrainerInfo:        func() states.State { return commands.SetEditTrainerInfo(botUrl, chatId, update, repo, state.GetID()) },
		states.StateConfirmTrainerEdit:     func() states.State { return handleCallback(botUrl, update.CallbackQuery, repo, state) },
		states.StateConfirmTrainerDelete:   func() states.State { return handleCallback(botUrl, update.CallbackQuery, repo, state) },

		states.StateEnterTrackName:       func() states.State { return commands.SetTrackName(botUrl, chatId, update, repo, state) },
		states.StateEnterTrackInfo:       func() states.State { return commands.SetTrackInfo(botUrl, chatId, update, repo, state) },
		states.StateConfirmTrackCreation: func() states.State { return handleCallback(botUrl, update.CallbackQuery, repo, state) },
		states.StateSelectTrackToEdit:    func() states.State { return handleCallback(botUrl, update.CallbackQuery, repo, state) },
		states.StateEditTrackName:        func() states.State { return commands.SetEditTrackName(botUrl, chatId, update, repo, state.GetID()) },
		states.StateEditTrackInfo:        func() states.State { return commands.SetEditTrackInfo(botUrl, chatId, update, repo, state.GetID()) },
		states.StateConfirmTrackEdit:     func() states.State { return handleCallback(botUrl, update.CallbackQuery, repo, state) },
		states.StateConfirmTrackDelete:   func() states.State { return handleCallback(botUrl, update.CallbackQuery, repo, state) },

		states.StateEnterUserName:           func() states.State { return commands.SetUserName(botUrl, chatId, update, repo, state) },
		states.StateConfirmUserRegistration: func() states.State { return handleCallback(botUrl, update.CallbackQuery, repo, state) },

		states.StateEnterTrainingTrainer:              func() states.State { return handleCallback(botUrl, update.CallbackQuery, repo, state) },
		states.StateEnterTrainingTrack:                func() states.State { return handleCallback(botUrl, update.CallbackQuery, repo, state) },
		states.StateEnterTrainingDate:                 func() states.State { return handleCallback(botUrl, update.CallbackQuery, repo, state) },
		states.StateEnterTrainingMaxParticipants:      func() states.State { return handleCallback(botUrl, update.CallbackQuery, repo, state) },
		states.StateConfirmTrainingCreation:           func() states.State { return handleCallback(botUrl, update.CallbackQuery, repo, state) },
		states.StateSelectTrainingToRegister:          func() states.State { return handleCallback(botUrl, update.CallbackQuery, repo, state) },
		states.StateConfirmTrainingRegistration:       func() states.State { return handleCallback(botUrl, update.CallbackQuery, repo, state) },
		states.StateSelectTrackForRegistration:        func() states.State { return handleCallback(botUrl, update.CallbackQuery, repo, state) },
		states.StateSelectTrainerForRegistration:      func() states.State { return handleCallback(botUrl, update.CallbackQuery, repo, state) },
		states.StateSelectTrainingTimeForRegistration: func() states.State { return handleCallback(botUrl, update.CallbackQuery, repo, state) },

		states.StateStart: func() states.State { return commands.Start(botUrl, chatId) },
		states.StateError: func() states.State { return commands.Help(botUrl, chatId) },
	}

	if handler, ok := handlers[state.Type]; ok {
		return handler()
	}

	log.Println("State doestn exist: ", state)
	return states.SetError()
}

func handleCallback(botUrl string, query *telegram.CallbackQuery, repo database.ContentRepositoryInterface, state states.State) states.State {
	if query == nil {
		return states.SetError()
	}

	chatId := query.Message.Chat.ChatId
	messageId := query.Message.MessageId
	data := query.Data

	telegram.AnswerCallbackQuery(botUrl, query.ID)

	log.Printf("Callback from user %d: %s", chatId, data)
	prefix := ""
	id_str := ""
	if idx := strings.Index(data, "_"); idx != -1 {
		prefix = data[:idx]
		id_str = data[idx+1:]
	}

	id := 0

	if id_str != "" {
		parsedId, err := strconv.ParseUint(id_str, 10, 32)
		if err != nil {
			log.Printf("Error parsing id from user %d: %s", chatId, err)
			return states.SetError()
		}
		id = int(parsedId)
	}
	log.Printf("prefix from user %d: %s", chatId, prefix)

	handlers := map[string]func() states.State{
		"selectTrainer":      func() states.State { return commands.SelectTrainerToEdit(botUrl, chatId, messageId, uint(id), repo) },
		"editTrainerName":    func() states.State { return commands.EditTrainerName(botUrl, chatId, messageId, uint(id)) },
		"editTrainerTgId":    func() states.State { return commands.EditTrainerTgId(botUrl, chatId, messageId, uint(id)) },
		"editTrainerInfo":    func() states.State { return commands.EditTrainerInfo(botUrl, chatId, messageId, uint(id)) },
		"deleteTrainer":      func() states.State { return commands.ConfirmTrainerDeletion(botUrl, chatId, messageId, uint(id), repo) },
		"confirmDelete":      func() states.State { return commands.ExecuteTrainerDeletion(botUrl, chatId, messageId, uint(id), repo) },
		"selectTrack":        func() states.State { return commands.SelectTrackToEdit(botUrl, chatId, messageId, uint(id), repo) },
		"editTrackName":      func() states.State { return commands.EditTrackName(botUrl, chatId, messageId, uint(id)) },
		"editTrackInfo":      func() states.State { return commands.EditTrackInfo(botUrl, chatId, messageId, uint(id)) },
		"deleteTrack":        func() states.State { return commands.ConfirmTrackDeletion(botUrl, chatId, messageId, uint(id), repo) },
		"confirmDeleteTrack": func() states.State { return commands.ExecuteTrackDeletion(botUrl, chatId, messageId, uint(id), repo) },

		"selectTraining": func() states.State {
			return commands.ConfirmTrainingRegistration(botUrl, chatId, messageId, uint(id), repo)
		},
		"confirmTrainingRegistration": func() states.State {
			return commands.ExecuteTrainingRegistration(botUrl, chatId, messageId, uint(id), repo)
		},
		"approveRegistration": func() states.State {
			return commands.ApproveTrainingRegistration(botUrl, chatId, messageId, uint(id), repo)
		},
		"rejectRegistration": func() states.State {
			return commands.RejectTrainingRegistration(botUrl, chatId, messageId, uint(id), repo)
		},
		"selectTrainerForTraining": func() states.State {
			return commands.SetTrainingTrainer(botUrl, chatId, messageId, uint(id), repo, state)
		},
		"selectTrackForTraining": func() states.State {
			return commands.SetTrainingTrack(botUrl, chatId, messageId, uint(id), repo, state)
		},

		"editTraining":         func() states.State { return commands.EditTraining(botUrl, chatId, messageId, uint(id), repo) },
		"toggleTrainingStatus": func() states.State { return commands.ToggleTrainingStatus(botUrl, chatId, messageId, uint(id), repo) },
		"deleteTraining":       func() states.State { return commands.DeleteTraining(botUrl, chatId, messageId, uint(id), repo) },
		"selectTrackForRegistration": func() states.State {
			return commands.SelectTrackForRegistration(botUrl, chatId, messageId, uint(id), repo, state)
		},
		"selectTrainerForRegistration": func() states.State {
			return commands.SelectTrainerForRegistration(botUrl, chatId, messageId, uint(id), repo, state)
		},
		"selectTrainingTimeForRegistration": func() states.State {
			return commands.SelectTrainingTimeForRegistration(botUrl, chatId, messageId, uint(id), repo, state)
		},
	}

	if handler, ok := handlers[prefix]; ok {
		return handler()
	}

	handlersMap := map[string]func() states.State{
		"start": func() states.State {
			return commands.ReturnToStart(botUrl, chatId, messageId)
		},
		"help": func() states.State {
			telegram.EditMessage(botUrl, chatId, messageId, "👋 <b>RVA Academy Bot</b>\n\n"+
				"📋 Команды:\n"+
				"/start - главное меню\n"+
				"/help - справка\n"+
				"/admin - админ-панель", telegram.CreateNavigationKeyboard())
			return states.SetStartKeyboard()
		},
		"admin": func() states.State {
			if !commands.IsAdmin(chatId, repo) {
				telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Доступ запрещен</b>\n"+
					"Нет прав администратора.", telegram.CreateBaseKeyboard())
				return states.SetStartKeyboard()
			}
			telegram.EditMessage(botUrl, chatId, messageId, "⚙️ <b>Админ-панель</b>\n"+
				"", telegram.CreateAdminKeyboard())
			return states.SetAdminKeyboard()
		},
		"trainersMenu": func() states.State {
			telegram.EditMessage(botUrl, chatId, messageId, "👨‍🏫 Управление тренерами\n\n"+
				"", telegram.CreateTrainersMenuKeyboard())
			return states.SetAdminKeyboard()
		},
		"tracksMenu": func() states.State {
			telegram.EditMessage(botUrl, chatId, messageId, "🏁 <b>Управление трассами</b>\n\n"+
				"🎯 <b>Доступные действия:</b>\n"+
				"➕ Добавление новых трасс\n"+
				"🏁 Просмотр списка трасс\n"+
				"✏️ Редактирование трасс\n"+
				"🗑️ Удаление трасс\n\n"+
				"", telegram.CreateTracksMenuKeyboard())
			return states.SetAdminKeyboard()
		},
		"scheduleMenu": func() states.State {
			telegram.EditMessage(botUrl, chatId, messageId, "📅 Управление расписанием\n\n"+
				"", telegram.CreateScheduleMenuKeyboard())
			return states.SetAdminKeyboard()
		},
		"createTrainer": func() states.State {
			return commands.CreateTrainer(botUrl, chatId, messageId)
		},
		"viewTrainers": func() states.State {
			return commands.ViewTrainers(botUrl, chatId, messageId, repo)
		},
		"editTrainer": func() states.State {
			return commands.EditTrainer(botUrl, chatId, messageId, repo)
		},
		"deleteTrainer": func() states.State {
			return commands.DeleteTrainer(botUrl, chatId, messageId, repo)
		},
		"createTrack": func() states.State {
			return commands.CreateTrack(botUrl, chatId, messageId)
		},
		"viewTracks": func() states.State {
			return commands.ViewTracks(botUrl, chatId, messageId, repo)
		},
		"editTrack": func() states.State {
			return commands.EditTrack(botUrl, chatId, messageId, repo)
		},
		"deleteTrack": func() states.State {
			return commands.DeleteTrack(botUrl, chatId, messageId, repo)
		},
		"createSchedule": func() states.State {
			return commands.CreateTraining(botUrl, chatId, messageId, repo)
		},
		"viewSchedule": func() states.State {
			return commands.ViewSchedule(botUrl, chatId, messageId, repo)
		},
		"editSchedule": func() states.State {
			return commands.EditSchedule(botUrl, chatId, messageId, repo)
		},
		"BookTraining": func() states.State {
			return commands.StartTrainingRegistration(botUrl, chatId, messageId, repo)
		},
		"Info": func() states.State {
			return commands.Info(botUrl, chatId, messageId)
		},
		"infoTrainer": func() states.State {
			return commands.InfoTrainer(botUrl, chatId, messageId, repo)
		},
		"infoTrack": func() states.State {
			return commands.InfoTrack(botUrl, chatId, messageId, repo)
		},
		"viewScheduleUser": func() states.State {
			return commands.ViewScheduleUser(botUrl, chatId, messageId, repo)
		},
		"infoFormat": func() states.State {
			return commands.InfoFormat(botUrl, chatId, messageId)
		},
		"backToTrackSelection": func() states.State {
			return commands.BackToTrackSelection(botUrl, chatId, messageId, repo, state)
		},
		"backToTrainerSelection": func() states.State {
			return commands.BackToTrainerSelection(botUrl, chatId, messageId, repo, state)
		},
	}

	if handler, ok := handlersMap[data]; ok {
		return handler()
	}

	switch data {
	case "confirm":
		switch state.Type {
		case states.StateConfirmTrainerCreation:
			tempData := state.GetTempTrainerData()
			if tempData.Name != "" && tempData.TgId != "" && tempData.Info != "" {
				return commands.ConfirmTrainerCreation(botUrl, chatId, messageId, repo, tempData)
			}
		case states.StateConfirmTrackCreation:
			tempData := state.GetTempTrackData()
			if tempData.Name != "" && tempData.Info != "" {
				return commands.ConfirmTrackCreation(botUrl, chatId, messageId, repo, tempData)
			}
		case states.StateConfirmUserRegistration:
			tempData := state.GetTempUserData()
			if tempData.Name != "" {
				return commands.ConfirmUserRegistration(botUrl, chatId, messageId, repo, tempData)
			}
		case states.StateConfirmTrainingCreation:
			tempData := state.GetTempTrainingData()
			if tempData.TrainerID != 0 && tempData.TrackID != 0 && tempData.Date != "" {
				return commands.ConfirmTrainingCreation(botUrl, chatId, messageId, repo, tempData)
			}
		case states.StateConfirmTrainingRegistration:
			trainingId := state.Data["trainingId"].(uint)
			return commands.ExecuteTrainingRegistration(botUrl, chatId, messageId, uint(trainingId), repo)
		}
		return states.SetError()

	case "cancel":
		switch state.Type {
		case states.StateConfirmTrainerCreation:
			return commands.CancelTrainerCreation(botUrl, chatId, messageId)
		case states.StateConfirmTrackCreation:
			return commands.CancelTrackCreation(botUrl, chatId, messageId)
		case states.StateConfirmUserRegistration:
			telegram.EditMessage(botUrl, chatId, messageId, "🚫 <b>Регистрация отменена</b>\n\n"+
				"💡 Вы можете зарегистрироваться позже.", telegram.CreateBaseKeyboard())
			return states.SetStartKeyboard()
		case states.StateConfirmTrainingCreation:
			telegram.EditMessage(botUrl, chatId, messageId, "🚫 <b>Создание тренировки отменено</b>\n\n"+
				"💡 Вы можете создать тренировку позже.", telegram.CreateBackToScheduleMenuKeyboard())
			return states.SetAdminKeyboard()
		case states.StateEnterTrainerName, states.StateEnterTrainerTgId, states.StateEnterTrainerChatId, states.StateEnterTrainerInfo:
			telegram.EditMessage(botUrl, chatId, messageId, "🚫 <b>Создание тренера отменено</b>\n\n"+
				"💡 Вы можете создать тренера позже.", telegram.CreateBackToTrainersMenuKeyboard())
			return states.SetAdminKeyboard()
		case states.StateEnterTrackName, states.StateEnterTrackInfo:
			telegram.EditMessage(botUrl, chatId, messageId, "🚫 <b>Создание трассы отменено</b>\n\n"+
				"💡 Вы можете создать трассу позже.", telegram.CreateBackToTracksMenuKeyboard())
			return states.SetAdminKeyboard()
		case states.StateEnterTrainingTrainer, states.StateEnterTrainingTrack, states.StateEnterTrainingDate, states.StateEnterTrainingMaxParticipants:
			telegram.EditMessage(botUrl, chatId, messageId, "🚫 <b>Создание тренировки отменено</b>\n\n"+
				"💡 Вы можете создать тренировку позже.", telegram.CreateBackToScheduleMenuKeyboard())
			return states.SetAdminKeyboard()
		case states.StateEnterUserName:
			telegram.EditMessage(botUrl, chatId, messageId, "🚫 <b>Регистрация отменена</b>\n\n"+
				"💡 Вы можете зарегистрироваться позже.", telegram.CreateBaseKeyboard())
			return states.SetStartKeyboard()
		case states.StateConfirmTrainingRegistration:
			telegram.EditMessage(botUrl, chatId, messageId, "🚫 <b>Запись на тренировку отменена</b>\n\n"+
				"💡 Вы можете записаться на тренировку позже.", telegram.CreateBaseKeyboard())
			return states.SetStartKeyboard()
		default:
			telegram.EditMessage(botUrl, chatId, messageId, "🚫 <b>Операция отменена</b>\n\n"+
				"", telegram.CreateBaseKeyboard())
			return states.SetStartKeyboard()
		}

	}

	return states.SetStart()
}
