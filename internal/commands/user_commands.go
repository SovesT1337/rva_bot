package commands

import (
	"fmt"
	"log"
	"strings"

	"x.localhost/rvabot/internal/database"
	"x.localhost/rvabot/internal/states"
	"x.localhost/rvabot/internal/telegram"
)

func sendErrorMessage(botUrl string, chatId int, messageId int) states.State {
	telegram.EditMessage(botUrl, chatId, messageId, "❌ Произошла ошибка, повторите попытку позже", telegram.CreateStartKeyboard())
	return states.SetStartKeyboard()
}

func Help(botUrl string, ChatId int) states.State {
	telegram.SendMessage(botUrl, ChatId, "🎓 <b>Добро пожаловать в RVA Academy Bot!</b>\n\n"+
		"🤖 Я помогу вам управлять тренировками и тренерами.\n\n"+
		"📋 <b>Доступные команды:</b>\n"+
		"🏠 /start - главное меню\n"+
		"❓ /help - эта справка\n"+
		"⚙️ /admin - панель администратора\n\n"+
		"💡 <i>Используйте кнопки ниже для навигации</i>", telegram.CreateNavigationKeyboard())
	return states.SetStartKeyboard()
}

func Start(botUrl string, chatId int) states.State {
	telegram.SendMessage(botUrl, chatId, "🎯 <b>RVA Academy Bot</b>\n\n"+
		"🏃‍♂️ Добро пожаловать в систему регистрации на тренировки!\n\n", telegram.CreateStartKeyboard())
	return states.SetStartKeyboard()
}

func ReturnToStart(botUrl string, chatId int, messageId int) states.State {
	telegram.EditMessage(botUrl, chatId, messageId, "🏁 Добро пожаловать в RVA Academy!\n\n"+
		"", telegram.CreateStartKeyboard())
	return states.SetStartKeyboard()
}

func Info(botUrl string, chatId int, messageId int) states.State {
	telegram.EditMessage(botUrl, chatId, messageId, "ℹ️ Информация о RVA Academy\n\n"+
		"", telegram.CreateInfoKeyboard())
	return states.SetStartKeyboard()
}

func InfoTrainer(botUrl string, chatId int, messageId int, repo database.ContentRepositoryInterface) states.State {
	trainers, err := repo.GetTrainers()
	if err != nil {
		return sendErrorMessage(botUrl, chatId, messageId)
	}

	message := formatTrainersListForUsers(trainers)
	telegram.EditMessage(botUrl, chatId, messageId, message, telegram.CreateBackToInfoKeyboard())
	return states.SetStartKeyboard()
}

func InfoTrack(botUrl string, chatId int, messageId int, repo database.ContentRepositoryInterface) states.State {
	tracks, err := repo.GetTracks()
	if err != nil {
		return sendErrorMessage(botUrl, chatId, messageId)
	}

	message := formatTracksListForUsers(tracks)
	telegram.EditMessage(botUrl, chatId, messageId, message, telegram.CreateBackToInfoKeyboard())
	return states.SetStartKeyboard()
}

func InfoFormat(botUrl string, chatId int, messageId int) states.State {
	message := "Тут пока пусто, потом будет информация о формате тренировок"

	telegram.EditMessage(botUrl, chatId, messageId, message, telegram.CreateBackToInfoKeyboard())
	return states.SetStartKeyboard()
}

func ViewScheduleUser(botUrl string, chatId int, messageId int, repo database.ContentRepositoryInterface) states.State {
	trainings, err := repo.GetActiveTrainings()
	if err != nil {
		log.Printf("ERROR: Failed to get active trainings: %v", err)
		return sendErrorMessage(botUrl, chatId, messageId)
	}

	if len(trainings) == 0 {
		telegram.EditMessage(botUrl, chatId, messageId, "📅 <b>Расписание тренировок</b>\n\n"+
			"📝 <b>Активных тренировок пока нет</b>\n\n"+
			"💡 Следите за обновлениями! Новые тренировки появятся в ближайшее время.", telegram.CreateBackToInfoKeyboard())
		return states.SetStartKeyboard()
	}

	message := formatTrainingsListForUsers(trainings, repo)
	telegram.EditMessage(botUrl, chatId, messageId, message, telegram.CreateBackToInfoKeyboard())
	return states.SetStartKeyboard()
}

func SetUserName(botUrl string, chatId int, update telegram.Update, repo database.ContentRepositoryInterface, state states.State) states.State {
	name := update.Message.Text

	tempData := state.GetTempUserData()
	tempData.Name = name

	message := fmt.Sprintf("✅ <b>Подтверждение регистрации</b>\n\n"+
		"📋 <b>Проверьте данные:</b>\n\n"+
		"👤 <b>ФИО:</b> %s\n\n"+
		"❓ <b>Зарегистрироваться с этими данными?</b>", tempData.Name)

	telegram.SendMessage(botUrl, chatId, message, telegram.CreateConfirmationKeyboard())

	newState := states.SetConfirmUserRegistration()
	return newState.SetTempUserData(tempData)
}

func ConfirmUserRegistration(botUrl string, chatId int, messageId int, repo database.ContentRepositoryInterface, tempData *states.TempUserData) states.State {
	log.Printf("Creating user: %s (ChatId: %d)", tempData.Name, chatId)

	user := &database.User{
		Name:     tempData.Name,
		ChatId:   chatId,
		IsActive: true,
	}

	id, err := repo.CreateUser(user)
	if err != nil {
		log.Printf("ERROR: Failed to create user %s: %v", tempData.Name, err)
		return sendErrorMessage(botUrl, chatId, messageId)
	}

	log.Printf("User created successfully: %s (ID: %d)", tempData.Name, id)
	telegram.EditMessage(botUrl, chatId, messageId, "🎉 <b>Регистрация завершена!</b>\n"+
		"Добро пожаловать, "+tempData.Name+"!", telegram.CreateStartKeyboard())
	return states.SetStartKeyboard()
}

func StartTrainingRegistration(botUrl string, chatId int, messageId int, repo database.ContentRepositoryInterface) states.State {
	user, err := repo.GetUserByChatId(chatId)
	if err != nil || user == nil {
		telegram.EditMessage(botUrl, chatId, messageId, "🏃‍♂️ <b>Запись на тренировку</b>\n\n"+
			"📝 Введите ваше ФИО\n"+
			"<i>Пример: Иванов Иван Иванович</i>", telegram.CreateCancelKeyboard())
		return states.SetEnterUserName()
	}

	tracks, err := repo.GetTracksWithActiveTrainings()
	if err != nil {
		log.Printf("ERROR: Failed to get tracks with active trainings: %v", err)
		return sendErrorMessage(botUrl, chatId, messageId)
	}

	if len(tracks) == 0 {
		telegram.EditMessage(botUrl, chatId, messageId, "🏁 <b>Нет доступных трасс</b>\n"+
			"Нет активных тренировок.", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	telegram.EditMessage(botUrl, chatId, messageId, "🏃‍♂️ <b>Запись на тренировку</b>\n\n"+
		"👤 "+user.Name+"\n"+
		"🏁 <b>Шаг 1/3:</b> Трасса", telegram.CreateTrackSelectionForRegistrationKeyboard(tracks))

	tempData := &states.TempRegistrationData{}
	state := states.SetSelectTrackForRegistration()
	return state.SetTempRegistrationData(tempData)
}

func ConfirmTrainingRegistration(botUrl string, chatId int, messageId int, trainingId uint, repo database.ContentRepositoryInterface) states.State {
	training, err := repo.GetTrainingById(trainingId)
	if err != nil || training == nil {
		return sendErrorMessage(botUrl, chatId, messageId)
	}

	user, err := repo.GetUserByChatId(chatId)
	if err != nil || user == nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Пользователь не найден</b>\n\n"+
			"🔍 Сначала зарегистрируйтесь в системе.", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	existingRegistration, _ := repo.GetTrainingRegistrationByUserAndTraining(user.ID, trainingId)
	if existingRegistration != nil {
		telegram.EditMessage(botUrl, chatId, messageId, "⚠️ <b>Вы уже зарегистрированы</b>\n\n"+
			"🏃‍♂️ Вы уже записаны на эту тренировку.\n"+
			"📊 <b>Статус:</b> "+existingRegistration.Status, telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	registrations, _ := repo.GetTrainingRegistrationsByTrainingID(trainingId)
	registeredCount := 0
	for _, reg := range registrations {
		if reg.Status == "confirmed" || reg.Status == "pending" {
			registeredCount++
		}
	}

	if registeredCount >= training.MaxParticipants {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Нет свободных мест</b>\n\n"+
			"🏃‍♂️ На эту тренировку уже записалось максимальное количество участников.\n"+
			"💡 Попробуйте выбрать другую тренировку.", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	trainer, _ := repo.GetTrainerByID(training.TrainerID)
	track, _ := repo.GetTrackByID(training.TrackID)

	trainerName := "Неизвестный тренер"
	trackName := "Неизвестная трасса"

	if trainer != nil {
		trainerName = trainer.Name
	}
	if track != nil {
		trackName = track.Name
	}

	message := fmt.Sprintf("✅ <b>Подтверждение записи на тренировку</b>\n\n"+
		"📋 <b>Детали тренировки:</b>\n\n"+
		"🏃‍♂️ <b>Тренировка:</b> %s\n"+
		"👨‍🏫 <b>Тренер:</b> %s\n"+
		"📅 <b>Дата и время:</b> %s\n"+
		"👥 <b>Свободных мест:</b> %d\n\n"+
		"❓ <b>Подтвердить запись на тренировку?</b>",
		trackName, trainerName, training.Time.Format("02.01.2006 15:04"), training.MaxParticipants-registeredCount)

	telegram.EditMessage(botUrl, chatId, messageId, message, telegram.CreateTrainingRegistrationConfirmationKeyboard(trainingId))
	return states.SetConfirmTrainingRegistration(trainingId)
}

func ExecuteTrainingRegistration(botUrl string, chatId int, messageId int, trainingId uint, repo database.ContentRepositoryInterface) states.State {
	user, err := repo.GetUserByChatId(chatId)
	if err != nil || user == nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Пользователь не найден</b>\n\n"+
			"🔍 Сначала зарегистрируйтесь в системе.", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	registration := &database.TrainingRegistration{
		TrainingID: trainingId,
		UserID:     user.ID,
		Status:     "pending",
	}

	regId, err := repo.CreateTrainingRegistration(registration)
	if err != nil {
		log.Printf("ERROR: Failed to create training registration: %v", err)
		return sendErrorMessage(botUrl, chatId, messageId)
	}

	training, _ := repo.GetTrainingById(trainingId)
	trainer, _ := repo.GetTrainerByID(training.TrainerID)
	track, _ := repo.GetTrackByID(training.TrackID)

	if trainer != nil && trainer.ChatId != 0 {
		trackName := "Неизвестная трасса"
		if track != nil {
			trackName = track.Name
		}

		notificationMessage := fmt.Sprintf("🔔 <b>Новая заявка</b>\n"+
			"👤 %s\n"+
			"🏃‍♂️ %s\n"+
			"📅 %s",
			user.Name, trackName, training.Time.Format("02.01.2006 15:04"))

		telegram.SendMessage(botUrl, trainer.ChatId, notificationMessage, telegram.CreateTrainingApprovalKeyboard(regId))
	}

	log.Printf("Training registration created successfully: ID=%d, UserID=%d, TrainingID=%d", regId, user.ID, trainingId)
	telegram.EditMessage(botUrl, chatId, messageId, "🎉 <b>Заявка на тренировку отправлена!</b>\n\n"+
		"✅ <b>Ваша заявка принята и отправлена тренеру на рассмотрение.</b>\n\n"+
		"📱 <b>Вы получите уведомление о решении тренера.</b>\n"+
		"⏰ <b>Обычно рассмотрение занимает несколько часов.</b>", telegram.CreateBaseKeyboard())
	return states.SetStartKeyboard()
}

func BackToTrackSelection(botUrl string, chatId int, messageId int, repo database.ContentRepositoryInterface, state states.State) states.State {
	user, err := repo.GetUserByChatId(chatId)
	if err != nil || user == nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Пользователь не найден</b>", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	tracks, err := repo.GetTracksWithActiveTrainings()
	if err != nil {
		log.Printf("ERROR: Failed to get tracks with active trainings: %v", err)
		return sendErrorMessage(botUrl, chatId, messageId)
	}

	if len(tracks) == 0 {
		telegram.EditMessage(botUrl, chatId, messageId, "🏁 <b>Нет доступных трасс</b>\n"+
			"Нет активных тренировок.", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	telegram.EditMessage(botUrl, chatId, messageId, "🏃‍♂️ <b>Запись на тренировку</b>\n\n"+
		"👤 "+user.Name+"\n"+
		"🏁 <b>Шаг 1/3:</b> Трасса", telegram.CreateTrackSelectionForRegistrationKeyboard(tracks))

	tempData := &states.TempRegistrationData{}
	newState := states.SetSelectTrackForRegistration()
	return newState.SetTempRegistrationData(tempData)
}

func BackToTrainerSelection(botUrl string, chatId int, messageId int, repo database.ContentRepositoryInterface, state states.State) states.State {
	tempData := state.GetTempRegistrationData()
	if tempData.TrackID == 0 {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Ошибка навигации</b>\n"+
			"Начните заново.", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	track, err := repo.GetTrackByID(tempData.TrackID)
	if err != nil || track == nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Трасса не найдена</b>", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	trainers, err := repo.GetTrainersByTrack(tempData.TrackID)
	if err != nil {
		log.Printf("ERROR: Failed to get trainers by track: %v", err)
		return sendErrorMessage(botUrl, chatId, messageId)
	}

	if len(trainers) == 0 {
		telegram.EditMessage(botUrl, chatId, messageId, "👨‍🏫 <b>Нет тренеров</b>\n"+
			"На трассе \""+track.Name+"\" нет тренировок.", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	telegram.EditMessage(botUrl, chatId, messageId, "🏃‍♂️ <b>Запись на тренировку</b>\n\n"+
		"✅ Трасса: "+track.Name+"\n"+
		"👨‍🏫 <b>Шаг 2/3:</b> Тренер", telegram.CreateTrainerSelectionForRegistrationKeyboard(trainers))

	newState := states.SetSelectTrainerForRegistration()
	return newState.SetTempRegistrationData(tempData)
}

func SelectTrackForRegistration(botUrl string, chatId int, messageId int, trackId uint, repo database.ContentRepositoryInterface, state states.State) states.State {
	tempData := state.GetTempRegistrationData()
	tempData.TrackID = trackId

	track, err := repo.GetTrackByID(trackId)
	if err != nil || track == nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Трасса не найдена</b>", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	trainers, err := repo.GetTrainersByTrack(trackId)
	if err != nil {
		log.Printf("ERROR: Failed to get trainers by track: %v", err)
		return sendErrorMessage(botUrl, chatId, messageId)
	}

	if len(trainers) == 0 {
		telegram.EditMessage(botUrl, chatId, messageId, "👨‍🏫 <b>Нет тренеров</b>\n"+
			"На трассе \""+track.Name+"\" нет тренировок.", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	telegram.EditMessage(botUrl, chatId, messageId, "🏃‍♂️ <b>Запись на тренировку</b>\n\n"+
		"✅ Трасса: "+track.Name+"\n"+
		"👨‍🏫 <b>Шаг 2/3:</b> Тренер", telegram.CreateTrainerSelectionForRegistrationKeyboard(trainers))

	newState := states.SetSelectTrainerForRegistration()
	return newState.SetTempRegistrationData(tempData)
}

func SelectTrainerForRegistration(botUrl string, chatId int, messageId int, trainerId uint, repo database.ContentRepositoryInterface, state states.State) states.State {
	tempData := state.GetTempRegistrationData()
	tempData.TrainerID = trainerId

	trainer, err := repo.GetTrainerByID(trainerId)
	if err != nil || trainer == nil {
		return sendErrorMessage(botUrl, chatId, messageId)
	}

	track, err := repo.GetTrackByID(tempData.TrackID)
	if err != nil || track == nil {
		return sendErrorMessage(botUrl, chatId, messageId)
	}

	trainings, err := repo.GetActiveTrainingsByTrackAndTrainer(tempData.TrackID, trainerId)
	if err != nil {
		log.Printf("ERROR: Failed to get trainings by track and trainer: %v", err)
		return sendErrorMessage(botUrl, chatId, messageId)
	}

	if len(trainings) == 0 {
		telegram.EditMessage(botUrl, chatId, messageId, "📅 <b>Нет доступных тренировок</b>\n\n"+
			"🏃‍♂️ <b>Тренер:</b> "+trainer.Name+"\n"+
			"🏁 <b>Трасса:</b> "+track.Name+"\n\n"+
			"📝 <b>У выбранного тренера нет активных тренировок на этой трассе.</b>\n"+
			"💡 Попробуйте выбрать другого тренера или трассу.", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	for i := 0; i < len(trainings)-1; i++ {
		for j := 0; j < len(trainings)-i-1; j++ {
			if trainings[j].Time.After(trainings[j+1].Time) {
				trainings[j], trainings[j+1] = trainings[j+1], trainings[j]
			}
		}
	}

	telegram.EditMessage(botUrl, chatId, messageId, "🏃‍♂️ <b>Запись на тренировку</b>\n\n"+
		"✅ Трасса: "+track.Name+"\n"+
		"✅ Тренер: "+trainer.Name+"\n"+
		"📅 <b>Шаг 3/3:</b> Время", telegram.CreateTrainingTimeSelectionKeyboard(trainings))

	newState := states.SetSelectTrainingTimeForRegistration()
	return newState.SetTempRegistrationData(tempData)
}

func SelectTrainingTimeForRegistration(botUrl string, chatId int, messageId int, trainingId uint, repo database.ContentRepositoryInterface, state states.State) states.State {
	training, err := repo.GetTrainingById(trainingId)
	if err != nil || training == nil {
		return sendErrorMessage(botUrl, chatId, messageId)
	}

	tempData := state.GetTempRegistrationData()
	if training.TrackID != tempData.TrackID || training.TrainerID != tempData.TrainerID {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Ошибка выбора тренировки</b>\n\n"+
			"🔍 Выбранная тренировка не соответствует выбранным трассе и тренеру.", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	return ConfirmTrainingRegistration(botUrl, chatId, messageId, trainingId, repo)
}

func ApproveTrainingRegistration(botUrl string, chatId int, messageId int, registrationId uint, repo database.ContentRepositoryInterface) states.State {
	registration, err := repo.GetTrainingRegistrationByID(registrationId)
	if err != nil || registration == nil {
		return sendErrorMessage(botUrl, chatId, messageId)
	}

	training, _ := repo.GetTrainingById(registration.TrainingID)
	if training == nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Тренировка не найдена</b>", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	trainer, _ := repo.GetTrainerByID(training.TrainerID)
	if trainer == nil || trainer.ChatId != chatId {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Нет прав</b>", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	registration.Status = "confirmed"
	err = repo.UpdateTrainingRegistration(registrationId, registration)
	if err != nil {
		log.Printf("ERROR: Failed to approve training registration %d: %v", registrationId, err)
		return sendErrorMessage(botUrl, chatId, messageId)
	}

	user, _ := repo.GetUserByID(registration.UserID)
	track, _ := repo.GetTrackByID(training.TrackID)

	trackName := "Неизвестная трасса"
	if track != nil {
		trackName = track.Name
	}

	if user != nil {
		userMessage := fmt.Sprintf("🎉 <b>Заявка на тренировку одобрена!</b>\n\n"+
			"✅ <b>Ваша заявка на тренировку была подтверждена тренером.</b>\n\n"+
			"🏃‍♂️ <b>Тренировка:</b> %s\n"+
			"📅 <b>Дата и время:</b> %s\n\n"+
			"💡 <b>До встречи на тренировке!</b>",
			trackName, training.Time.Format("02.01.2006 15:04"))

		telegram.SendMessage(botUrl, user.ChatId, userMessage, telegram.CreateBaseKeyboard())
	}

	log.Printf("Training registration %d approved by trainer %d", registrationId, chatId)
	telegram.EditMessage(botUrl, chatId, messageId, "✅ <b>Заявка подтверждена</b>", telegram.CreateBaseKeyboard())
	return states.SetStartKeyboard()
}

func RejectTrainingRegistration(botUrl string, chatId int, messageId int, registrationId uint, repo database.ContentRepositoryInterface) states.State {
	registration, err := repo.GetTrainingRegistrationByID(registrationId)
	if err != nil || registration == nil {
		return sendErrorMessage(botUrl, chatId, messageId)
	}

	training, _ := repo.GetTrainingById(registration.TrainingID)
	if training == nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Тренировка не найдена</b>", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	trainer, _ := repo.GetTrainerByID(training.TrainerID)
	if trainer == nil || trainer.ChatId != chatId {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Нет прав</b>", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	registration.Status = "rejected"
	err = repo.UpdateTrainingRegistration(registrationId, registration)
	if err != nil {
		log.Printf("ERROR: Failed to reject training registration %d: %v", registrationId, err)
		return sendErrorMessage(botUrl, chatId, messageId)
	}

	user, _ := repo.GetUserByID(registration.UserID)
	track, _ := repo.GetTrackByID(training.TrackID)

	trackName := "Неизвестная трасса"
	if track != nil {
		trackName = track.Name
	}

	if user != nil {
		userMessage := fmt.Sprintf("❌ <b>Заявка на тренировку отклонена</b>\n\n"+
			"🏃‍♂️ <b>Тренировка:</b> %s\n"+
			"📅 <b>Дата и время:</b> %s\n\n"+
			"💡 <b>Попробуйте записаться на другую тренировку.</b>",
			trackName, training.Time.Format("02.01.2006 15:04"))

		telegram.SendMessage(botUrl, user.ChatId, userMessage, telegram.CreateBaseKeyboard())
	}

	log.Printf("Training registration %d rejected by trainer %d", registrationId, chatId)
	telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Заявка отклонена</b>", telegram.CreateBaseKeyboard())
	return states.SetStartKeyboard()
}

func formatTrainingsListForUsers(trainings []database.Training, repo database.ContentRepositoryInterface) string {

	var builder strings.Builder
	builder.WriteString("📅 <b>Расписание тренировок RVA Academy</b>\n\n")

	for i, training := range trainings {
		trainer, _ := repo.GetTrainerByID(training.TrainerID)
		trainerName := "Неизвестный тренер"
		if trainer != nil {
			trainerName = trainer.Name
		}

		track, _ := repo.GetTrackByID(training.TrackID)
		trackName := "Неизвестная трасса"
		if track != nil {
			trackName = track.Name
		}

		registrations, _ := repo.GetTrainingRegistrationsByTrainingID(training.ID)
		confirmedCount := 0
		var confirmedUsers []string

		for _, reg := range registrations {
			if reg.Status == "confirmed" {
				confirmedCount++
				user, _ := repo.GetUserByID(reg.UserID)
				userName := "Участник"
				if user != nil {
					userName = user.Name
				}
				confirmedUsers = append(confirmedUsers, userName)
			}
		}

		availableSpots := training.MaxParticipants - confirmedCount
		spotsText := fmt.Sprintf("%d мест", availableSpots)
		if availableSpots <= 0 {
			spotsText = "❌ Мест нет"
		} else if availableSpots == 1 {
			spotsText = "1 место"
		}

		builder.WriteString(fmt.Sprintf("🏃‍♂️ <b>%d. Тренировка</b>\n", i+1))
		builder.WriteString(fmt.Sprintf("👨‍🏫 <b>Тренер:</b> %s\n", trainerName))
		builder.WriteString(fmt.Sprintf("🏁 <b>Трасса:</b> %s\n", trackName))
		builder.WriteString(fmt.Sprintf("📅 <b>Дата и время:</b> %s\n", training.Time.Format("02.01.2006 15:04")))
		builder.WriteString(fmt.Sprintf("👥 <b>Свободно:</b> %s\n", spotsText))

		if len(confirmedUsers) > 0 {
			builder.WriteString("✅ <b>Участники:</b> ")
			var displayNames []string
			for _, fullName := range confirmedUsers {
				parts := strings.Fields(fullName)
				if len(parts) > 0 {
					displayNames = append(displayNames, parts[0])
				}
			}
			builder.WriteString(strings.Join(displayNames, ", "))
			builder.WriteString("\n")
		}

		builder.WriteString("\n")
	}

	builder.WriteString("💡 <i>Для записи на тренировку используйте кнопку \"Записаться на тренировку\" в главном меню.</i>")

	return builder.String()
}

func formatTrainersListForUsers(trainers []database.Trainer) string {

	var builder strings.Builder
	builder.WriteString("👥 Тренерский состав RVA Academy\n\n")

	for i, trainer := range trainers {
		builder.WriteString(fmt.Sprintf("👨‍🏫 <b>%d. %s</b>\n", i+1, trainer.Name))
		builder.WriteString(fmt.Sprintf("📱 %s\n", trainer.TgId))
		builder.WriteString(fmt.Sprintf("📝 %s\n\n", trainer.Info))
	}

	return builder.String()
}

func formatTracksListForUsers(tracks []database.Track) string {
	var builder strings.Builder
	builder.WriteString("🏁 Трассы RVA Academy\n\n")

	for i, track := range tracks {
		builder.WriteString(fmt.Sprintf("🏁 <b>%d. %s</b>\n", i+1, track.Name))
		builder.WriteString(fmt.Sprintf("📄 %s\n\n", track.Info))
	}

	return builder.String()
}
