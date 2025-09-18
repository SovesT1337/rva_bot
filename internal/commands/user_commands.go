package commands

import (
	"fmt"
	"log"
	"strings"

	"x.localhost/rvabot/internal/database"
	"x.localhost/rvabot/internal/states"
	"x.localhost/rvabot/internal/telegram"
)

// Основные команды
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

// Информационные команды
func Info(botUrl string, chatId int, messageId int) states.State {
	telegram.EditMessage(botUrl, chatId, messageId, "ℹ️ Информация о RVA Academy\n\n"+
		"", telegram.CreateInfoKeyboard())
	return states.SetStartKeyboard()
}

func ViewELORating(botUrl string, chatId int, messageId int, repo database.ContentRepositoryInterface) states.State {
	telegram.EditMessage(botUrl, chatId, messageId, "📊 <b>Рейтинг ELO RVA Academy</b>\n\n"+
		"🚧 <b>В разработке</b>\n\n"+
		"💡 Система рейтинга будет доступна в ближайшее время.", telegram.CreateBaseKeyboard())
	return states.SetStartKeyboard()
}

func InfoTrainer(botUrl string, chatId int, messageId int, repo database.ContentRepositoryInterface) states.State {
	trainers, err := repo.GetTrainers()
	if err != nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ Ошибка при получении информации о тренерах.\n\n"+
			"Попробуйте позже.", telegram.CreateBackToInfoKeyboard())
		return states.SetStartKeyboard()
	}

	message := formatTrainersListForUsers(trainers)
	telegram.EditMessage(botUrl, chatId, messageId, message, telegram.CreateBackToInfoKeyboard())
	return states.SetStartKeyboard()
}

func InfoTrack(botUrl string, chatId int, messageId int, repo database.ContentRepositoryInterface) states.State {
	tracks, err := repo.GetTracks()
	if err != nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ Ошибка при получении информации о трассах.\n\n"+
			"Попробуйте позже.", telegram.CreateBackToInfoKeyboard())
		return states.SetStartKeyboard()
	}

	message := formatTracksListForUsers(tracks)
	telegram.EditMessage(botUrl, chatId, messageId, message, telegram.CreateBackToInfoKeyboard())
	return states.SetStartKeyboard()
}

func InfoFormat(botUrl string, chatId int, messageId int) states.State {
	message := "📋 <b>Формат тренировок RVA Academy</b>\n\n" +
		"🏃‍♂️ <b>Структура тренировки:</b>\n" +
		"• Разминка (15-20 минут)\n" +
		"• Основная часть (40-60 минут)\n" +
		"• Заминка (10-15 минут)\n\n" +
		"⏰ <b>Продолжительность:</b> 1.5-2 часа\n\n" +
		"👥 <b>Группы:</b>\n" +
		"• Начинающие (до 6 месяцев опыта)\n" +
		"• Продвинутые (от 6 месяцев)\n" +
		"• Профессионалы (соревновательный уровень)\n\n" +
		"🎯 <b>Что включено:</b>\n" +
		"• Техническая подготовка\n" +
		"• Физическая подготовка\n" +
		"• Тактическая подготовка\n" +
		"• Анализ результатов\n\n" +
		"📝 <b>Что взять с собой:</b>\n" +
		"• Спортивная форма\n" +
		"• Сменная обувь\n" +
		"• Вода\n" +
		"• Полотенце\n\n" +
		"💡 <i>Все необходимое оборудование предоставляется академией</i>"

	telegram.EditMessage(botUrl, chatId, messageId, message, telegram.CreateBackToInfoKeyboard())
	return states.SetStartKeyboard()
}

// Просмотр расписания для пользователей
func ViewScheduleUser(botUrl string, chatId int, messageId int, repo database.ContentRepositoryInterface) states.State {
	// Получаем активные тренировки
	trainings, err := repo.GetActiveTrainings()
	if err != nil {
		log.Printf("ERROR: Failed to get active trainings: %v", err)
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Ошибка при получении расписания</b>\n\n"+
			"Попробуйте позже.", telegram.CreateBackToInfoKeyboard())
		return states.SetStartKeyboard()
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

// Регистрация пользователя
func SetUserName(botUrl string, chatId int, update telegram.Update, repo database.ContentRepositoryInterface, state states.State) states.State {
	name := update.Message.Text

	// Получаем временные данные из состояния и обновляем их
	tempData := state.GetTempUserData()
	tempData.Name = name

	// Показываем подтверждение регистрации
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

	// Создаем пользователя в базе данных
	user := &database.User{
		Name:     tempData.Name,
		ChatId:   chatId,
		IsActive: true,
	}

	id, err := repo.CreateUser(user)
	if err != nil {
		log.Printf("ERROR: Failed to create user %s: %v", tempData.Name, err)
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Ошибка при регистрации</b>\n\n"+
			"Попробуйте позже.\n"+
			"Обратитесь к администратору.", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	log.Printf("User created successfully: %s (ID: %d)", tempData.Name, id)
	telegram.EditMessage(botUrl, chatId, messageId, "🎉 <b>Регистрация завершена!</b>\n"+
		"Добро пожаловать, "+tempData.Name+"!", telegram.CreateStartKeyboard())
	return states.SetStartKeyboard()
}

// Регистрация на тренировки
func StartTrainingRegistration(botUrl string, chatId int, messageId int, repo database.ContentRepositoryInterface) states.State {
	// Проверяем, зарегистрирован ли пользователь
	user, err := repo.GetUserByChatId(chatId)
	if err != nil || user == nil {
		// Пользователь не найден, нужно зарегистрироваться
		telegram.EditMessage(botUrl, chatId, messageId, "🏃‍♂️ <b>Запись на тренировку</b>\n\n"+
			"📝 Введите ваше ФИО\n"+
			"<i>Пример: Иванов Иван Иванович</i>", telegram.CreateCancelUserRegistrationKeyboard())
		return states.SetEnterUserName()
	}

	// Пользователь уже зарегистрирован, начинаем пошаговый выбор
	// Шаг 1: Выбор трассы (только те, на которых есть активные тренировки)
	tracks, err := repo.GetTracksWithActiveTrainings()
	if err != nil {
		log.Printf("ERROR: Failed to get tracks with active trainings: %v", err)
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Ошибка получения трасс</b>\n"+
			"Попробуйте позже.", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	if len(tracks) == 0 {
		telegram.EditMessage(botUrl, chatId, messageId, "🏁 <b>Нет доступных трасс</b>\n"+
			"Нет активных тренировок.", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	telegram.EditMessage(botUrl, chatId, messageId, "🏃‍♂️ <b>Запись на тренировку</b>\n\n"+
		"👤 "+user.Name+"\n"+
		"🏁 <b>Шаг 1/3:</b> Трасса", telegram.CreateTrackSelectionForRegistrationKeyboard(tracks))

	// Создаем временные данные для регистрации
	tempData := &states.TempRegistrationData{}
	state := states.SetSelectTrackForRegistration()
	return state.SetTempRegistrationData(tempData)
}

func ConfirmTrainingRegistration(botUrl string, chatId int, messageId int, trainingId uint, repo database.ContentRepositoryInterface) states.State {
	// Получаем информацию о тренировке
	training, err := repo.GetTrainingById(trainingId)
	if err != nil || training == nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Тренировка не найдена</b>\n\n"+
			"🔍 Возможно, тренировка была удалена.", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	// Проверяем, не зарегистрирован ли уже пользователь
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

	// Проверяем, есть ли свободные места
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

	// Получаем информацию о тренере и трассе
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

	// Показываем подтверждение регистрации
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
	// Получаем пользователя
	user, err := repo.GetUserByChatId(chatId)
	if err != nil || user == nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Пользователь не найден</b>\n\n"+
			"🔍 Сначала зарегистрируйтесь в системе.", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	// Создаем регистрацию
	registration := &database.TrainingRegistration{
		TrainingID: trainingId,
		UserID:     user.ID,
		Status:     "pending", // Ожидает подтверждения тренера
	}

	regId, err := repo.CreateTrainingRegistration(registration)
	if err != nil {
		log.Printf("ERROR: Failed to create training registration: %v", err)
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Ошибка при записи на тренировку</b>\n\n"+
			"Попробуйте позже.", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	// Получаем информацию о тренировке и тренере для уведомления
	training, _ := repo.GetTrainingById(trainingId)
	trainer, _ := repo.GetTrainerByID(training.TrainerID)
	track, _ := repo.GetTrackByID(training.TrackID)

	// Отправляем уведомление тренеру
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

// Функции навигации назад при записи на тренировки
func BackToTrackSelection(botUrl string, chatId int, messageId int, repo database.ContentRepositoryInterface, state states.State) states.State {
	// Получаем пользователя
	user, err := repo.GetUserByChatId(chatId)
	if err != nil || user == nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Пользователь не найден</b>", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	// Получаем трассы с активными тренировками
	tracks, err := repo.GetTracksWithActiveTrainings()
	if err != nil {
		log.Printf("ERROR: Failed to get tracks with active trainings: %v", err)
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Ошибка получения трасс</b>\n"+
			"Попробуйте позже.", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	if len(tracks) == 0 {
		telegram.EditMessage(botUrl, chatId, messageId, "🏁 <b>Нет доступных трасс</b>\n"+
			"Нет активных тренировок.", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	telegram.EditMessage(botUrl, chatId, messageId, "🏃‍♂️ <b>Запись на тренировку</b>\n\n"+
		"👤 "+user.Name+"\n"+
		"🏁 <b>Шаг 1/3:</b> Трасса", telegram.CreateTrackSelectionForRegistrationKeyboard(tracks))

	// Создаем временные данные для регистрации
	tempData := &states.TempRegistrationData{}
	newState := states.SetSelectTrackForRegistration()
	return newState.SetTempRegistrationData(tempData)
}

func BackToTrainerSelection(botUrl string, chatId int, messageId int, repo database.ContentRepositoryInterface, state states.State) states.State {
	// Получаем временные данные из состояния
	tempData := state.GetTempRegistrationData()
	if tempData.TrackID == 0 {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Ошибка навигации</b>\n"+
			"Начните заново.", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	// Получаем информацию о выбранной трассе
	track, err := repo.GetTrackByID(tempData.TrackID)
	if err != nil || track == nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Трасса не найдена</b>", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	// Получаем список тренеров, у которых есть активные тренировки на выбранной трассе
	trainers, err := repo.GetTrainersByTrack(tempData.TrackID)
	if err != nil {
		log.Printf("ERROR: Failed to get trainers by track: %v", err)
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Ошибка получения тренеров</b>\n"+
			"Попробуйте позже.", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
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

// Новые функции для пошаговой записи на тренировки
func SelectTrackForRegistration(botUrl string, chatId int, messageId int, trackId uint, repo database.ContentRepositoryInterface, state states.State) states.State {
	// Получаем временные данные из состояния и обновляем их
	tempData := state.GetTempRegistrationData()
	tempData.TrackID = trackId

	// Получаем информацию о выбранной трассе
	track, err := repo.GetTrackByID(trackId)
	if err != nil || track == nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Трасса не найдена</b>", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	// Получаем список тренеров, у которых есть активные тренировки на выбранной трассе
	trainers, err := repo.GetTrainersByTrack(trackId)
	if err != nil {
		log.Printf("ERROR: Failed to get trainers by track: %v", err)
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Ошибка получения тренеров</b>\n"+
			"Попробуйте позже.", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
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
	// Получаем временные данные из состояния и обновляем их
	tempData := state.GetTempRegistrationData()
	tempData.TrainerID = trainerId

	// Получаем информацию о выбранном тренере
	trainer, err := repo.GetTrainerByID(trainerId)
	if err != nil || trainer == nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Тренер не найден</b>\n\n"+
			"🔍 Возможно, тренер был удален.", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	// Получаем информацию о выбранной трассе
	track, err := repo.GetTrackByID(tempData.TrackID)
	if err != nil || track == nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Трасса не найдена</b>", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	// Получаем доступные тренировки для выбранной трассы и тренера
	trainings, err := repo.GetActiveTrainingsByTrackAndTrainer(tempData.TrackID, trainerId)
	if err != nil {
		log.Printf("ERROR: Failed to get trainings by track and trainer: %v", err)
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Ошибка при получении расписания</b>\n\n"+
			"Попробуйте позже.", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	if len(trainings) == 0 {
		telegram.EditMessage(botUrl, chatId, messageId, "📅 <b>Нет доступных тренировок</b>\n\n"+
			"🏃‍♂️ <b>Тренер:</b> "+trainer.Name+"\n"+
			"🏁 <b>Трасса:</b> "+track.Name+"\n\n"+
			"📝 <b>У выбранного тренера нет активных тренировок на этой трассе.</b>\n"+
			"💡 Попробуйте выбрать другого тренера или трассу.", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	// Сортируем тренировки по времени
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
	// Получаем информацию о выбранной тренировке
	training, err := repo.GetTrainingById(trainingId)
	if err != nil || training == nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Тренировка не найдена</b>\n\n"+
			"🔍 Возможно, тренировка была удалена.", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	// Проверяем, что тренировка соответствует выбранным трассе и тренеру
	tempData := state.GetTempRegistrationData()
	if training.TrackID != tempData.TrackID || training.TrainerID != tempData.TrainerID {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Ошибка выбора тренировки</b>\n\n"+
			"🔍 Выбранная тренировка не соответствует выбранным трассе и тренеру.", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	// Переходим к подтверждению регистрации
	return ConfirmTrainingRegistration(botUrl, chatId, messageId, trainingId, repo)
}

// Команды для тренеров (подтверждение/отклонение заявок)
func ApproveTrainingRegistration(botUrl string, chatId int, messageId int, registrationId uint, repo database.ContentRepositoryInterface) states.State {
	// Получаем регистрацию
	registration, err := repo.GetTrainingRegistrationByID(registrationId)
	if err != nil || registration == nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Заявка не найдена</b>", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	// Проверяем, что тренер имеет право подтверждать эту заявку
	training, _ := repo.GetTrainingById(registration.TrainingID)
	if training == nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Тренировка не найдена</b>\n\n"+
			"🔍 Возможно, тренировка была удалена.", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	trainer, _ := repo.GetTrainerByID(training.TrainerID)
	if trainer == nil || trainer.ChatId != chatId {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Нет прав</b>\n\n"+
			"🔒 У вас нет прав для подтверждения этой заявки.", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	// Обновляем статус регистрации
	registration.Status = "confirmed"
	err = repo.UpdateTrainingRegistration(registrationId, registration)
	if err != nil {
		log.Printf("ERROR: Failed to approve training registration %d: %v", registrationId, err)
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Ошибка при подтверждении заявки</b>\n\n"+
			"Попробуйте позже.", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	// Получаем информацию о пользователе для уведомления
	user, _ := repo.GetUserByID(registration.UserID)
	track, _ := repo.GetTrackByID(training.TrackID)

	trackName := "Неизвестная трасса"
	if track != nil {
		trackName = track.Name
	}

	// Отправляем уведомление пользователю
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
	// Получаем регистрацию
	registration, err := repo.GetTrainingRegistrationByID(registrationId)
	if err != nil || registration == nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Заявка не найдена</b>", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	// Проверяем, что тренер имеет право отклонять эту заявку
	training, _ := repo.GetTrainingById(registration.TrainingID)
	if training == nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Тренировка не найдена</b>\n\n"+
			"🔍 Возможно, тренировка была удалена.", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	trainer, _ := repo.GetTrainerByID(training.TrainerID)
	if trainer == nil || trainer.ChatId != chatId {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Нет прав</b>\n\n"+
			"🔒 У вас нет прав для отклонения этой заявки.", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	// Обновляем статус регистрации
	registration.Status = "rejected"
	err = repo.UpdateTrainingRegistration(registrationId, registration)
	if err != nil {
		log.Printf("ERROR: Failed to reject training registration %d: %v", registrationId, err)
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Ошибка при отклонении заявки</b>\n\n"+
			"Попробуйте позже.", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	// Получаем информацию о пользователе для уведомления
	user, _ := repo.GetUserByID(registration.UserID)
	track, _ := repo.GetTrackByID(training.TrackID)

	trackName := "Неизвестная трасса"
	if track != nil {
		trackName = track.Name
	}

	// Отправляем уведомление пользователю
	if user != nil {
		userMessage := fmt.Sprintf("❌ <b>Заявка на тренировку отклонена</b>\n\n"+
			"😔 <b>К сожалению, ваша заявка на тренировку была отклонена тренером.</b>\n\n"+
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

// Форматирование списка тренировок для пользователей
func formatTrainingsListForUsers(trainings []database.Training, repo database.ContentRepositoryInterface) string {
	if len(trainings) == 0 {
		return "📅 <b>Расписание тренировок</b>\n\n" +
			"📝 <b>Активных тренировок пока нет</b>\n\n" +
			"💡 Следите за обновлениями!"
	}

	var builder strings.Builder
	builder.WriteString("📅 <b>Расписание тренировок RVA Academy</b>\n\n")

	for i, training := range trainings {
		// Получаем информацию о тренере
		trainer, _ := repo.GetTrainerByID(training.TrainerID)
		trainerName := "Неизвестный тренер"
		if trainer != nil {
			trainerName = trainer.Name
		}

		// Получаем информацию о трассе
		track, _ := repo.GetTrackByID(training.TrackID)
		trackName := "Неизвестная трасса"
		if track != nil {
			trackName = track.Name
		}

		// Получаем количество зарегистрированных участников
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

		// Проверяем, есть ли свободные места
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

		// Показываем список участников (только имена, без фамилий для приватности)
		if len(confirmedUsers) > 0 {
			builder.WriteString("✅ <b>Участники:</b> ")
			// Показываем только первые имена для приватности
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

// Форматирование списка тренеров для пользователей
func formatTrainersListForUsers(trainers []database.Trainer) string {
	if len(trainers) == 0 {
		return "👥 Тренерский состав\n\n" +
			"Информация о тренерах будет доступна в ближайшее время."
	}

	var builder strings.Builder
	builder.WriteString("👥 Тренерский состав RVA Academy\n\n")

	for i, trainer := range trainers {
		builder.WriteString(fmt.Sprintf("👨‍🏫 <b>%d. %s</b>\n", i+1, trainer.Name))

		// Добавляем ссылку на тренера в Telegram, если есть TgId
		if trainer.TgId != "" {
			// Создаем ссылку на пользователя в Telegram
			builder.WriteString(fmt.Sprintf("📱 <a href=\"https://t.me/%s\">Написать тренеру</a>\n", trainer.TgId))
		}

		if trainer.Info != "" {
			builder.WriteString(fmt.Sprintf("📝 %s\n", trainer.Info))
		}

		builder.WriteString("\n")
	}

	return builder.String()
}

// Форматирование списка трасс для пользователей
func formatTracksListForUsers(tracks []database.Track) string {
	if len(tracks) == 0 {
		return "🏁 Информация о трассах\n\n" +
			"Информация о трассах будет доступна в ближайшее время."
	}

	var builder strings.Builder
	builder.WriteString("🏁 Трассы RVA Academy\n\n")

	for i, track := range tracks {
		builder.WriteString(fmt.Sprintf("🏁 <b>%d. %s</b>\n", i+1, track.Name))

		if track.Info != "" {
			builder.WriteString(fmt.Sprintf("📄 %s\n", track.Info))
		}

		builder.WriteString("\n")
	}

	return builder.String()
}
