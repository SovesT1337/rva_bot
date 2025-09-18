package commands

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"x.localhost/rvabot/internal/database"
	"x.localhost/rvabot/internal/states"
	"x.localhost/rvabot/internal/telegram"
)

func Admin(botUrl string, chatId int, repo database.ContentRepositoryInterface) states.State {
	if !IsAdmin(chatId, repo) {
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

func IsAdmin(chatId int, repo database.ContentRepositoryInterface) bool {
	admin, err := repo.GetAdminByChatId(chatId)
	if err != nil {
		log.Printf("Admin check failed for user %d: %v", chatId, err)
		return false
	}
	if admin == nil {
		log.Printf("User %d not found in admins database", chatId)
		return false
	}

	log.Printf("User %d (%s) has admin rights", chatId, admin.Name)
	return true
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

	trainerChatId, err := strconv.Atoi(chatIdStr)
	if err != nil {
		telegram.SendMessage(botUrl, chatId, "❌ <b>Неверный формат Chat ID</b>\n\n"+
			"💡 <b>Введите числовой ID чата:</b>\n"+
			"📱 <i>Пример: 123456789</i>\n\n"+
			"🔄 Попробуйте еще раз:", telegram.CreateCancelKeyboard())
		return state
	}

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
	log.Printf("Creating trainer: %s (TgId: %s, ChatId: %d)", tempData.Name, tempData.TgId, tempData.ChatId)

	trainer := &database.Trainer{
		Name:   tempData.Name,
		TgId:   tempData.TgId,
		ChatId: tempData.ChatId,
		Info:   tempData.Info,
	}

	_, err := repo.CreateTrainer(trainer)
	if err != nil {
		log.Printf("ERROR: Failed to create trainer %s: %v", tempData.Name, err)
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Ошибка при создании тренера</b>\n"+
			"Попробуйте позже.", telegram.CreateBackToTrainersMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	log.Printf("Trainer created successfully: %s", tempData.Name)
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

func EditTrainer(botUrl string, chatId int, messageId int, repo database.ContentRepositoryInterface) states.State {
	trainers, err := repo.GetTrainers()
	if err != nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Ошибка получения списка тренеров</b>\n\n"+
			"Попробуйте позже.", telegram.CreateBackToTrainersMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	if len(trainers) == 0 {
		telegram.EditMessage(botUrl, chatId, messageId, "📝 <b>Список тренеров пуст</b>\n\n"+
			"👨‍🏫 Сначала добавьте тренеров через меню создания.", telegram.CreateBackToTrainersMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	message := "✏️ <b>Редактирование тренеров</b>\n\n" +
		"👨‍🏫 Тренеры:\n\n"
	for i, trainer := range trainers {
		message += fmt.Sprintf("%d. %s\n", i+1, trainer.Name)
	}

	telegram.EditMessage(botUrl, chatId, messageId, message, telegram.CreateTrainerSelectionKeyboard(trainers))
	return states.SetSelectTrainerToEdit()
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

func DeleteTrainer(botUrl string, chatId int, messageId int, repo database.ContentRepositoryInterface) states.State {
	trainers, err := repo.GetTrainers()
	if err != nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Ошибка получения списка тренеров</b>\n\n"+
			"Попробуйте позже.", telegram.CreateBackToTrainersMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	if len(trainers) == 0 {
		telegram.EditMessage(botUrl, chatId, messageId, "📝 <b>Список тренеров пуст</b>\n\n"+
			"👨‍🏫 Нет тренеров для удаления.", telegram.CreateBackToTrainersMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	message := "🗑️ <b>Удаление тренера</b>\n\n" +
		"⚠️ <b>Внимание!</b> Это действие нельзя отменить.\n\n" +
		"👨‍🏫 Тренеры:\n\n"
	for i, trainer := range trainers {
		message += fmt.Sprintf("%d. %s\n", i+1, trainer.Name)
	}

	telegram.EditMessage(botUrl, chatId, messageId, message, telegram.CreateTrainerDeletionKeyboard(trainers))
	return states.SetSelectTrainerToEdit()
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

func SelectTrainerToEdit(botUrl string, chatId int, messageId int, trainerId uint, repo database.ContentRepositoryInterface) states.State {
	trainer, err := repo.GetTrainerByID(trainerId)
	if err != nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Тренер не найден</b>\n\n"+
			"🔍 Возможно, тренер был удален.", telegram.CreateBackToTrainersMenuKeyboard())
		return states.SetAdminKeyboard()
	}
	if trainer == nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Тренер не найден</b>\n\n"+
			"🔍 Возможно, тренер был удален.", telegram.CreateBackToTrainersMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	message := fmt.Sprintf("✏️ <b>Редактирование тренера</b>\n\n"+
		"👤 <b>Тренер:</b> %s\n\n"+
		"📋 <b>Текущие данные:</b>\n"+
		"📝 <b>ФИО:</b> %s\n"+
		"📱 <b>Telegram ID:</b> %s\n"+
		"📄 <b>Информация:</b> %s\n\n"+
		"🎯 <b>Поля для редактирования:</b>",
		trainer.Name, trainer.Name, trainer.TgId, trainer.Info)

	telegram.EditMessage(botUrl, chatId, messageId, message, telegram.CreateTrainerEditKeyboard(trainerId))
	return states.SetAdminKeyboard()
}

func EditTrainerName(botUrl string, chatId int, messageId int, trainerId uint) states.State {
	telegram.EditMessage(botUrl, chatId, messageId, "✏️ <b>Редактирование ФИО тренера</b>\n\n"+
		"📝 Введите новое ФИО тренера:\n\n"+
		"💡 <i>Пример: Иванов Иван Иванович</i>", telegram.CreateBackToTrainersMenuKeyboard())
	return states.SetEditTrainerName(trainerId)
}

func SetEditTrainerName(botUrl string, chatId int, update telegram.Update, repo database.ContentRepositoryInterface, trainerId uint) states.State {
	name := update.Message.Text
	log.Printf("User %d updating trainer %d name to: %s", chatId, trainerId, name)

	trainer := &database.Trainer{Name: name}
	err := repo.UpdateTrainer(trainerId, trainer)
	if err != nil {
		log.Printf("ERROR: Failed to update trainer %d name: %v", trainerId, err)
		telegram.SendMessage(botUrl, chatId, "❌ <b>Ошибка обновления имени тренера</b>\n\n"+
			"Ошибка сохранения.\n"+
			"Обратитесь к администратору.", telegram.CreateBackToTrainersMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	log.Printf("Trainer %d name updated to: %s", trainerId, name)
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

	trainer := &database.Trainer{TgId: tgId}
	err := repo.UpdateTrainer(trainerId, trainer)
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

	trainer := &database.Trainer{Info: info}
	err := repo.UpdateTrainer(trainerId, trainer)
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
	log.Printf("User %d setting track name: %s", chatId, name)

	tempData := &states.TempTrackData{Name: name}
	newState := states.SetEnterTrackInfo(0).SetTempTrackData(tempData)

	telegram.SendMessage(botUrl, chatId, "📋 Введите описание трассы:\n\n"+
		"💡 <i>Пример: Легкая трасса для начинающих, длина 1 км</i>", telegram.CreateBackToTracksMenuKeyboard())
	return newState
}

func SetTrackInfo(botUrl string, chatId int, update telegram.Update, repo database.ContentRepositoryInterface, state states.State) states.State {
	info := update.Message.Text
	log.Printf("User %d setting track info: %s", chatId, info)

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
		log.Printf("ERROR: Failed to create track: %v", err)
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Ошибка создания трассы</b>\n\n"+
			"Ошибка сохранения.\n"+
			"Обратитесь к администратору.", telegram.CreateBackToTracksMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	log.Printf("Track created: %s", track.Name)
	telegram.EditMessage(botUrl, chatId, messageId, "✅ <b>Трасса создана!</b>\n\n"+
		"🏁 Название: "+track.Name, telegram.CreateBackToTracksMenuKeyboard())
	return states.SetAdminKeyboard()
}

func CancelTrackCreation(botUrl string, chatId int, messageId int) states.State {
	telegram.EditMessage(botUrl, chatId, messageId, "🚫 <b>Создание трассы отменено</b>\n\n"+
		"💡 Вы можете создать трассу позже.", telegram.CreateBackToTracksMenuKeyboard())
	return states.SetAdminKeyboard()
}

func EditTrack(botUrl string, chatId int, messageId int, repo database.ContentRepositoryInterface) states.State {
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

	message := "🏁 <b>Выберите трассу для редактирования:</b>\n\n"
	message += formatTracksListForAdmin(tracks)

	telegram.EditMessage(botUrl, chatId, messageId, message, telegram.CreateTrackSelectionKeyboard(tracks))
	return states.SetSelectTrackToEdit()
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

func DeleteTrack(botUrl string, chatId int, messageId int, repo database.ContentRepositoryInterface) states.State {
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

	message := "🗑️ <b>Выберите трассу для удаления:</b>\n\n"
	message += formatTracksListForAdmin(tracks)

	telegram.EditMessage(botUrl, chatId, messageId, message, telegram.CreateTrackDeletionKeyboard(tracks))
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

	telegram.EditMessage(botUrl, chatId, messageId, message, telegram.CreateTrainerSelectionKeyboard(trainers))
	return states.SetEnterTrainingTrainer(0)
}

func SetTrainingTrainer(botUrl string, chatId int, messageId int, trainerId uint, repo database.ContentRepositoryInterface, state states.State) states.State {
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

	telegram.EditMessage(botUrl, chatId, messageId, message, telegram.CreateTrackSelectionKeyboard(tracks))
	return states.SetEnterTrainingTrack(trainerId)
}

func SetTrainingTrack(botUrl string, chatId int, messageId int, trackId uint, repo database.ContentRepositoryInterface, state states.State) states.State {
	telegram.EditMessage(botUrl, chatId, messageId, "📅 Введите дату и время тренировки:\n\n"+
		"💡 <i>Пример: 2024-01-15 18:00</i>", telegram.CreateBackToScheduleMenuKeyboard())
	return states.SetEnterTrainingDate(0)
}

func SetTrainingDate(botUrl string, chatId int, update telegram.Update, repo database.ContentRepositoryInterface, state states.State) states.State {
	date := update.Message.Text
	log.Printf("User %d setting training date: %s", chatId, date)

	telegram.SendMessage(botUrl, chatId, "👥 Введите максимальное количество участников:\n\n"+
		"💡 <i>Пример: 10</i>", telegram.CreateBackToScheduleMenuKeyboard())
	return states.SetEnterTrainingMaxParticipants(0)
}

func SetTrainingMaxParticipants(botUrl string, chatId int, update telegram.Update, repo database.ContentRepositoryInterface, state states.State) states.State {
	maxParticipantsStr := update.Message.Text
	maxParticipants, err := strconv.Atoi(maxParticipantsStr)
	if err != nil {
		telegram.SendMessage(botUrl, chatId, "❌ <b>Неверный формат числа</b>\n\n"+
			"Введите число участников:", telegram.CreateBackToScheduleMenuKeyboard())
		return states.SetEnterTrainingMaxParticipants(0)
	}

	date := state.Data["date"].(string)
	trainerId := state.Data["trainerId"].(uint)
	trackId := state.Data["trackId"].(uint)

	message := fmt.Sprintf("📅 <b>Подтверждение создания тренировки</b>\n\n"+
		"👨‍🏫 <b>Тренер:</b> ID %d\n"+
		"🏁 <b>Трасса:</b> ID %d\n"+
		"📅 <b>Дата:</b> %s\n"+
		"👥 <b>Макс. участников:</b> %d\n\n"+
		"❓ <b>Создать тренировку?</b>",
		trainerId, trackId, date, maxParticipants)

	telegram.SendMessage(botUrl, chatId, message, telegram.CreateConfirmationKeyboard())
	return states.SetConfirmTrainingCreation()
}

func ConfirmTrainingCreation(botUrl string, chatId int, messageId int, repo database.ContentRepositoryInterface, tempData *states.TempTrainingData) states.State {
	training := &database.Training{
		TrainerID:       tempData.TrainerID,
		TrackID:         tempData.TrackID,
		Time:            time.Now(),
		MaxParticipants: tempData.MaxParticipants,
		IsActive:        true,
	}

	_, err := repo.CreateTraining(training)
	if err != nil {
		log.Printf("ERROR: Failed to create training: %v", err)
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Ошибка создания тренировки</b>\n\n"+
			"Ошибка сохранения.\n"+
			"Обратитесь к администратору.", telegram.CreateBackToScheduleMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	log.Printf("Training created: ID %d", training.ID)
	telegram.EditMessage(botUrl, chatId, messageId, "✅ <b>Тренировка создана!</b>\n\n"+
		"📅 Дата: "+training.Time.Format("2006-01-02 15:04"), telegram.CreateBackToScheduleMenuKeyboard())
	return states.SetAdminKeyboard()
}

func SelectTrackToEdit(botUrl string, chatId int, messageId int, trackId uint, repo database.ContentRepositoryInterface) states.State {
	track, err := repo.GetTrackByID(trackId)
	if err != nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Трасса не найдена</b>\n\n"+
			"🔍 Возможно, трасса была удалена.", telegram.CreateBackToTracksMenuKeyboard())
		return states.SetAdminKeyboard()
	}
	if track == nil {
		telegram.EditMessage(botUrl, chatId, messageId, "❌ <b>Трасса не найдена</b>\n\n"+
			"🔍 Возможно, трасса была удалена.", telegram.CreateBackToTracksMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	message := fmt.Sprintf("✏️ <b>Редактирование трассы</b>\n\n"+
		"🏁 <b>Трасса:</b> %s\n\n"+
		"📋 <b>Текущие данные:</b>\n"+
		"📝 <b>Название:</b> %s\n"+
		"📄 <b>Описание:</b> %s\n\n"+
		"🎯 <b>Поля для редактирования:</b>",
		track.Name, track.Name, track.Info)

	telegram.EditMessage(botUrl, chatId, messageId, message, telegram.CreateTrackEditKeyboard(trackId))
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
	log.Printf("User %d updating track %d name to: %s", chatId, trackId, name)

	track := &database.Track{Name: name}
	err := repo.UpdateTrack(trackId, track)
	if err != nil {
		log.Printf("ERROR: Failed to update track %d name: %v", trackId, err)
		telegram.SendMessage(botUrl, chatId, "❌ <b>Ошибка обновления названия трассы</b>\n\n"+
			"Ошибка сохранения.\n"+
			"Обратитесь к администратору.", telegram.CreateBackToTracksMenuKeyboard())
		return states.SetAdminKeyboard()
	}

	log.Printf("Track %d name updated to: %s", trackId, name)
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

	track := &database.Track{Info: info}
	err := repo.UpdateTrack(trackId, track)
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
		training.Time.Format("2006-01-02 15:04"), training.MaxParticipants,
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
		"📅 Дата: %s", status, training.Time.Format("2006-01-02 15:04")), telegram.CreateBackToScheduleMenuKeyboard())
	return states.SetAdminKeyboard()
}

func DeleteTraining(botUrl string, chatId int, messageId int, trainingId uint, repo database.ContentRepositoryInterface) states.State {
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
		"📅 Дата: %s", training.Time.Format("2006-01-02 15:04")), telegram.CreateBackToScheduleMenuKeyboard())
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
		builder.WriteString(fmt.Sprintf("%d. 📅 <b>%s</b>\n", i+1, training.Time.Format("2006-01-02 15:04")))
		builder.WriteString(fmt.Sprintf("   👥 Участников: %d\n", training.MaxParticipants))
		builder.WriteString(fmt.Sprintf("   🔄 Статус: %s\n\n", map[bool]string{true: "Активна", false: "Неактивна"}[training.IsActive]))
	}

	return builder.String()
}
