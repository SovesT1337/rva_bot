package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"x.localhost/rvabot/internal/database"
	"x.localhost/rvabot/internal/logger"
	"x.localhost/rvabot/internal/states"
	"x.localhost/rvabot/internal/telegram"
	"x.localhost/rvabot/internal/validation"
)

func Admin(botUrl string, chatId int, repo database.ContentRepositoryInterface) states.State {
	if !database.IsAdmin(chatId, repo) {
		telegram.SendMessage(botUrl, chatId, "🚫 <b>Доступ запрещен</b>\n\n"+
			"❌ У вас нет прав администратора для доступа к этой панели.\n\n"+
			"💡 Обратитесь к администратору для получения доступа.", telegram.CreateBaseKeyboard())
		return states.SetStartKeyboard()
	}

	telegram.SendMessage(botUrl, chatId, "⚙️ <b>Панель администратора</b>\n\n"+
		"🎛️ Добро пожаловать в систему управления!\n\n"+
		"📋 <b>Доступные разделы:</b>\n"+
		"👨‍🏫 Управление тренерами\n"+
		"🏁 Управление трассами\n"+
		"📅 Управление расписанием\n\n"+
		"", telegram.CreateAdminKeyboard())
	return states.SetAdminKeyboard()
}

func CreateTrainer(botUrl string, chatId int, messageId int) states.State {
	telegram.EditMessage(botUrl, chatId, messageId, "👨‍🏫 <b>Добавление нового тренера</b>\n\n"+
		"📝 <b>Шаг 1 из 3:</b> Введите ФИО тренера\n\n"+
		"💡 <i>Пример: Иванов Иван Иванович</i>", telegram.CreateBackToTrainersMenuKeyboard())

	tempData := &states.TempTrainerData{}
	state := states.SetEnterTrainerName(0)
	return state.SetTempTrainerData(tempData)
}

func SetTrainerName(botUrl string, chatId int, update telegram.Update, repo database.ContentRepositoryInterface, state states.State) states.State {
	name := update.Message.Text

	// Валидация имени тренера
	validator := validation.NewValidator()
	result := validator.ValidateTrainerName(name)
	if !result.IsValid {
		errorMsg := strings.Join(result.GetErrorMessages(), "\n")
		telegram.SendMessage(botUrl, chatId, "❌ <b>Ошибка валидации</b>\n\n"+errorMsg+"\n\n🔄 Попробуйте еще раз:", telegram.CreateCancelKeyboard())
		return state
	}

	tempData := state.GetTempTrainerData()
	tempData.Name = name

	telegram.SendMessage(botUrl, chatId, "👨‍🏫 <b>Добавление нового тренера</b>\n\n"+
		"📱 <b>Шаг 2 из 4:</b> Введите Telegram ID тренера\n"+
		"💡 <i>Пример: @username или 123456789</i>", telegram.CreateCancelKeyboard())

	newState := states.SetEnterTrainerTgId(0)
	return newState.SetTempTrainerData(tempData)
}

func SetTrainerTgId(botUrl string, chatId int, update telegram.Update, repo database.ContentRepositoryInterface, state states.State) states.State {
	tgid := update.Message.Text

	// Валидация Telegram ID
	validator := validation.NewValidator()
	result := validator.ValidateTelegramID(tgid)
	if !result.IsValid {
		errorMsg := strings.Join(result.GetErrorMessages(), "\n")
		telegram.SendMessage(botUrl, chatId, "❌ <b>Ошибка валидации</b>\n\n"+errorMsg+"\n\n🔄 Попробуйте еще раз:", telegram.CreateCancelKeyboard())
		return state
	}

	tempData := state.GetTempTrainerData()
	tempData.TgId = tgid

	telegram.SendMessage(botUrl, chatId, "👨‍🏫 <b>Добавление нового тренера</b>\n\n"+
		"💬 <b>Шаг 3 из 4:</b> Введите Chat ID тренера\n"+
		"💡 <i>Пример: 123456789 (числовой ID чата)</i>", telegram.CreateCancelKeyboard())

	newState := states.SetEnterTrainerChatId(0)
	return newState.SetTempTrainerData(tempData)
}

func SetTrainerChatId(botUrl string, chatId int, update telegram.Update, repo database.ContentRepositoryInterface, state states.State) states.State {
	chatIdStr := update.Message.Text

	// Валидация Chat ID
	validator := validation.NewValidator()
	result := validator.ValidateChatID(chatIdStr)
	if !result.IsValid {
		errorMsg := strings.Join(result.GetErrorMessages(), "\n")
		telegram.SendMessage(botUrl, chatId, "❌ <b>Ошибка валидации</b>\n\n"+errorMsg+"\n\n🔄 Попробуйте еще раз:", telegram.CreateCancelKeyboard())
		return state
	}

	trainerChatId, _ := strconv.Atoi(chatIdStr) // Ошибка уже проверена в валидации

	tempData := state.GetTempTrainerData()
	tempData.ChatId = trainerChatId

	telegram.SendMessage(botUrl, chatId, "👨‍🏫 <b>Добавление нового тренера</b>\n\n"+
		"📝 <b>Шаг 4 из 4:</b> Введите информацию о тренере\n"+
		"💡 <i>Пример: Опытный тренер по бегу, стаж 5 лет</i>", telegram.CreateCancelKeyboard())

	newState := states.SetEnterTrainerInfo(0)
	return newState.SetTempTrainerData(tempData)
}

func SetTrainerInfo(botUrl string, chatId int, update telegram.Update, repo database.ContentRepositoryInterface, state states.State) states.State {
	info := update.Message.Text

	// Валидация информации о тренере
	validator := validation.NewValidator()
	result := validator.ValidateTrainerInfo(info)
	if !result.IsValid {
		errorMsg := strings.Join(result.GetErrorMessages(), "\n")
		telegram.SendMessage(botUrl, chatId, "❌ <b>Ошибка валидации</b>\n\n"+errorMsg+"\n\n🔄 Попробуйте еще раз:", telegram.CreateCancelKeyboard())
		return state
	}

	tempData := state.GetTempTrainerData()
	tempData.Info = info

	message := fmt.Sprintf("✅ <b>Подтверждение создания тренера</b>\n\n"+
		"📋 <b>Проверьте данные:</b>\n\n"+
		"👤 <b>ФИО:</b> %s\n"+
		"📱 <b>Telegram ID:</b> %s\n"+
		"💬 <b>Chat ID:</b> %d\n"+
		"📝 <b>Информация:</b> %s\n\n"+
		"❓ <b>Создать тренера с этими данными?</b>", tempData.Name, tempData.TgId, tempData.ChatId, tempData.Info)

	telegram.SendMessage(botUrl, chatId, message, telegram.CreateConfirmationKeyboard())

	newState := states.SetConfirmTrainerCreation()
	return newState.SetTempTrainerData(tempData)
}

func ConfirmTrainerCreation(botUrl string, chatId int, messageId int, repo database.ContentRepositoryInterface, tempData *states.TempTrainerData) states.State {
	logger.AdminInfo(chatId, "Создание тренера: %s", tempData.Name)

	trainer := &database.Trainer{
		Name:   tempData.Name,
		TgId:   tempData.TgId,
		ChatId: tempData.ChatId,
		Info:   tempData.Info,
	}

	_, err := repo.CreateTrainer(trainer)
	if err != nil {
		logger.AdminError(chatId, "Создание тренера %s: %v", tempData.Name, err)
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Ошибка при создании тренера</b>\n"+
			"Попробуйте позже.", telegram.CreateBackToTrainersMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	logger.AdminInfo(chatId, "Тренер создан: %s", tempData.Name)
	telegram.EditMessage(botUrl, chatId, messageId, "🎉 <b>Тренер создан!</b>\n\n"+
		"👤 <b>Имя:</b> "+tempData.Name+"\n"+
		"📱 <b>Telegram ID:</b> "+tempData.TgId+"\n"+
		"💬 <b>Chat ID:</b> "+fmt.Sprintf("%d", tempData.ChatId)+"\n"+
		"📝 <b>Информация:</b> "+tempData.Info+"\n\n"+
		"✨ Тренер добавлен в систему!", telegram.CreateBackToTrainersMenuKeyboard())
	return states.SetAdminKeyboard()
}

func CancelTrainerCreation(botUrl string, chatId int, messageId int) states.State {
	telegram.EditMessage(botUrl, chatId, messageId, "🚫 <b>Создание тренера отменено</b>\n\n"+
		"💡 Вы можете создать тренера позже через меню управления.\n"+
		"🔄 Все введенные данные были сброшены.", telegram.CreateBackToTrainersMenuKeyboard())
	return states.SetAdminKeyboard()
}

func ViewTrainers(botUrl string, chatId int, messageId int, repo database.ContentRepositoryInterface) states.State {
	trainers, err := repo.GetTrainers()
	if err != nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Ошибка получения тренеров</b>\n"+
			"Попробуйте позже.", telegram.CreateBackToTrainersMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	message := formatTrainersListForAdmin(trainers)
	telegram.EditMessage(botUrl, chatId, messageId, message, telegram.CreateBackToTrainersMenuKeyboard())
	return states.SetAdminKeyboard()
}

func formatTrainersListForAdmin(trainers []database.Trainer) string {
	if len(trainers) == 0 {
		return "👥 <b>Список тренеров пуст</b>\n\n" +
			"👨‍🏫 Добавьте первого тренера через админ-панель."
	}

	var builder strings.Builder
	builder.WriteString("👥 <b>Список тренеров RVA Academy</b>\n\n")

	for i, trainer := range trainers {
		builder.WriteString(fmt.Sprintf("👤 <b>%d. %s</b>\n", i+1, trainer.Name))

		if trainer.TgId != "" {
			builder.WriteString(fmt.Sprintf("📱 <b>Telegram ID:</b> <code>%s</code>\n", trainer.TgId))
		}

		if trainer.Info != "" {
			builder.WriteString(fmt.Sprintf("📄 <b>Информация:</b> %s\n", trainer.Info))
		}

		builder.WriteString(fmt.Sprintf("📅 <b>Добавлен:</b> %s\n", trainer.CreatedAt.Format("02.01.2006")))
		builder.WriteString("\n")
	}

	return builder.String()
}

func EditTrainerName(botUrl string, chatId int, messageId int, trainerId uint) states.State {
	telegram.EditMessage(botUrl, chatId, messageId, "✏️ <b>Редактирование ФИО тренера</b>\n\n"+
		"📝 Введите новое ФИО тренера:\n\n"+
		"💡 <i>Пример: Иванов Иван Иванович</i>", telegram.CreateBackToTrainersMenuKeyboard())
	return states.SetEditTrainerName(trainerId)
}

func SetEditTrainerName(botUrl string, chatId int, update telegram.Update, repo database.ContentRepositoryInterface, trainerId uint) states.State {
	name := update.Message.Text
	logger.AdminInfo(chatId, "Обновление тренера %d: %s", trainerId, name)

	// Получаем существующего тренера
	trainer, err := repo.GetTrainerByID(trainerId)
	if err != nil {
		logger.AdminError(chatId, "Получение тренера %d: %v", trainerId, err)
		telegram.SendMessage(botUrl, chatId, "❌ <b>Тренер не найден</b>\n\n"+
			"🔍 Возможно, тренер был удален.", telegram.CreateBackToTrainersMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	// Обновляем только имя
	trainer.Name = name
	err = repo.UpdateTrainer(trainerId, trainer)
	if err != nil {
		logger.AdminError(chatId, "Обновление тренера %d: %v", trainerId, err)
		telegram.SendMessage(botUrl, chatId, "❌ <b>Ошибка обновления имени тренера</b>\n\n"+
			"Ошибка сохранения.\n"+
			"Обратитесь к администратору.", telegram.CreateBackToTrainersMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	logger.AdminInfo(chatId, "Тренер %d обновлен: %s", trainerId, name)
	telegram.SendMessage(botUrl, chatId, "✅ <b>ФИО тренера обновлено!</b>\n\n"+
		"👤 Новое имя: "+name, telegram.CreateBackToTrainersMenuKeyboard())
	return states.SetAdminKeyboard()
}

func EditTrainerTgId(botUrl string, chatId int, messageId int, trainerId uint) states.State {
	telegram.EditMessage(botUrl, chatId, messageId, "✏️ <b>Редактирование Telegram ID</b>\n\n"+
		"📱 Введите новый Telegram ID тренера:\n\n"+
		"💡 <i>Пример: @username или 123456789</i>", telegram.CreateBackToTrainersMenuKeyboard())
	return states.SetEditTrainerTgId(trainerId)
}

func SetEditTrainerTgId(botUrl string, chatId int, update telegram.Update, repo database.ContentRepositoryInterface, trainerId uint) states.State {
	tgId := update.Message.Text

	// Получаем существующего тренера
	trainer, err := repo.GetTrainerByID(trainerId)
	if err != nil {
		logger.AdminError(chatId, "Получение тренера %d: %v", trainerId, err)
		telegram.SendMessage(botUrl, chatId, "❌ <b>Тренер не найден</b>\n\n"+
			"🔍 Возможно, тренер был удален.", telegram.CreateBackToTrainersMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	// Обновляем только Telegram ID
	trainer.TgId = tgId
	err = repo.UpdateTrainer(trainerId, trainer)
	if err != nil {
		telegram.SendMessage(botUrl, chatId, "❌ <b>Ошибка обновления Telegram ID</b>\n\n"+
			"Ошибка сохранения.\n"+
			"Обратитесь к администратору.", telegram.CreateBackToTrainersMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	telegram.SendMessage(botUrl, chatId, "✅ <b>Telegram ID обновлен!</b>\n\n"+
		"📱 Новый ID: "+tgId, telegram.CreateBackToTrainersMenuKeyboard())
	return states.SetAdminKeyboard()
}

func EditTrainerInfo(botUrl string, chatId int, messageId int, trainerId uint) states.State {
	telegram.EditMessage(botUrl, chatId, messageId, "✏️ <b>Редактирование информации о тренере</b>\n\n"+
		"📋 Введите новую информацию о тренере:\n\n"+
		"💡 <i>Пример: Опытный тренер по бегу, стаж 5 лет</i>", telegram.CreateBackToTrainersMenuKeyboard())
	return states.SetEditTrainerInfo(trainerId)
}

func SetEditTrainerInfo(botUrl string, chatId int, update telegram.Update, repo database.ContentRepositoryInterface, trainerId uint) states.State {
	info := update.Message.Text

	// Получаем существующего тренера
	trainer, err := repo.GetTrainerByID(trainerId)
	if err != nil {
		logger.AdminError(chatId, "Получение тренера %d: %v", trainerId, err)
		telegram.SendMessage(botUrl, chatId, "❌ <b>Тренер не найден</b>\n\n"+
			"🔍 Возможно, тренер был удален.", telegram.CreateBackToTrainersMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	// Обновляем только информацию
	trainer.Info = info
	err = repo.UpdateTrainer(trainerId, trainer)
	if err != nil {
		telegram.SendMessage(botUrl, chatId, "❌ <b>Ошибка обновления информации о тренере</b>\n\n"+
			"Ошибка сохранения.\n"+
			"Обратитесь к администратору.", telegram.CreateBackToTrainersMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	telegram.SendMessage(botUrl, chatId, "✅ <b>Информация о тренере обновлена!</b>\n\n"+
		"📄 Новая информация: "+info, telegram.CreateBackToTrainersMenuKeyboard())
	return states.SetAdminKeyboard()
}

func ConfirmTrainerDeletion(botUrl string, chatId int, messageId int, trainerId uint, repo database.ContentRepositoryInterface) states.State {
	trainer, err := repo.GetTrainerByID(trainerId)
	if err != nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ Тренер не найден.", telegram.CreateBackToTrainersMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	message := fmt.Sprintf("⚠️ <b>Подтверждение удаления тренера</b>\n\n"+
		"👤 <b>Тренер:</b> %s\n"+
		"📱 <b>Telegram ID:</b> %s\n"+
		"📄 <b>Информация:</b> %s\n\n"+
		"🚨 <b>ВНИМАНИЕ!</b> Это действие нельзя отменить!\n\n"+
		"❓ <b>Вы уверены, что хотите удалить этого тренера?</b>",
		trainer.Name, trainer.TgId, trainer.Info)

	telegram.EditMessage(botUrl, chatId, messageId, message, telegram.CreateDeletionConfirmationKeyboard(trainerId))
	return states.SetConfirmTrainerDelete(trainerId)
}

func ExecuteTrainerDeletion(botUrl string, chatId int, messageId int, trainerId uint, repo database.ContentRepositoryInterface) states.State {
	trainer, err := repo.GetTrainerByID(trainerId)
	if err != nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Тренер не найден</b>\n\n"+
			"🔍 Возможно, тренер уже был удален.", telegram.CreateBackToTrainersMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	err = repo.DeleteTrainer(trainerId)
	if err != nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Ошибка удаления тренера</b>\n\n"+
			"Попробуйте позже.", telegram.CreateBackToTrainersMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	telegram.EditMessage(botUrl, chatId, messageId, fmt.Sprintf("🗑️ <b>Тренер %s удален</b>", trainer.Name), telegram.CreateBackToTrainersMenuKeyboard())
	return states.SetAdminKeyboard()
}

func CreateTrack(botUrl string, chatId int, messageId int) states.State {
	telegram.EditMessage(botUrl, chatId, messageId, "🏁 <b>Создание новой трассы</b>\n\n"+
		"📝 Введите название трассы:\n\n"+
		"💡 <i>Пример: Трасса №1 - Легкая</i>", telegram.CreateBackToTracksMenuKeyboard())
	return states.SetEnterTrackName(0)
}

func SetTrackName(botUrl string, chatId int, update telegram.Update, repo database.ContentRepositoryInterface, state states.State) states.State {
	name := update.Message.Text
	logger.AdminInfo(chatId, "Название трека: %s", name)

	tempData := &states.TempTrackData{Name: name}
	newState := states.SetEnterTrackInfo(0).SetTempTrackData(tempData)

	telegram.SendMessage(botUrl, chatId, "📋 Введите описание трассы:\n\n"+
		"💡 <i>Пример: Легкая трасса для начинающих, длина 1 км</i>", telegram.CreateBackToTracksMenuKeyboard())
	return newState
}

func SetTrackInfo(botUrl string, chatId int, update telegram.Update, repo database.ContentRepositoryInterface, state states.State) states.State {
	info := update.Message.Text
	logger.AdminInfo(chatId, "Информация о треке: %s", info)

	tempData := state.GetTempTrackData()
	tempData.Info = info

	message := fmt.Sprintf("🏁 <b>Подтверждение создания трассы</b>\n\n"+
		"📝 <b>Название:</b> %s\n"+
		"📋 <b>Описание:</b> %s\n\n"+
		"❓ <b>Создать трассу?</b>",
		tempData.Name, tempData.Info)

	telegram.SendMessage(botUrl, chatId, message, telegram.CreateConfirmationKeyboard())
	return states.SetConfirmTrackCreation().SetTempTrackData(tempData)
}

func ConfirmTrackCreation(botUrl string, chatId int, messageId int, repo database.ContentRepositoryInterface, tempData *states.TempTrackData) states.State {
	track := &database.Track{
		Name: tempData.Name,
		Info: tempData.Info,
	}

	_, err := repo.CreateTrack(track)
	if err != nil {
		logger.AdminError(chatId, "Создание трека: %v", err)
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Ошибка создания трассы</b>\n\n"+
			"Ошибка сохранения.\n"+
			"Обратитесь к администратору.", telegram.CreateBackToTracksMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	logger.AdminInfo(chatId, "Трек создан: %s", track.Name)
	telegram.EditMessage(botUrl, chatId, messageId, "✅ <b>Трасса создана!</b>\n\n"+
		"🏁 Название: "+track.Name, telegram.CreateBackToTracksMenuKeyboard())
	return states.SetAdminKeyboard()
}

func CancelTrackCreation(botUrl string, chatId int, messageId int) states.State {
	telegram.EditMessage(botUrl, chatId, messageId, "🚫 <b>Создание трассы отменено</b>\n\n"+
		"💡 Вы можете создать трассу позже.", telegram.CreateBackToTracksMenuKeyboard())
	return states.SetAdminKeyboard()
}

func ViewTracks(botUrl string, chatId int, messageId int, repo database.ContentRepositoryInterface) states.State {
	tracks, err := repo.GetTracks()
	if err != nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Ошибка загрузки трасс</b>\n\n"+
			"Попробуйте позже.", telegram.CreateBackToTracksMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	if len(tracks) == 0 {
		telegram.EditMessage(botUrl, chatId, messageId, "📭 <b>Трассы не найдены</b>\n\n"+
			"Сначала создайте трассы.", telegram.CreateBackToTracksMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	message := "🏁 <b>Список трасс:</b>\n\n"
	message += formatTracksListForAdmin(tracks)

	telegram.EditMessage(botUrl, chatId, messageId, message, telegram.CreateBackToTracksMenuKeyboard())
	return states.SetAdminKeyboard()
}

func ViewSchedule(botUrl string, chatId int, messageId int, repo database.ContentRepositoryInterface) states.State {
	trainings, err := repo.GetTrainings()
	if err != nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Ошибка загрузки расписания</b>\n\n"+
			"Попробуйте позже.", telegram.CreateBackToScheduleMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	if len(trainings) == 0 {
		telegram.EditMessage(botUrl, chatId, messageId, "📭 <b>Тренировки не найдены</b>\n\n"+
			"Сначала создайте тренировки.", telegram.CreateBackToScheduleMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	message := "📅 <b>Расписание тренировок:</b>\n\n"
	message += formatTrainingsListForAdmin(trainings, repo)

	telegram.EditMessage(botUrl, chatId, messageId, message, telegram.CreateBackToScheduleMenuKeyboard())
	return states.SetAdminKeyboard()
}

func EditSchedule(botUrl string, chatId int, messageId int, repo database.ContentRepositoryInterface) states.State {
	trainings, err := repo.GetTrainings()
	if err != nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Ошибка загрузки тренировок</b>\n\n"+
			"Попробуйте позже.", telegram.CreateBackToScheduleMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	if len(trainings) == 0 {
		telegram.EditMessage(botUrl, chatId, messageId, "📭 <b>Тренировки не найдены</b>\n\n"+
			"Сначала создайте тренировки.", telegram.CreateBackToScheduleMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	message := "✏️ <b>Выберите тренировку для редактирования:</b>\n\n"
	message += formatTrainingsListForAdmin(trainings, repo)

	telegram.EditMessage(botUrl, chatId, messageId, message, telegram.CreateTrainingEditKeyboard(0))
	return states.SetAdminKeyboard()
}

func CreateTraining(botUrl string, chatId int, messageId int, repo database.ContentRepositoryInterface) states.State {
	tracks, err := repo.GetTracks()
	if err != nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Ошибка загрузки трасс</b>\n\n"+
			"Попробуйте позже.", telegram.CreateBackToScheduleMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	if len(tracks) == 0 {
		telegram.EditMessage(botUrl, chatId, messageId, "📭 <b>Трассы не найдены</b>\n\n"+
			"Сначала создайте трассы.", telegram.CreateBackToScheduleMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	message := "🏁 <b>Выберите трассу для тренировки:</b>\n\n"
	message += formatTracksListForAdmin(tracks)

	telegram.EditMessage(botUrl, chatId, messageId, message, telegram.CreateTrackSelectionForTrainingKeyboard(tracks))
	return states.SetSetTrainingTrack(0)
}

func SetTrainingTrainer(botUrl string, chatId int, messageId int, trainerId uint, repo database.ContentRepositoryInterface, state states.State) states.State {
	telegram.EditMessage(botUrl, chatId, messageId, "🕐 Введите время начала тренировки:\n\n"+
		"💡 <i>Пример: 2024-01-15 18:00</i>", telegram.CreateBackToScheduleMenuKeyboard())

	// Сохраняем данные в состоянии
	newState := states.SetSetTrainingStartTime(0)
	newState.Data["trackId"] = state.Data["trackId"]
	newState.Data["trainerId"] = trainerId
	return newState
}

func SetTrainingTrack(botUrl string, chatId int, messageId int, trackId uint, repo database.ContentRepositoryInterface, state states.State) states.State {
	trainers, err := repo.GetTrainers()
	if err != nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Ошибка загрузки тренеров</b>\n\n"+
			"Попробуйте позже.", telegram.CreateBackToScheduleMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	if len(trainers) == 0 {
		telegram.EditMessage(botUrl, chatId, messageId, "📭 <b>Тренеры не найдены</b>\n\n"+
			"Сначала создайте тренеров.", telegram.CreateBackToScheduleMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	message := "👨‍🏫 <b>Выберите тренера для тренировки:</b>\n\n"
	message += formatTrainersListForAdmin(trainers)

	telegram.EditMessage(botUrl, chatId, messageId, message, telegram.CreateTrainerSelectionForTrainingKeyboard(trainers))

	// Сохраняем trackId в состоянии
	newState := states.SetSetTrainingTrainer(0)
	newState.Data["trackId"] = trackId
	return newState
}

func SetTrainingStartTime(botUrl string, chatId int, update telegram.Update, repo database.ContentRepositoryInterface, state states.State) states.State {
	startTime := update.Message.Text
	logger.AdminInfo(chatId, "Время начала: %s", startTime)

	// Валидируем введенное время
	validator := validation.NewValidator()
	if result := validator.ValidateDateTime(startTime); !result.IsValid {
		errorMsg := "❌ <b>Неверный формат времени</b>\n\n"
		for _, err := range result.Errors {
			errorMsg += fmt.Sprintf("• %s\n", err.Error())
		}
		errorMsg += "\n💡 <i>Пример: 2024-01-15 20:00</i>"

		telegram.SendMessage(botUrl, chatId, errorMsg, telegram.CreateBackToScheduleMenuKeyboard())
		// Сохраняем данные из текущего состояния
		newState := states.SetSetTrainingStartTime(0)
		newState.Data["trackId"] = state.Data["trackId"]
		newState.Data["trainerId"] = state.Data["trainerId"]
		return newState
	}

	telegram.SendMessage(botUrl, chatId, "🕕 Введите время окончания тренировки:\n\n"+
		"💡 <i>Пример: 2024-01-15 20:00</i>", telegram.CreateBackToScheduleMenuKeyboard())

	// Сохраняем данные в состоянии
	newState := states.SetSetTrainingEndTime(0)
	newState.Data["trackId"] = state.Data["trackId"]
	newState.Data["trainerId"] = state.Data["trainerId"]
	newState.Data["startTime"] = startTime
	return newState
}

func SetTrainingEndTime(botUrl string, chatId int, update telegram.Update, repo database.ContentRepositoryInterface, state states.State) states.State {
	endTime := update.Message.Text
	logger.AdminInfo(chatId, "Время окончания: %s", endTime)

	// Валидируем введенное время
	validator := validation.NewValidator()
	if result := validator.ValidateDateTime(endTime); !result.IsValid {
		errorMsg := "❌ <b>Неверный формат времени</b>\n\n"
		for _, err := range result.Errors {
			errorMsg += fmt.Sprintf("• %s\n", err.Error())
		}
		errorMsg += "\n💡 <i>Пример: 2024-01-15 20:00</i>"

		telegram.SendMessage(botUrl, chatId, errorMsg, telegram.CreateBackToScheduleMenuKeyboard())
		// Сохраняем данные из текущего состояния
		newState := states.SetSetTrainingEndTime(0)
		newState.Data["trackId"] = state.Data["trackId"]
		newState.Data["trainerId"] = state.Data["trainerId"]
		newState.Data["startTime"] = state.Data["startTime"]
		return newState
	}

	// Проверяем, что время окончания после времени начала
	startTimeStr, ok := state.Data["startTime"].(string)
	if ok {
		startTime, err1 := time.Parse("2006-01-02 15:04", startTimeStr)
		endTimeParsed, err2 := time.Parse("2006-01-02 15:04", endTime)

		if err1 == nil && err2 == nil {
			if endTimeParsed.Before(startTime) || endTimeParsed.Equal(startTime) {
				telegram.SendMessage(botUrl, chatId, "❌ <b>Неверное время окончания</b>\n\n"+
					"Время окончания должно быть после времени начала.\n"+
					"💡 <i>Пример: 2024-01-15 20:00</i>", telegram.CreateBackToScheduleMenuKeyboard())
				// Сохраняем данные из текущего состояния
				newState := states.SetSetTrainingEndTime(0)
				newState.Data["trackId"] = state.Data["trackId"]
				newState.Data["trainerId"] = state.Data["trainerId"]
				newState.Data["startTime"] = state.Data["startTime"]
				return newState
			}
		}
	}

	telegram.SendMessage(botUrl, chatId, "👥 Введите максимальное количество участников:\n\n"+
		"💡 <i>Пример: 10</i>", telegram.CreateBackToScheduleMenuKeyboard())

	// Сохраняем данные в состоянии
	newState := states.SetSetTrainingMaxParticipants(0)
	newState.Data["trackId"] = state.Data["trackId"]
	newState.Data["trainerId"] = state.Data["trainerId"]
	newState.Data["startTime"] = state.Data["startTime"]
	newState.Data["endTime"] = endTime
	return newState
}

func SetTrainingMaxParticipants(botUrl string, chatId int, update telegram.Update, repo database.ContentRepositoryInterface, state states.State) states.State {
	maxParticipantsStr := update.Message.Text

	// Валидируем введенное количество участников
	validator := validation.NewValidator()
	if result := validator.ValidateMaxParticipants(maxParticipantsStr); !result.IsValid {
		errorMsg := "❌ <b>Неверное количество участников</b>\n\n"
		for _, err := range result.Errors {
			errorMsg += fmt.Sprintf("• %s\n", err.Error())
		}
		errorMsg += "\n💡 <i>Пример: 10</i>"

		telegram.SendMessage(botUrl, chatId, errorMsg, telegram.CreateBackToScheduleMenuKeyboard())
		// Сохраняем данные из текущего состояния
		newState := states.SetSetTrainingMaxParticipants(0)
		newState.Data["trackId"] = state.Data["trackId"]
		newState.Data["trainerId"] = state.Data["trainerId"]
		newState.Data["startTime"] = state.Data["startTime"]
		newState.Data["endTime"] = state.Data["endTime"]
		return newState
	}

	maxParticipants, err := strconv.Atoi(maxParticipantsStr)
	if err != nil {
		telegram.SendMessage(botUrl, chatId, "❌ <b>Неверный формат числа</b>\n\n"+
			"Введите число участников:", telegram.CreateBackToScheduleMenuKeyboard())
		// Сохраняем данные из текущего состояния
		newState := states.SetSetTrainingMaxParticipants(0)
		newState.Data["trackId"] = state.Data["trackId"]
		newState.Data["trainerId"] = state.Data["trainerId"]
		newState.Data["startTime"] = state.Data["startTime"]
		newState.Data["endTime"] = state.Data["endTime"]
		return newState
	}

	// Получаем данные из состояния
	trackId, ok1 := state.Data["trackId"].(uint)
	trainerId, ok2 := state.Data["trainerId"].(uint)
	startTime, ok3 := state.Data["startTime"].(string)
	endTime, ok4 := state.Data["endTime"].(string)

	if !ok1 || !ok2 || !ok3 || !ok4 {
		logger.AdminError(chatId, "Неверные типы данных в состоянии для создания тренировки")
		return states.SetError()
	}

	// Создаем временные данные для тренировки
	tempData := &states.TempTrainingData{
		TrackID:         trackId,
		TrainerID:       trainerId,
		StartTime:       startTime,
		EndTime:         endTime,
		MaxParticipants: maxParticipants,
	}

	// Получаем информацию о тренере и трассе для отображения
	trainer, _ := repo.GetTrainerByID(trainerId)
	track, _ := repo.GetTrackByID(trackId)

	trainerName := "Неизвестный тренер"
	if trainer != nil {
		trainerName = trainer.Name
	}

	trackName := "Неизвестная трасса"
	if track != nil {
		trackName = track.Name
	}

	message := fmt.Sprintf("📅 <b>Подтверждение создания тренировки</b>\n\n"+
		"👨‍🏫 <b>Тренер:</b> %s\n"+
		"🏁 <b>Трасса:</b> %s\n"+
		"🕐 <b>Начало:</b> %s\n"+
		"🕕 <b>Окончание:</b> %s\n"+
		"👥 <b>Макс. участников:</b> %d\n\n"+
		"❓ <b>Создать тренировку?</b>",
		trainerName, trackName, startTime, endTime, maxParticipants)

	telegram.SendMessage(botUrl, chatId, message, telegram.CreateConfirmationKeyboard())
	return states.SetConfirmTrainingCreation().SetTempTrainingData(tempData)
}

func ConfirmTrainingCreation(botUrl string, chatId int, messageId int, repo database.ContentRepositoryInterface, tempData *states.TempTrainingData) states.State {
	// Парсим время начала и окончания
	startTime, err := time.Parse("2006-01-02 15:04", tempData.StartTime)
	if err != nil {
		logger.AdminError(chatId, "Парсинг времени начала: %v", err)
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Ошибка создания тренировки</b>\n\n"+
			"Неверный формат времени начала.\n"+
			"Обратитесь к администратору.", telegram.CreateBackToScheduleMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	endTime, err := time.Parse("2006-01-02 15:04", tempData.EndTime)
	if err != nil {
		logger.AdminError(chatId, "Парсинг времени окончания: %v", err)
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Ошибка создания тренировки</b>\n\n"+
			"Неверный формат времени окончания.\n"+
			"Обратитесь к администратору.", telegram.CreateBackToScheduleMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	// Проверяем, что время начала не в прошлом
	if startTime.Before(time.Now()) {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Ошибка создания тренировки</b>\n\n"+
			"Время начала тренировки не может быть в прошлом.\n"+
			"Выберите будущую дату и время.", telegram.CreateBackToScheduleMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	// Проверяем, что время окончания после времени начала
	if endTime.Before(startTime) || endTime.Equal(startTime) {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Ошибка создания тренировки</b>\n\n"+
			"Время окончания должно быть после времени начала.\n"+
			"Проверьте введенные данные.", telegram.CreateBackToScheduleMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	training := &database.Training{
		TrainerID:       tempData.TrainerID,
		TrackID:         tempData.TrackID,
		StartTime:       startTime,
		EndTime:         endTime,
		MaxParticipants: tempData.MaxParticipants,
		IsActive:        true,
	}

	_, err = repo.CreateTraining(training)
	if err != nil {
		logger.AdminError(chatId, "Создание тренировки: %v", err)
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Ошибка создания тренировки</b>\n\n"+
			"Ошибка сохранения.\n"+
			"Обратитесь к администратору.", telegram.CreateBackToScheduleMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	logger.AdminInfo(chatId, "Тренировка создана: %d", training.ID)
	telegram.EditMessage(botUrl, chatId, messageId, "✅ <b>Тренировка создана!</b>\n\n"+
		"🕐 Начало: "+training.StartTime.Format("2006-01-02 15:04")+"\n"+
		"🕕 Окончание: "+training.EndTime.Format("2006-01-02 15:04"), telegram.CreateBackToScheduleMenuKeyboard())
	return states.SetAdminKeyboard()
}

func EditTrackName(botUrl string, chatId int, messageId int, trackId uint) states.State {
	telegram.EditMessage(botUrl, chatId, messageId, "✏️ <b>Редактирование названия трассы</b>\n\n"+
		"📝 Введите новое название трассы:\n\n"+
		"💡 <i>Пример: Трасса №1 - Легкая</i>", telegram.CreateBackToTracksMenuKeyboard())
	return states.SetEditTrackName(trackId)
}

func SetEditTrackName(botUrl string, chatId int, update telegram.Update, repo database.ContentRepositoryInterface, trackId uint) states.State {
	name := update.Message.Text
	logger.AdminInfo(chatId, "Обновление трека %d: %s", trackId, name)

	// Получаем существующую трассу
	track, err := repo.GetTrackByID(trackId)
	if err != nil {
		logger.AdminError(chatId, "Получение трека %d: %v", trackId, err)
		telegram.SendMessage(botUrl, chatId, "❌ <b>Трасса не найдена</b>\n\n"+
			"🔍 Возможно, трасса была удалена.", telegram.CreateBackToTracksMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	// Обновляем только название
	track.Name = name
	err = repo.UpdateTrack(trackId, track)
	if err != nil {
		logger.AdminError(chatId, "Обновление трека %d: %v", trackId, err)
		telegram.SendMessage(botUrl, chatId, "❌ <b>Ошибка обновления названия трассы</b>\n\n"+
			"Ошибка сохранения.\n"+
			"Обратитесь к администратору.", telegram.CreateBackToTracksMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	logger.AdminInfo(chatId, "Трек %d обновлен: %s", trackId, name)
	telegram.SendMessage(botUrl, chatId, "✅ <b>Название трассы обновлено!</b>\n\n"+
		"🏁 Новое название: "+name, telegram.CreateBackToTracksMenuKeyboard())
	return states.SetAdminKeyboard()
}

func EditTrackInfo(botUrl string, chatId int, messageId int, trackId uint) states.State {
	telegram.EditMessage(botUrl, chatId, messageId, "✏️ <b>Редактирование описания трассы</b>\n\n"+
		"📋 Введите новое описание трассы:\n\n"+
		"💡 <i>Пример: Легкая трасса для начинающих, длина 1 км</i>", telegram.CreateBackToTracksMenuKeyboard())
	return states.SetEditTrackInfo(trackId)
}

func SetEditTrackInfo(botUrl string, chatId int, update telegram.Update, repo database.ContentRepositoryInterface, trackId uint) states.State {
	info := update.Message.Text

	// Получаем существующую трассу
	track, err := repo.GetTrackByID(trackId)
	if err != nil {
		logger.AdminError(chatId, "Получение трека %d: %v", trackId, err)
		telegram.SendMessage(botUrl, chatId, "❌ <b>Трасса не найдена</b>\n\n"+
			"🔍 Возможно, трасса была удалена.", telegram.CreateBackToTracksMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	// Обновляем только описание
	track.Info = info
	err = repo.UpdateTrack(trackId, track)
	if err != nil {
		telegram.SendMessage(botUrl, chatId, "❌ <b>Ошибка обновления описания трассы</b>\n\n"+
			"Ошибка сохранения.\n"+
			"Обратитесь к администратору.", telegram.CreateBackToTracksMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	telegram.SendMessage(botUrl, chatId, "✅ <b>Описание трассы обновлено!</b>\n\n"+
		"📄 Новое описание: "+info, telegram.CreateBackToTracksMenuKeyboard())
	return states.SetAdminKeyboard()
}

func ConfirmTrackDeletion(botUrl string, chatId int, messageId int, trackId uint, repo database.ContentRepositoryInterface) states.State {
	track, err := repo.GetTrackByID(trackId)
	if err != nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ Трасса не найдена.", telegram.CreateBackToTracksMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	message := fmt.Sprintf("⚠️ <b>Подтверждение удаления трассы</b>\n\n"+
		"🏁 <b>Трасса:</b> %s\n"+
		"📄 <b>Описание:</b> %s\n\n"+
		"🚨 <b>ВНИМАНИЕ!</b> Это действие нельзя отменить!\n\n"+
		"❓ <b>Вы уверены, что хотите удалить эту трассу?</b>",
		track.Name, track.Info)

	telegram.EditMessage(botUrl, chatId, messageId, message, telegram.CreateTrackDeletionConfirmationKeyboard(trackId))
	return states.SetConfirmTrackDelete(trackId)
}

func ExecuteTrackDeletion(botUrl string, chatId int, messageId int, trackId uint, repo database.ContentRepositoryInterface) states.State {
	track, err := repo.GetTrackByID(trackId)
	if err != nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Трасса не найдена</b>\n\n"+
			"🔍 Возможно, трасса уже была удалена.", telegram.CreateBackToTracksMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	err = repo.DeleteTrack(trackId)
	if err != nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Ошибка удаления трассы</b>\n\n"+
			"Попробуйте позже.", telegram.CreateBackToTracksMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	telegram.EditMessage(botUrl, chatId, messageId, fmt.Sprintf("🗑️ <b>Трасса %s удалена</b>", track.Name), telegram.CreateBackToTracksMenuKeyboard())
	return states.SetAdminKeyboard()
}

func EditTraining(botUrl string, chatId int, messageId int, trainingId uint, repo database.ContentRepositoryInterface) states.State {
	training, err := repo.GetTrainingById(trainingId)
	if err != nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Тренировка не найдена</b>\n\n"+
			"🔍 Возможно, тренировка была удалена.", telegram.CreateBackToScheduleMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	message := fmt.Sprintf("✏️ <b>Редактирование тренировки</b>\n\n"+
		"📅 <b>Дата:</b> %s\n"+
		"👥 <b>Макс. участников:</b> %d\n"+
		"🔄 <b>Статус:</b> %s\n\n"+
		"🎯 <b>Доступные действия:</b>",
		training.StartTime.Format("2006-01-02 15:04"), training.MaxParticipants,
		map[bool]string{true: "Активна", false: "Неактивна"}[training.IsActive])

	telegram.EditMessage(botUrl, chatId, messageId, message, telegram.CreateTrainingEditKeyboard(trainingId))
	return states.SetAdminKeyboard()
}

func ToggleTrainingStatus(botUrl string, chatId int, messageId int, trainingId uint, repo database.ContentRepositoryInterface) states.State {
	training, err := repo.GetTrainingById(trainingId)
	if err != nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Тренировка не найдена</b>\n\n"+
			"🔍 Возможно, тренировка была удалена.", telegram.CreateBackToScheduleMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	training.IsActive = !training.IsActive
	err = repo.UpdateTraining(trainingId, training)
	if err != nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Ошибка обновления статуса</b>\n\n"+
			"Попробуйте позже.", telegram.CreateBackToScheduleMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	status := map[bool]string{true: "активна", false: "неактивна"}[training.IsActive]
	telegram.EditMessage(botUrl, chatId, messageId, fmt.Sprintf("✅ <b>Тренировка %s</b>\n\n"+
		"📅 Дата: %s", status, training.StartTime.Format("2006-01-02 15:04")), telegram.CreateBackToScheduleMenuKeyboard())
	return states.SetAdminKeyboard()
}

func ConfirmTrainingDeletion(botUrl string, chatId int, messageId int, trainingId uint, repo database.ContentRepositoryInterface) states.State {
	training, err := repo.GetTrainingById(trainingId)
	if err != nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Тренировка не найдена</b>\n\n"+
			"🔍 Возможно, тренировка уже была удалена.", telegram.CreateBackToScheduleMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	// Получаем информацию о тренере и трассе
	trainer, _ := repo.GetTrainerByID(training.TrainerID)
	track, _ := repo.GetTrackByID(training.TrackID)

	trainerName := "Неизвестный тренер"
	if trainer != nil {
		trainerName = trainer.Name
	}

	trackName := "Неизвестная трасса"
	if track != nil {
		trackName = track.Name
	}

	message := fmt.Sprintf("⚠️ <b>Подтверждение удаления тренировки</b>\n\n"+
		"📅 <b>Дата и время:</b> %s\n"+
		"👨‍🏫 <b>Тренер:</b> %s\n"+
		"🏁 <b>Трасса:</b> %s\n"+
		"👥 <b>Макс. участников:</b> %d\n"+
		"🔄 <b>Статус:</b> %s\n\n"+
		"🚨 <b>ВНИМАНИЕ!</b> Это действие нельзя отменить!\n\n"+
		"❓ <b>Вы уверены, что хотите удалить эту тренировку?</b>",
		training.StartTime.Format("2006-01-02 15:04"), trainerName, trackName, training.MaxParticipants,
		map[bool]string{true: "Активна", false: "Неактивна"}[training.IsActive])

	telegram.EditMessage(botUrl, chatId, messageId, message, telegram.CreateTrainingDeletionConfirmationKeyboard(trainingId))
	return states.SetConfirmTrainingDelete(trainingId)
}

func ExecuteTrainingDeletion(botUrl string, chatId int, messageId int, trainingId uint, repo database.ContentRepositoryInterface) states.State {
	training, err := repo.GetTrainingById(trainingId)
	if err != nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Тренировка не найдена</b>\n\n"+
			"🔍 Возможно, тренировка уже была удалена.", telegram.CreateBackToScheduleMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	err = repo.DeleteTraining(trainingId)
	if err != nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Ошибка удаления тренировки</b>\n\n"+
			"Попробуйте позже.", telegram.CreateBackToScheduleMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	telegram.EditMessage(botUrl, chatId, messageId, fmt.Sprintf("🗑️ <b>Тренировка удалена</b>\n\n"+
		"📅 Дата: %s", training.StartTime.Format("2006-01-02 15:04")), telegram.CreateBackToScheduleMenuKeyboard())
	return states.SetAdminKeyboard()
}

func formatTracksListForAdmin(tracks []database.Track) string {
	if len(tracks) == 0 {
		return "📭 Трассы не найдены"
	}

	var builder strings.Builder
	for i, track := range tracks {
		builder.WriteString(fmt.Sprintf("%d. 🏁 <b>%s</b>\n", i+1, track.Name))
		builder.WriteString(fmt.Sprintf("   📄 %s\n\n", track.Info))
	}

	return builder.String()
}

func formatTrainingsListForAdmin(trainings []database.Training, repo database.ContentRepositoryInterface) string {
	if len(trainings) == 0 {
		return "📭 Тренировки не найдены"
	}

	var builder strings.Builder
	for i, training := range trainings {
		// Получаем информацию о тренере
		trainer, err := repo.GetTrainerByID(training.TrainerID)
		trainerName := "❓ Неизвестный"
		if err == nil && trainer != nil {
			trainerName = trainer.Name
		}

		// Получаем информацию о трассе
		track, err := repo.GetTrackByID(training.TrackID)
		trackName := "❓ Неизвестная"
		if err == nil && track != nil {
			trackName = track.Name
		}

		// Определяем статус и иконку
		statusIcon := "🟢"
		if !training.IsActive {
			statusIcon = "🔴"
		}

		// Форматируем дату и время
		dateStr := training.StartTime.Format("02.01")
		startTimeStr := training.StartTime.Format("15:04")
		endTimeStr := training.EndTime.Format("15:04")

		// Создаем компактную запись
		builder.WriteString(fmt.Sprintf("%d. %s <b>%s %s-%s</b>\n",
			i+1, statusIcon, dateStr, startTimeStr, endTimeStr))
		builder.WriteString(fmt.Sprintf("   👨‍🏫 %s | 🏁 %s | 👥 %d\n\n",
			trainerName, trackName, training.MaxParticipants))
	}

	return builder.String()
}

// ViewTrainingRequests - просмотр запросов тренировок
func ViewTrainingRequests(botUrl string, chatId int, messageId int, repo database.ContentRepositoryInterface) states.State {
	requests, err := repo.GetUnreviewedTrainingRequests()
	if err != nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Ошибка загрузки запросов</b>\n\n"+
			"Попробуйте позже.", telegram.CreateBackToAdminKeyboard())
		return states.SetAdminKeyboard()
	}

	if len(requests) == 0 {
		telegram.EditMessage(botUrl, chatId, messageId, "📭 <b>Новых запросов нет</b>\n\n"+
			"Все запросы рассмотрены.", telegram.CreateBackToAdminKeyboard())
		return states.SetAdminKeyboard()
	}

	message := "💬 <b>Запросы тренировок</b>\n\n"
	message += formatTrainingRequestsList(requests, repo)

	telegram.EditMessage(botUrl, chatId, messageId, message, telegram.CreateTrainingRequestsKeyboard(requests))
	return states.SetAdminKeyboard()
}

// MarkTrainingRequestAsReviewed - отметить запрос как рассмотренный
func MarkTrainingRequestAsReviewed(botUrl string, chatId int, messageId int, requestId uint, repo database.ContentRepositoryInterface) states.State {
	request, err := repo.GetTrainingRequestByID(requestId)
	if err != nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Запрос не найден</b>\n\n"+
			"🔍 Возможно, запрос уже был удален.", telegram.CreateBackToAdminKeyboard())
		return states.SetAdminKeyboard()
	}

	request.IsReviewed = true
	err = repo.UpdateTrainingRequest(requestId, request)
	if err != nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Ошибка обновления запроса</b>\n\n"+
			"Попробуйте позже.", telegram.CreateBackToAdminKeyboard())
		return states.SetAdminKeyboard()
	}

	telegram.EditMessage(botUrl, chatId, messageId, "✅ <b>Запрос отмечен как рассмотренный</b>\n\n"+
		"📝 Запрос больше не будет отображаться в очереди.", telegram.CreateBackToAdminKeyboard())
	return states.SetAdminKeyboard()
}

// formatTrainingRequestsList - форматирование списка запросов
func formatTrainingRequestsList(requests []database.TrainingRequest, repo database.ContentRepositoryInterface) string {
	if len(requests) == 0 {
		return "📭 Запросы не найдены"
	}

	var builder strings.Builder
	for i, request := range requests {
		// Получаем информацию о пользователе
		user, err := repo.GetUserByID(request.UserID)
		userName := "❓ Неизвестный"
		if err == nil && user != nil {
			userName = user.Name
		}

		// Форматируем дату
		dateStr := request.CreatedAt.Format("02.01 15:04")

		builder.WriteString(fmt.Sprintf("%d. 👤 <b>%s</b> (%s)\n",
			i+1, userName, dateStr))
		builder.WriteString(fmt.Sprintf("💬 %s\n\n", request.Message))
	}

	return builder.String()
}
