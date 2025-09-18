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
	var statesMutex sync.RWMutex // Мьютекс для защиты userStates

	// Канал для обработки обновлений
	updateChan := make(chan telegram.Update, 100) // Буфер на 100 обновлений

	// Запускаем воркер для обработки обновлений
	go processUpdates(botUrl, repo, userStates, &statesMutex, updateChan)

	for {
		updates, err := telegram.GetUpdates(botUrl, offSet)
		if err != nil {
			log.Panicln("telegram.GetUpdates error: ", err)
			continue
		}

		for _, update := range updates {
			// Отправляем обновление в канал для обработки
			select {
			case updateChan <- update:
				log.Printf("Update %d queued for processing", update.UpdateId)
			default:
				// Канал переполнен, логируем предупреждение
				log.Printf("WARNING: Update channel full! Dropping update %d", update.UpdateId)
			}
			offSet = update.UpdateId + 1
		}
	}
}

// processUpdates обрабатывает обновления из канала
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

		// Безопасное чтение и обновление состояния
		statesMutex.Lock()
		if _, ok := userStates[chatId]; !ok {
			userStates[chatId] = states.SetStart()
			log.Printf("New user %d initialized with start state", chatId)
		}
		currentState := userStates[chatId]
		statesMutex.Unlock()

		// Обрабатываем обновление
		newState := respond(botUrl, update, currentState, repo)

		// Безопасно сохраняем новое состояние
		statesMutex.Lock()
		userStates[chatId] = newState
		statesMutex.Unlock()

		log.Printf("User %d state updated: %s", chatId, getStateName(newState.Type))
	}
}

func respond(botUrl string, update telegram.Update, state states.State, repo database.ContentRepositoryInterface) states.State {
	chatId := update.Message.Chat.ChatId

	if update.Message.Text == "/help" {
		return commands.Help(botUrl, chatId)
	}
	if update.Message.Text == "/start" {
		return commands.Start(botUrl, chatId)
	}
	if update.Message.Text == "/admin" {
		return commands.Admin(botUrl, chatId, repo)
	}

	switch state.Type {
	case states.StateAdminKeyboard, states.StateStartKeyboard:
		return handleCallback(botUrl, update.CallbackQuery, repo, state)

	// Тренер
	case states.StateEnterTrainerName:
		return commands.SetTrainerName(botUrl, chatId, update, repo, state)

	case states.StateEnterTrainerTgId:
		return commands.SetTrainerTgId(botUrl, chatId, update, repo, state)

	case states.StateEnterTrainerChatId:
		return commands.SetTrainerChatId(botUrl, chatId, update, repo, state)

	case states.StateEnterTrainerInfo:
		return commands.SetTrainerInfo(botUrl, chatId, update, repo, state)

	case states.StateConfirmTrainerCreation:
		// Обработка подтверждения будет в callback'ах
		return handleCallback(botUrl, update.CallbackQuery, repo, state)

	// Состояния редактирования тренеров
	case states.StateSelectTrainerToEdit:
		// Обработка выбора тренера будет в callback'ах
		return handleCallback(botUrl, update.CallbackQuery, repo, state)
	case states.StateEditTrainerName:
		return commands.SetEditTrainerName(botUrl, chatId, update, repo, state.GetID())

	case states.StateEditTrainerTgId:
		return commands.SetEditTrainerTgId(botUrl, chatId, update, repo, state.GetID())

	case states.StateEditTrainerInfo:
		return commands.SetEditTrainerInfo(botUrl, chatId, update, repo, state.GetID())

	case states.StateConfirmTrainerEdit:
		// Обработка будет в callback'ах
		return handleCallback(botUrl, update.CallbackQuery, repo, state)

	case states.StateConfirmTrainerDelete:
		// Обработка будет в callback'ах
		return handleCallback(botUrl, update.CallbackQuery, repo, state)

	// Состояния создания трасс
	case states.StateEnterTrackName:
		return commands.SetTrackName(botUrl, chatId, update, repo, state)

	case states.StateEnterTrackInfo:
		return commands.SetTrackInfo(botUrl, chatId, update, repo, state)

	case states.StateConfirmTrackCreation:
		// Обработка подтверждения будет в callback'ах
		return handleCallback(botUrl, update.CallbackQuery, repo, state)

	// Состояния редактирования трасс
	case states.StateSelectTrackToEdit:
		// Обработка выбора трассы будет в callback'ах
		return handleCallback(botUrl, update.CallbackQuery, repo, state)

	case states.StateEditTrackName:
		return commands.SetEditTrackName(botUrl, chatId, update, repo, state.GetID())

	case states.StateEditTrackInfo:
		return commands.SetEditTrackInfo(botUrl, chatId, update, repo, state.GetID())

	case states.StateConfirmTrackEdit:
		// Обработка будет в callback'ах
		return handleCallback(botUrl, update.CallbackQuery, repo, state)

	case states.StateConfirmTrackDelete:
		// Обработка будет в callback'ах
		return handleCallback(botUrl, update.CallbackQuery, repo, state)

	// Расписание

	// Регистрация пользователя
	case states.StateEnterUserName:
		return commands.SetUserName(botUrl, chatId, update, repo, state)

	case states.StateConfirmUserRegistration:
		// Обработка подтверждения будет в callback'ах
		return handleCallback(botUrl, update.CallbackQuery, repo, state)

	// Создание тренировки
	case states.StateEnterTrainingTrainer:
		// Обработка выбора тренера будет в callback'ах
		if update.CallbackQuery != nil {
			return handleCallback(botUrl, update.CallbackQuery, repo, state)
		}
		return states.SetError()

	case states.StateEnterTrainingTrack:
		// Обработка выбора трассы будет в callback'ах
		if update.CallbackQuery != nil {
			return handleCallback(botUrl, update.CallbackQuery, repo, state)
		}
		return states.SetError()

	case states.StateEnterTrainingDate:
		return commands.SetTrainingDate(botUrl, chatId, update, repo, state)

	case states.StateEnterTrainingMaxParticipants:
		return commands.SetTrainingMaxParticipants(botUrl, chatId, update, repo, state)

	case states.StateConfirmTrainingCreation:
		// Обработка подтверждения будет в callback'ах
		return handleCallback(botUrl, update.CallbackQuery, repo, state)

	// Регистрация на тренировку
	case states.StateSelectTrainingToRegister:
		// Обработка выбора тренировки будет в callback'ах
		return handleCallback(botUrl, update.CallbackQuery, repo, state)

	case states.StateConfirmTrainingRegistration:
		// Обработка подтверждения будет в callback'ах
		return handleCallback(botUrl, update.CallbackQuery, repo, state)

	// Новые состояния для пошаговой записи на тренировки
	case states.StateSelectTrackForRegistration:
		// Обработка выбора трассы будет в callback'ах
		return handleCallback(botUrl, update.CallbackQuery, repo, state)

	case states.StateSelectTrainerForRegistration:
		// Обработка выбора тренера будет в callback'ах
		return handleCallback(botUrl, update.CallbackQuery, repo, state)

	case states.StateSelectTrainingTimeForRegistration:
		// Обработка выбора времени тренировки будет в callback'ах
		return handleCallback(botUrl, update.CallbackQuery, repo, state)

	case states.StateStart:
		return commands.Start(botUrl, chatId)

	case states.StateError:
		log.Println("Error")
		return commands.Help(botUrl, chatId)
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

	// Обработка динамических callback'ов для редактирования тренеров
	if strings.HasPrefix(data, "selectTrainer_") {
		trainerIdStr := strings.TrimPrefix(data, "selectTrainer_")
		if trainerId, err := strconv.ParseUint(trainerIdStr, 10, 32); err == nil {
			return commands.SelectTrainerToEdit(botUrl, chatId, messageId, uint(trainerId), repo)
		}
	}

	if strings.HasPrefix(data, "editTrainerName_") {
		trainerIdStr := strings.TrimPrefix(data, "editTrainerName_")
		if trainerId, err := strconv.ParseUint(trainerIdStr, 10, 32); err == nil {
			return commands.EditTrainerName(botUrl, chatId, messageId, uint(trainerId))
		}
	}

	if strings.HasPrefix(data, "editTrainerTgId_") {
		trainerIdStr := strings.TrimPrefix(data, "editTrainerTgId_")
		if trainerId, err := strconv.ParseUint(trainerIdStr, 10, 32); err == nil {
			return commands.EditTrainerTgId(botUrl, chatId, messageId, uint(trainerId))
		}
	}

	if strings.HasPrefix(data, "editTrainerInfo_") {
		trainerIdStr := strings.TrimPrefix(data, "editTrainerInfo_")
		if trainerId, err := strconv.ParseUint(trainerIdStr, 10, 32); err == nil {
			return commands.EditTrainerInfo(botUrl, chatId, messageId, uint(trainerId))
		}
	}

	if strings.HasPrefix(data, "deleteTrainer_") {
		trainerIdStr := strings.TrimPrefix(data, "deleteTrainer_")
		if trainerId, err := strconv.ParseUint(trainerIdStr, 10, 32); err == nil {
			return commands.ConfirmTrainerDeletion(botUrl, chatId, messageId, uint(trainerId), repo)
		}
	}

	if strings.HasPrefix(data, "confirmDelete_") {
		trainerIdStr := strings.TrimPrefix(data, "confirmDelete_")
		if trainerId, err := strconv.ParseUint(trainerIdStr, 10, 32); err == nil {
			return commands.ExecuteTrainerDeletion(botUrl, chatId, messageId, uint(trainerId), repo)
		}
	}

	// Обработка динамических callback'ов для редактирования трасс
	if strings.HasPrefix(data, "selectTrack_") {
		trackIdStr := strings.TrimPrefix(data, "selectTrack_")
		if trackId, err := strconv.ParseUint(trackIdStr, 10, 32); err == nil {
			return commands.SelectTrackToEdit(botUrl, chatId, messageId, uint(trackId), repo)
		}
	}

	if strings.HasPrefix(data, "editTrackName_") {
		trackIdStr := strings.TrimPrefix(data, "editTrackName_")
		if trackId, err := strconv.ParseUint(trackIdStr, 10, 32); err == nil {
			return commands.EditTrackName(botUrl, chatId, messageId, uint(trackId))
		}
	}

	if strings.HasPrefix(data, "editTrackInfo_") {
		trackIdStr := strings.TrimPrefix(data, "editTrackInfo_")
		if trackId, err := strconv.ParseUint(trackIdStr, 10, 32); err == nil {
			return commands.EditTrackInfo(botUrl, chatId, messageId, uint(trackId))
		}
	}

	if strings.HasPrefix(data, "deleteTrack_") {
		trackIdStr := strings.TrimPrefix(data, "deleteTrack_")
		if trackId, err := strconv.ParseUint(trackIdStr, 10, 32); err == nil {
			return commands.ConfirmTrackDeletion(botUrl, chatId, messageId, uint(trackId), repo)
		}
	}

	if strings.HasPrefix(data, "confirmDeleteTrack_") {
		trackIdStr := strings.TrimPrefix(data, "confirmDeleteTrack_")
		if trackId, err := strconv.ParseUint(trackIdStr, 10, 32); err == nil {
			return commands.ExecuteTrackDeletion(botUrl, chatId, messageId, uint(trackId), repo)
		}
	}

	// Обработка динамических callback'ов для регистрации на тренировки
	if strings.HasPrefix(data, "selectTraining_") {
		trainingIdStr := strings.TrimPrefix(data, "selectTraining_")
		if trainingId, err := strconv.ParseUint(trainingIdStr, 10, 32); err == nil {
			return commands.ConfirmTrainingRegistration(botUrl, chatId, messageId, uint(trainingId), repo)
		}
	}

	if strings.HasPrefix(data, "confirmTrainingRegistration_") {
		trainingIdStr := strings.TrimPrefix(data, "confirmTrainingRegistration_")
		if trainingId, err := strconv.ParseUint(trainingIdStr, 10, 32); err == nil {
			return commands.ExecuteTrainingRegistration(botUrl, chatId, messageId, uint(trainingId), repo)
		}
	}

	// Обработка callback'ов для тренеров (подтверждение/отклонение заявок)
	if strings.HasPrefix(data, "approveRegistration_") {
		registrationIdStr := strings.TrimPrefix(data, "approveRegistration_")
		if registrationId, err := strconv.ParseUint(registrationIdStr, 10, 32); err == nil {
			return commands.ApproveTrainingRegistration(botUrl, chatId, messageId, uint(registrationId), repo)
		}
	}

	if strings.HasPrefix(data, "rejectRegistration_") {
		registrationIdStr := strings.TrimPrefix(data, "rejectRegistration_")
		if registrationId, err := strconv.ParseUint(registrationIdStr, 10, 32); err == nil {
			return commands.RejectTrainingRegistration(botUrl, chatId, messageId, uint(registrationId), repo)
		}
	}

	// Обработка callback'ов для создания тренировок
	if strings.HasPrefix(data, "selectTrainerForTraining_") {
		trainerIdStr := strings.TrimPrefix(data, "selectTrainerForTraining_")
		if trainerId, err := strconv.ParseUint(trainerIdStr, 10, 32); err == nil {
			return commands.SetTrainingTrainer(botUrl, chatId, messageId, uint(trainerId), repo, state)
		}
	}

	if strings.HasPrefix(data, "selectTrackForTraining_") {
		trackIdStr := strings.TrimPrefix(data, "selectTrackForTraining_")
		if trackId, err := strconv.ParseUint(trackIdStr, 10, 32); err == nil {
			return commands.SetTrainingTrack(botUrl, chatId, messageId, uint(trackId), repo, state)
		}
	}

	if strings.HasPrefix(data, "editTraining_") {
		trainingIdStr := strings.TrimPrefix(data, "editTraining_")
		if trainingId, err := strconv.ParseUint(trainingIdStr, 10, 32); err == nil {
			return commands.EditTraining(botUrl, chatId, messageId, uint(trainingId), repo)
		}
	}

	if strings.HasPrefix(data, "toggleTrainingStatus_") {
		trainingIdStr := strings.TrimPrefix(data, "toggleTrainingStatus_")
		if trainingId, err := strconv.ParseUint(trainingIdStr, 10, 32); err == nil {
			return commands.ToggleTrainingStatus(botUrl, chatId, messageId, uint(trainingId), repo)
		}
	}

	if strings.HasPrefix(data, "deleteTraining_") {
		trainingIdStr := strings.TrimPrefix(data, "deleteTraining_")
		if trainingId, err := strconv.ParseUint(trainingIdStr, 10, 32); err == nil {
			return commands.DeleteTraining(botUrl, chatId, messageId, uint(trainingId), repo)
		}
	}

	// Обработка callback'ов для пошаговой записи на тренировки
	if strings.HasPrefix(data, "selectTrackForRegistration_") {
		trackIdStr := strings.TrimPrefix(data, "selectTrackForRegistration_")
		if trackId, err := strconv.ParseUint(trackIdStr, 10, 32); err == nil {
			return commands.SelectTrackForRegistration(botUrl, chatId, messageId, uint(trackId), repo, state)
		}
	}

	if strings.HasPrefix(data, "selectTrainerForRegistration_") {
		trainerIdStr := strings.TrimPrefix(data, "selectTrainerForRegistration_")
		if trainerId, err := strconv.ParseUint(trainerIdStr, 10, 32); err == nil {
			return commands.SelectTrainerForRegistration(botUrl, chatId, messageId, uint(trainerId), repo, state)
		}
	}

	if strings.HasPrefix(data, "selectTrainingTimeForRegistration_") {
		trainingIdStr := strings.TrimPrefix(data, "selectTrainingTimeForRegistration_")
		if trainingId, err := strconv.ParseUint(trainingIdStr, 10, 32); err == nil {
			return commands.SelectTrainingTimeForRegistration(botUrl, chatId, messageId, uint(trainingId), repo, state)
		}
	}

	switch data {
	// Навигация
	case "start":
		return commands.ReturnToStart(botUrl, chatId, messageId)
	case "help":
		telegram.EditMessage(botUrl, chatId, messageId, "👋 <b>RVA Academy Bot</b>\n\n"+
			"📋 Команды:\n"+
			"/start - главное меню\n"+
			"/help - справка\n"+
			"/admin - админ-панель", telegram.CreateNavigationKeyboard())
		return states.SetStartKeyboard()
	case "admin":
		// Проверяем права администратора
		if !commands.IsAdmin(chatId, repo) {
			telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Доступ запрещен</b>\n"+
				"Нет прав администратора.", telegram.CreateBaseKeyboard())
			return states.SetStartKeyboard()
		}

		telegram.EditMessage(botUrl, chatId, messageId, "⚙️ <b>Админ-панель</b>\n"+
			"", telegram.CreateAdminKeyboard())
		return states.SetAdminKeyboard()

	// Админские меню
	case "trainersMenu":
		telegram.EditMessage(botUrl, chatId, messageId, "👨‍🏫 Управление тренерами\n\n"+
			"", telegram.CreateTrainersMenuKeyboard())
		return states.SetAdminKeyboard()

	case "tracksMenu":
		telegram.EditMessage(botUrl, chatId, messageId, "🏁 <b>Управление трассами</b>\n\n"+
			"🎯 <b>Доступные действия:</b>\n"+
			"➕ Добавление новых трасс\n"+
			"🏁 Просмотр списка трасс\n"+
			"✏️ Редактирование трасс\n"+
			"🗑️ Удаление трасс\n\n"+
			"", telegram.CreateTracksMenuKeyboard())
		return states.SetAdminKeyboard()

	case "scheduleMenu":
		telegram.EditMessage(botUrl, chatId, messageId, "📅 Управление расписанием\n\n"+
			"", telegram.CreateScheduleMenuKeyboard())
		return states.SetAdminKeyboard()

	// Тренеры
	case "createTrainer":
		return commands.CreateTrainer(botUrl, chatId, messageId)
	case "viewTrainers":
		return commands.ViewTrainers(botUrl, chatId, messageId, repo)
	case "editTrainer":
		return commands.EditTrainer(botUrl, chatId, messageId, repo)
	case "deleteTrainer":
		return commands.DeleteTrainer(botUrl, chatId, messageId, repo)

	// Трассы
	case "createTrack":
		return commands.CreateTrack(botUrl, chatId, messageId)
	case "viewTracks":
		return commands.ViewTracks(botUrl, chatId, messageId, repo)
	case "editTrack":
		return commands.EditTrack(botUrl, chatId, messageId, repo)
	case "deleteTrack":
		return commands.DeleteTrack(botUrl, chatId, messageId, repo)

	// Расписание
	case "createSchedule":
		return commands.CreateTraining(botUrl, chatId, messageId, repo)
	case "viewSchedule":
		return commands.ViewSchedule(botUrl, chatId, messageId, repo)
	case "editSchedule":
		return commands.EditSchedule(botUrl, chatId, messageId, repo)

	// Пользовательские функции
	case "BookTraining":
		return commands.StartTrainingRegistration(botUrl, chatId, messageId, repo)

	case "Info":
		return commands.Info(botUrl, chatId, messageId)

	case "infoTrainer":
		return commands.InfoTrainer(botUrl, chatId, messageId, repo)

	case "infoTrack":
		return commands.InfoTrack(botUrl, chatId, messageId, repo)

	case "viewScheduleUser":
		return commands.ViewScheduleUser(botUrl, chatId, messageId, repo)

	case "infoFormat":
		return commands.InfoFormat(botUrl, chatId, messageId)

	case "Raiting":
		return commands.ViewELORating(botUrl, chatId, messageId, repo)

	// Навигация назад при записи на тренировки
	case "backToTrackSelection":
		return commands.BackToTrackSelection(botUrl, chatId, messageId, repo, state)
	case "backToTrainerSelection":
		return commands.BackToTrainerSelection(botUrl, chatId, messageId, repo, state)

	// Подтверждение создания тренера
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

// getStateName возвращает читаемое имя состояния
// TODO: Move this to states package
func getStateName(stateType states.StateType) string {
	switch stateType {
	case states.StateStart:
		return "Start"
	case states.StateError:
		return "Error"
	case states.StateStartKeyboard:
		return "StartKeyboard"
	case states.StateAdminKeyboard:
		return "AdminKeyboard"
	case states.StateEnterTrainerName:
		return "EnterTrainerName"
	case states.StateEnterTrainerTgId:
		return "EnterTrainerTgId"
	case states.StateEnterTrainerChatId:
		return "EnterTrainerChatId"
	case states.StateEnterTrainerInfo:
		return "EnterTrainerInfo"
	case states.StateConfirmTrainerCreation:
		return "ConfirmTrainerCreation"
	case states.StateSelectTrainerToEdit:
		return "SelectTrainerToEdit"
	case states.StateEditTrainerName:
		return "EditTrainerName"
	case states.StateEditTrainerTgId:
		return "EditTrainerTgId"
	case states.StateEditTrainerInfo:
		return "EditTrainerInfo"
	case states.StateConfirmTrainerEdit:
		return "ConfirmTrainerEdit"
	case states.StateConfirmTrainerDelete:
		return "ConfirmTrainerDelete"
	case states.StateEnterTrackName:
		return "EnterTrackName"
	case states.StateEnterTrackInfo:
		return "EnterTrackInfo"
	case states.StateConfirmTrackCreation:
		return "ConfirmTrackCreation"
	case states.StateSelectTrackToEdit:
		return "SelectTrackToEdit"
	case states.StateEditTrackName:
		return "EditTrackName"
	case states.StateEditTrackInfo:
		return "EditTrackInfo"
	case states.StateConfirmTrackEdit:
		return "ConfirmTrackEdit"
	case states.StateConfirmTrackDelete:
		return "ConfirmTrackDelete"
	case states.StateEnterUserName:
		return "EnterUserName"
	case states.StateConfirmUserRegistration:
		return "ConfirmUserRegistration"
	case states.StateEnterTrainingTrainer:
		return "EnterTrainingTrainer"
	case states.StateEnterTrainingTrack:
		return "EnterTrainingTrack"
	case states.StateEnterTrainingDate:
		return "EnterTrainingDate"
	case states.StateEnterTrainingMaxParticipants:
		return "EnterTrainingMaxParticipants"
	case states.StateConfirmTrainingCreation:
		return "ConfirmTrainingCreation"
	case states.StateSelectTrainingToRegister:
		return "SelectTrainingToRegister"
	case states.StateConfirmTrainingRegistration:
		return "ConfirmTrainingRegistration"
	case states.StateSelectTrackForRegistration:
		return "SelectTrackForRegistration"
	case states.StateSelectTrainerForRegistration:
		return "SelectTrainerForRegistration"
	case states.StateSelectTrainingTimeForRegistration:
		return "SelectTrainingTimeForRegistration"
	default:
		return "Unknown"
	}
}
