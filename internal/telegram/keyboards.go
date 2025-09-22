package telegram

import (
	"fmt"

	"x.localhost/rvabot/internal/database"
)

func CreateBaseKeyboard() inlineKeyboardMarkup {
	return inlineKeyboardMarkup{
		InlineKeyboard: [][]inlineKeyboardButton{
			{
				{Text: "🏠 Главное меню", CallbackData: "start"},
			},
		},
	}
}

func CreateNavigationKeyboard() inlineKeyboardMarkup {
	return inlineKeyboardMarkup{
		InlineKeyboard: [][]inlineKeyboardButton{
			{
				{Text: "🏠 Главное меню", CallbackData: "start"},
				{Text: "❓ Помощь", CallbackData: "help"},
			},
		},
	}
}

func CreateBackToAdminKeyboard() inlineKeyboardMarkup {
	return inlineKeyboardMarkup{
		InlineKeyboard: [][]inlineKeyboardButton{
			{
				{Text: "🔙 Назад к админке", CallbackData: "admin"},
			},
		},
	}
}

func CreateBackToInfoKeyboard() inlineKeyboardMarkup {
	return inlineKeyboardMarkup{
		InlineKeyboard: [][]inlineKeyboardButton{
			{
				{Text: "🔙 Назад", CallbackData: "Info"},
			},
		},
	}
}

func CreateBackToTrainersMenuKeyboard() inlineKeyboardMarkup {
	return inlineKeyboardMarkup{
		InlineKeyboard: [][]inlineKeyboardButton{
			{
				{Text: "🔙 Назад к тренерам", CallbackData: "trainersMenu"},
			},
		},
	}
}

func CreateBackToTracksMenuKeyboard() inlineKeyboardMarkup {
	return inlineKeyboardMarkup{
		InlineKeyboard: [][]inlineKeyboardButton{
			{
				{Text: "🔙 Назад", CallbackData: "tracksMenu"},
			},
		},
	}
}

func CreateBackToScheduleMenuKeyboard() inlineKeyboardMarkup {
	return inlineKeyboardMarkup{
		InlineKeyboard: [][]inlineKeyboardButton{
			{
				{Text: "🔙 Назад", CallbackData: "scheduleMenu"},
			},
		},
	}
}

func CreateCancelKeyboard() inlineKeyboardMarkup {
	return inlineKeyboardMarkup{
		InlineKeyboard: [][]inlineKeyboardButton{
			{
				{Text: "❌ Отмена", CallbackData: "cancel"},
			},
		},
	}
}

func CreateStartKeyboard(chatId int, repo database.ContentRepositoryInterface) inlineKeyboardMarkup {
	keyboard := [][]inlineKeyboardButton{
		{
			{Text: "🏃‍♂️ Записаться на тренировку", CallbackData: "BookTraining"},
		},
		{
			{Text: "ℹ️ Информация о занятиях", CallbackData: "Info"},
		},
		{
			{Text: "🛒 Экипировка", URL: "https://dudarevmotorsport.ru/"},
		},
	}

	// Проверяем, является ли пользователь администратором
	if database.IsAdmin(chatId, repo) {
		keyboard = append(keyboard, []inlineKeyboardButton{
			{Text: "⚙️ Админ-панель", CallbackData: "admin"},
		})
	}

	return inlineKeyboardMarkup{InlineKeyboard: keyboard}
}

func CreateAdminKeyboard() inlineKeyboardMarkup {
	return inlineKeyboardMarkup{
		InlineKeyboard: [][]inlineKeyboardButton{
			{
				{Text: "👨‍🏫 Тренеры", CallbackData: "trainersMenu"},
				{Text: "🏁 Трассы", CallbackData: "tracksMenu"},
			},
			{
				{Text: "📅 Расписание", CallbackData: "scheduleMenu"},
			},
			{
				{Text: "🏠 Главное меню", CallbackData: "start"},
			},
		},
	}
}

func CreateTrainersListWithActionsKeyboard(trainers []database.Trainer) inlineKeyboardMarkup {
	var buttons [][]inlineKeyboardButton

	// Добавляем кнопку "Добавить тренера" в начале
	buttons = append(buttons, []inlineKeyboardButton{
		{Text: "➕ Добавить тренера", CallbackData: "createTrainer"},
	})

	// Добавляем кнопки для каждого тренера
	for _, trainer := range trainers {
		buttons = append(buttons, []inlineKeyboardButton{
			{Text: fmt.Sprintf("✏️ %s", trainer.Name), CallbackData: fmt.Sprintf("editTrainerName_%d", trainer.ID)},
			{Text: "📱", CallbackData: fmt.Sprintf("editTrainerTgId_%d", trainer.ID)},
			{Text: "📄", CallbackData: fmt.Sprintf("editTrainerInfo_%d", trainer.ID)},
			{Text: "🗑️", CallbackData: fmt.Sprintf("deleteTrainer_%d", trainer.ID)},
		})
	}

	// Добавляем кнопку "Назад к админке"
	buttons = append(buttons, []inlineKeyboardButton{
		{Text: "🔙 Назад к админке", CallbackData: "admin"},
	})

	return inlineKeyboardMarkup{InlineKeyboard: buttons}
}

func CreateTracksListWithActionsKeyboard(tracks []database.Track) inlineKeyboardMarkup {
	var buttons [][]inlineKeyboardButton

	// Добавляем кнопку "Добавить трассу" в начале
	buttons = append(buttons, []inlineKeyboardButton{
		{Text: "➕ Добавить трассу", CallbackData: "createTrack"},
	})

	// Добавляем кнопки для каждой трассы
	for _, track := range tracks {
		buttons = append(buttons, []inlineKeyboardButton{
			{Text: fmt.Sprintf("✏️ %s", track.Name), CallbackData: fmt.Sprintf("editTrackName_%d", track.ID)},
			{Text: "📄", CallbackData: fmt.Sprintf("editTrackInfo_%d", track.ID)},
			{Text: "🗑️", CallbackData: fmt.Sprintf("deleteTrack_%d", track.ID)},
		})
	}

	// Добавляем кнопку "Назад к админке"
	buttons = append(buttons, []inlineKeyboardButton{
		{Text: "🔙 Назад к админке", CallbackData: "admin"},
	})

	return inlineKeyboardMarkup{InlineKeyboard: buttons}
}

func CreateTrainingsListWithActionsKeyboard(trainings []database.Training) inlineKeyboardMarkup {
	var buttons [][]inlineKeyboardButton

	// Добавляем кнопку "Добавить тренировку" в начале
	buttons = append(buttons, []inlineKeyboardButton{
		{Text: "➕ Добавить тренировку", CallbackData: "createSchedule"},
	})

	// Добавляем кнопки для каждой тренировки
	for _, training := range trainings {
		statusIcon := "🟢"
		if !training.IsActive {
			statusIcon = "🔴"
		}
		buttons = append(buttons, []inlineKeyboardButton{
			{Text: fmt.Sprintf("%s %s", statusIcon, training.StartTime.Format("02.01 15:04")), CallbackData: fmt.Sprintf("editTraining_%d", training.ID)},
			{Text: "🗑️", CallbackData: fmt.Sprintf("deleteTraining_%d", training.ID)},
		})
	}

	// Добавляем кнопку "Назад к админке"
	buttons = append(buttons, []inlineKeyboardButton{
		{Text: "🔙 Назад к админке", CallbackData: "admin"},
	})

	return inlineKeyboardMarkup{InlineKeyboard: buttons}
}

func CreateInfoKeyboard() inlineKeyboardMarkup {
	return inlineKeyboardMarkup{
		InlineKeyboard: [][]inlineKeyboardButton{
			{
				{Text: "👨‍🏫 Тренерский состав", CallbackData: "infoTrainer"},
			},
			{
				{Text: "🏁 Трассы", CallbackData: "infoTrack"},
			},
			{
				{Text: "📅 Расписание тренировок", CallbackData: "viewScheduleUser"},
			},
			{
				{Text: "📋 Формат тренировок", CallbackData: "infoFormat"},
			},
			{
				{Text: "🏠 Главное меню", CallbackData: "start"},
			},
		},
	}
}

func CreateConfirmationKeyboard() inlineKeyboardMarkup {
	return inlineKeyboardMarkup{
		InlineKeyboard: [][]inlineKeyboardButton{
			{
				{Text: "✅ Подтвердить", CallbackData: "confirm"},
				{Text: "❌ Отменить", CallbackData: "cancel"},
			},
		},
	}
}

func CreateTrainerEditKeyboard(trainerId uint) inlineKeyboardMarkup {
	return inlineKeyboardMarkup{
		InlineKeyboard: [][]inlineKeyboardButton{
			{
				{Text: "✏️ ФИО", CallbackData: fmt.Sprintf("editTrainerName_%d", trainerId)},
				{Text: "📱 Telegram ID", CallbackData: fmt.Sprintf("editTrainerTgId_%d", trainerId)},
			},
			{
				{Text: "📄 Информация", CallbackData: fmt.Sprintf("editTrainerInfo_%d", trainerId)},
			},
			{
				{Text: "🔙 Назад к тренерам", CallbackData: "trainersMenu"},
			},
		},
	}
}

func CreateDeletionConfirmationKeyboard(trainerId uint) inlineKeyboardMarkup {
	return inlineKeyboardMarkup{
		InlineKeyboard: [][]inlineKeyboardButton{
			{
				{Text: "🗑️ Удалить", CallbackData: fmt.Sprintf("confirmDelete_%d", trainerId)},
				{Text: "❌ Отменить", CallbackData: "trainersMenu"},
			},
		},
	}
}

func CreateTrainingDeletionConfirmationKeyboard(trainingId uint) inlineKeyboardMarkup {
	return inlineKeyboardMarkup{
		InlineKeyboard: [][]inlineKeyboardButton{
			{
				{Text: "🗑️ Удалить", CallbackData: fmt.Sprintf("confirmDeleteTraining_%d", trainingId)},
				{Text: "❌ Отменить", CallbackData: "scheduleMenu"},
			},
		},
	}
}

func CreateTrackEditKeyboard(trackId uint) inlineKeyboardMarkup {
	return inlineKeyboardMarkup{
		InlineKeyboard: [][]inlineKeyboardButton{
			{
				{Text: "✏️ Название", CallbackData: fmt.Sprintf("editTrackName_%d", trackId)},
				{Text: "📄 Информация", CallbackData: fmt.Sprintf("editTrackInfo_%d", trackId)},
			},
			{
				{Text: "🔙 Назад к трассам", CallbackData: "tracksMenu"},
			},
		},
	}
}

func CreateTrackDeletionConfirmationKeyboard(trackId uint) inlineKeyboardMarkup {
	return inlineKeyboardMarkup{
		InlineKeyboard: [][]inlineKeyboardButton{
			{
				{Text: "🗑️ Удалить", CallbackData: fmt.Sprintf("confirmDeleteTrack_%d", trackId)},
				{Text: "❌ Отменить", CallbackData: "tracksMenu"},
			},
		},
	}
}

func CreateTrainingRegistrationConfirmationKeyboard(trainingId uint) inlineKeyboardMarkup {
	return inlineKeyboardMarkup{
		InlineKeyboard: [][]inlineKeyboardButton{
			{
				{Text: "✅ Записаться", CallbackData: fmt.Sprintf("confirmTrainingRegistration_%d", trainingId)},
				{Text: "❌ Отменить", CallbackData: "cancel"},
			},
		},
	}
}

func CreateTrainingApprovalKeyboard(registrationId uint) inlineKeyboardMarkup {
	return inlineKeyboardMarkup{
		InlineKeyboard: [][]inlineKeyboardButton{
			{
				{Text: "✅ Одобрить", CallbackData: fmt.Sprintf("approveRegistration_%d", registrationId)},
				{Text: "❌ Отклонить", CallbackData: fmt.Sprintf("rejectRegistration_%d", registrationId)},
			},
		},
	}
}

func CreateTrainerSelectionForTrainingKeyboard(trainers []database.Trainer) inlineKeyboardMarkup {
	var buttons [][]inlineKeyboardButton

	for _, t := range trainers {
		buttons = append(buttons, []inlineKeyboardButton{{
			Text:         t.Name,
			CallbackData: fmt.Sprintf("selectTrainerForTraining_%d", t.ID),
		}})
	}

	buttons = append(buttons, []inlineKeyboardButton{
		{Text: "🔙 Назад к расписанию", CallbackData: "scheduleMenu"},
		{Text: "❌ Отменить", CallbackData: "cancel"},
	})

	return inlineKeyboardMarkup{InlineKeyboard: buttons}
}

func CreateTrackSelectionForTrainingKeyboard(tracks []database.Track) inlineKeyboardMarkup {
	var buttons [][]inlineKeyboardButton

	for _, t := range tracks {
		buttons = append(buttons, []inlineKeyboardButton{{
			Text:         t.Name,
			CallbackData: fmt.Sprintf("selectTrackForTraining_%d", t.ID),
		}})
	}

	buttons = append(buttons, []inlineKeyboardButton{
		{Text: "🔙 Назад к расписанию", CallbackData: "scheduleMenu"},
		{Text: "❌ Отменить", CallbackData: "cancel"},
	})

	return inlineKeyboardMarkup{InlineKeyboard: buttons}
}

func CreateTrainingEditKeyboard(trainingId uint) inlineKeyboardMarkup {
	return inlineKeyboardMarkup{
		InlineKeyboard: [][]inlineKeyboardButton{
			{
				{Text: "📅 Изменить дату", CallbackData: fmt.Sprintf("editTrainingDate_%d", trainingId)},
				{Text: "👥 Изменить участников", CallbackData: fmt.Sprintf("editTrainingParticipants_%d", trainingId)},
			},
			{
				{Text: "🔄 Активировать/Деактивировать", CallbackData: fmt.Sprintf("toggleTrainingStatus_%d", trainingId)},
			},
			{
				{Text: "🗑️ Удалить тренировку", CallbackData: fmt.Sprintf("deleteTraining_%d", trainingId)},
			},
			{
				{Text: "🔙 Назад к расписанию", CallbackData: "scheduleMenu"},
			},
		},
	}
}

func CreateTrackSelectionForRegistrationKeyboard(tracks []database.Track) inlineKeyboardMarkup {
	var buttons [][]inlineKeyboardButton

	for _, t := range tracks {
		buttons = append(buttons, []inlineKeyboardButton{{
			Text:         t.Name,
			CallbackData: fmt.Sprintf("selectTrackForRegistration_%d", t.ID),
		}})
	}

	buttons = append(buttons, []inlineKeyboardButton{
		{Text: "🏠 Главное меню", CallbackData: "start"},
	})

	return inlineKeyboardMarkup{InlineKeyboard: buttons}
}

func CreateTrainerSelectionForRegistrationKeyboard(trainers []database.Trainer) inlineKeyboardMarkup {
	var buttons [][]inlineKeyboardButton

	for _, t := range trainers {
		buttons = append(buttons, []inlineKeyboardButton{{
			Text:         t.Name,
			CallbackData: fmt.Sprintf("selectTrainerForRegistration_%d", t.ID),
		}})
	}

	buttons = append(buttons, []inlineKeyboardButton{
		{Text: "🔙 Назад к выбору трассы", CallbackData: "backToTrackSelection"},
	})

	return inlineKeyboardMarkup{InlineKeyboard: buttons}
}

func CreateTrainingTimeSelectionKeyboard(trainings []database.Training) inlineKeyboardMarkup {
	var buttons [][]inlineKeyboardButton

	for _, t := range trainings {
		buttons = append(buttons, []inlineKeyboardButton{{
			Text:         t.StartTime.Format("02.01 15:04"),
			CallbackData: fmt.Sprintf("selectTrainingTimeForRegistration_%d", t.ID),
		}})
	}

	buttons = append(buttons, []inlineKeyboardButton{
		{Text: "🔙 Назад к выбору тренера", CallbackData: "backToTrainerSelection"},
	})

	return inlineKeyboardMarkup{InlineKeyboard: buttons}
}
