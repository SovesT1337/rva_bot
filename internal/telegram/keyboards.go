package telegram

import (
	"fmt"

	"x.localhost/rvabot/internal/database"
)

func CreateEmptyKeyboard() inlineKeyboardMarkup {
	return inlineKeyboardMarkup{
		InlineKeyboard: [][]inlineKeyboardButton{},
	}
}

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

func CreateStartKeyboard() inlineKeyboardMarkup {
	return inlineKeyboardMarkup{
		InlineKeyboard: [][]inlineKeyboardButton{
			{
				{Text: "🏃‍♂️ Записаться на тренировку", CallbackData: "BookTraining"},
			},
			{
				{Text: "ℹ️ Информация о занятиях", CallbackData: "Info"},
			},
			{
				{Text: "📊 Мой рейтинг ELO", CallbackData: "Raiting"},
			},
			{
				{Text: "🛒 Экипировка", URL: "https://dudarevmotorsport.ru/"},
			},
			{
				{Text: "⚙️ Админ-панель", CallbackData: "admin"},
			},
		},
	}
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

func CreateTrainersMenuKeyboard() inlineKeyboardMarkup {
	return inlineKeyboardMarkup{
		InlineKeyboard: [][]inlineKeyboardButton{
			{
				{Text: "➕ Добавить тренера", CallbackData: "createTrainer"},
				{Text: "👥 Список тренеров", CallbackData: "viewTrainers"},
			},
			{
				{Text: "✏️ Редактировать", CallbackData: "editTrainer"},
				{Text: "🗑️ Удалить", CallbackData: "deleteTrainer"},
			},
			{
				{Text: "🔙 Назад к админке", CallbackData: "admin"},
			},
		},
	}
}

func CreateTracksMenuKeyboard() inlineKeyboardMarkup {
	return inlineKeyboardMarkup{
		InlineKeyboard: [][]inlineKeyboardButton{
			{
				{Text: "➕ Добавить трассу", CallbackData: "createTrack"},
				{Text: "🏁 Список трасс", CallbackData: "viewTracks"},
			},
			{
				{Text: "✏️ Редактировать", CallbackData: "editTrack"},
				{Text: "🗑️ Удалить", CallbackData: "deleteTrack"},
			},
			{
				{Text: "🔙 Назад к админке", CallbackData: "admin"},
			},
		},
	}
}

func CreateScheduleMenuKeyboard() inlineKeyboardMarkup {
	return inlineKeyboardMarkup{
		InlineKeyboard: [][]inlineKeyboardButton{
			{
				{Text: "➕ Добавить в расписание", CallbackData: "createSchedule"},
			},
			{
				{Text: "📅 Просмотр расписания", CallbackData: "viewSchedule"},
			},
			{
				{Text: "✏️ Редактировать", CallbackData: "editSchedule"},
			},
			{
				{Text: "🔙 Назад к админке", CallbackData: "admin"},
			},
		},
	}
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

func CreateTrainerSelectionKeyboard(trainers []database.Trainer) inlineKeyboardMarkup {
	var buttons [][]inlineKeyboardButton

	for i := 0; i < len(trainers); i += 5 {
		var row []inlineKeyboardButton
		for j := i; j < i+5 && j < len(trainers); j++ {
			trainer := trainers[j]
			buttonText := fmt.Sprintf("%d. %s", j+1, trainer.Name)
			if len(buttonText) > 20 {
				buttonText = buttonText[:17] + "..."
			}
			row = append(row, inlineKeyboardButton{
				Text:         buttonText,
				CallbackData: fmt.Sprintf("selectTrainer_%d", trainer.ID),
			})
		}
		buttons = append(buttons, row)
	}

	buttons = append(buttons, []inlineKeyboardButton{
		{Text: "🔙 Назад к тренерам", CallbackData: "trainersMenu"},
	})

	return inlineKeyboardMarkup{InlineKeyboard: buttons}
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

func CreateTrainerDeletionKeyboard(trainers []database.Trainer) inlineKeyboardMarkup {
	var buttons [][]inlineKeyboardButton

	for i := 0; i < len(trainers); i += 3 {
		var row []inlineKeyboardButton
		for j := i; j < i+3 && j < len(trainers); j++ {
			trainer := trainers[j]
			buttonText := fmt.Sprintf("🗑️ %s", trainer.Name)
			if len(buttonText) > 20 {
				buttonText = buttonText[:17] + "..."
			}
			row = append(row, inlineKeyboardButton{
				Text:         buttonText,
				CallbackData: fmt.Sprintf("deleteTrainer_%d", trainer.ID),
			})
		}
		buttons = append(buttons, row)
	}

	buttons = append(buttons, []inlineKeyboardButton{
		{Text: "🔙 Назад к тренерам", CallbackData: "trainersMenu"},
	})

	return inlineKeyboardMarkup{InlineKeyboard: buttons}
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

func CreateTrackSelectionKeyboard(tracks []database.Track) inlineKeyboardMarkup {
	var buttons [][]inlineKeyboardButton

	for i := 0; i < len(tracks); i += 5 {
		var row []inlineKeyboardButton
		for j := i; j < i+5 && j < len(tracks); j++ {
			track := tracks[j]
			buttonText := fmt.Sprintf("%d. %s", j+1, track.Name)
			if len(buttonText) > 20 {
				buttonText = buttonText[:17] + "..."
			}
			row = append(row, inlineKeyboardButton{
				Text:         buttonText,
				CallbackData: fmt.Sprintf("selectTrack_%d", track.ID),
			})
		}
		buttons = append(buttons, row)
	}

	// Добавляем кнопку "Назад"
	buttons = append(buttons, []inlineKeyboardButton{
		{Text: "🔙 Назад к трассам", CallbackData: "tracksMenu"},
	})

	return inlineKeyboardMarkup{InlineKeyboard: buttons}
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

func CreateTrackDeletionKeyboard(tracks []database.Track) inlineKeyboardMarkup {
	var buttons [][]inlineKeyboardButton

	// Создаем кнопки для каждой трассы (максимум 3 на ряд для удаления)
	for i := 0; i < len(tracks); i += 3 {
		var row []inlineKeyboardButton
		for j := i; j < i+3 && j < len(tracks); j++ {
			track := tracks[j]
			buttonText := fmt.Sprintf("%d. %s", j+1, track.Name)
			if len(buttonText) > 15 {
				buttonText = buttonText[:12] + "..."
			}
			row = append(row, inlineKeyboardButton{
				Text:         buttonText,
				CallbackData: fmt.Sprintf("deleteTrack_%d", track.ID),
			})
		}
		buttons = append(buttons, row)
	}

	// Добавляем кнопку "Назад"
	buttons = append(buttons, []inlineKeyboardButton{
		{Text: "🔙 Назад к трассам", CallbackData: "tracksMenu"},
	})

	return inlineKeyboardMarkup{InlineKeyboard: buttons}
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

// Клавиатуры для регистрации на тренировки
func CreateTrainingSelectionKeyboard(trainings []database.Training) inlineKeyboardMarkup {
	var buttons [][]inlineKeyboardButton

	// Создаем кнопки для каждой тренировки (максимум 2 на ряд)
	for i := 0; i < len(trainings); i += 2 {
		var row []inlineKeyboardButton
		for j := i; j < i+2 && j < len(trainings); j++ {
			training := trainings[j]
			buttonText := fmt.Sprintf("🏃‍♂️ %d", j+1)
			row = append(row, inlineKeyboardButton{
				Text:         buttonText,
				CallbackData: fmt.Sprintf("selectTraining_%d", training.ID),
			})
		}
		buttons = append(buttons, row)
	}

	// Добавляем кнопку "Назад"
	buttons = append(buttons, []inlineKeyboardButton{
		{Text: "🏠 Главное меню", CallbackData: "start"},
	})

	return inlineKeyboardMarkup{InlineKeyboard: buttons}
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

// Клавиатуры для тренеров (подтверждение/отклонение заявок)
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

// Клавиатуры для создания тренировок
func CreateTrainerSelectionForTrainingKeyboard(trainers []database.Trainer) inlineKeyboardMarkup {
	var buttons [][]inlineKeyboardButton

	// Создаем кнопки для каждого тренера (максимум 3 на ряд)
	for i := 0; i < len(trainers); i += 3 {
		var row []inlineKeyboardButton
		for j := i; j < i+3 && j < len(trainers); j++ {
			trainer := trainers[j]
			buttonText := fmt.Sprintf("%d. %s", j+1, trainer.Name)
			if len(buttonText) > 20 {
				buttonText = buttonText[:17] + "..."
			}
			row = append(row, inlineKeyboardButton{
				Text:         buttonText,
				CallbackData: fmt.Sprintf("selectTrainerForTraining_%d", trainer.ID),
			})
		}
		buttons = append(buttons, row)
	}

	// Добавляем кнопки "Назад" и "Отменить"
	buttons = append(buttons, []inlineKeyboardButton{
		{Text: "🔙 Назад к расписанию", CallbackData: "scheduleMenu"},
		{Text: "❌ Отменить", CallbackData: "cancel"},
	})

	return inlineKeyboardMarkup{InlineKeyboard: buttons}
}

func CreateTrackSelectionForTrainingKeyboard(tracks []database.Track) inlineKeyboardMarkup {
	var buttons [][]inlineKeyboardButton

	// Создаем кнопки для каждой трассы (максимум 3 на ряд)
	for i := 0; i < len(tracks); i += 3 {
		var row []inlineKeyboardButton
		for j := i; j < i+3 && j < len(tracks); j++ {
			track := tracks[j]
			buttonText := fmt.Sprintf("%d. %s", j+1, track.Name)
			if len(buttonText) > 20 {
				buttonText = buttonText[:17] + "..."
			}
			row = append(row, inlineKeyboardButton{
				Text:         buttonText,
				CallbackData: fmt.Sprintf("selectTrackForTraining_%d", track.ID),
			})
		}
		buttons = append(buttons, row)
	}

	buttons = append(buttons, []inlineKeyboardButton{
		{Text: "🔙 Назад к расписанию", CallbackData: "scheduleMenu"},
		{Text: "❌ Отменить", CallbackData: "cancel"},
	})

	return inlineKeyboardMarkup{InlineKeyboard: buttons}
}

func CreateTrainingEditSelectionKeyboard(trainings []database.Training) inlineKeyboardMarkup {
	var buttons [][]inlineKeyboardButton

	for i := 0; i < len(trainings); i += 2 {
		var row []inlineKeyboardButton
		for j := i; j < i+2 && j < len(trainings); j++ {
			training := trainings[j]
			buttonText := fmt.Sprintf("🏃‍♂️ %d", j+1)
			row = append(row, inlineKeyboardButton{
				Text:         buttonText,
				CallbackData: fmt.Sprintf("editTraining_%d", training.ID),
			})
		}
		buttons = append(buttons, row)
	}

	buttons = append(buttons, []inlineKeyboardButton{
		{Text: "🔙 Назад к расписанию", CallbackData: "scheduleMenu"},
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

	for i := 0; i < len(tracks); i += 2 {
		var row []inlineKeyboardButton
		for j := i; j < i+2 && j < len(tracks); j++ {
			track := tracks[j]
			buttonText := fmt.Sprintf("🏁 %s", track.Name)
			if len(buttonText) > 20 {
				buttonText = buttonText[:17] + "..."
			}
			row = append(row, inlineKeyboardButton{
				Text:         buttonText,
				CallbackData: fmt.Sprintf("selectTrackForRegistration_%d", track.ID),
			})
		}
		buttons = append(buttons, row)
	}

	buttons = append(buttons, []inlineKeyboardButton{
		{Text: "🏠 Главное меню", CallbackData: "start"},
	})

	return inlineKeyboardMarkup{InlineKeyboard: buttons}
}

func CreateTrainerSelectionForRegistrationKeyboard(trainers []database.Trainer) inlineKeyboardMarkup {
	var buttons [][]inlineKeyboardButton

	for i := 0; i < len(trainers); i += 2 {
		var row []inlineKeyboardButton
		for j := i; j < i+2 && j < len(trainers); j++ {
			trainer := trainers[j]
			buttonText := fmt.Sprintf("👨‍🏫 %s", trainer.Name)
			if len(buttonText) > 20 {
				buttonText = buttonText[:17] + "..."
			}
			row = append(row, inlineKeyboardButton{
				Text:         buttonText,
				CallbackData: fmt.Sprintf("selectTrainerForRegistration_%d", trainer.ID),
			})
		}
		buttons = append(buttons, row)
	}

	buttons = append(buttons, []inlineKeyboardButton{
		{Text: "🔙 Назад к выбору трассы", CallbackData: "backToTrackSelection"},
	})

	return inlineKeyboardMarkup{InlineKeyboard: buttons}
}

func CreateTrainingTimeSelectionKeyboard(trainings []database.Training) inlineKeyboardMarkup {
	var buttons [][]inlineKeyboardButton

	for _, training := range trainings {
		buttonText := fmt.Sprintf("📅 %s", training.Time.Format("02.01 15:04"))
		row := []inlineKeyboardButton{
			{
				Text:         buttonText,
				CallbackData: fmt.Sprintf("selectTrainingTimeForRegistration_%d", training.ID),
			},
		}
		buttons = append(buttons, row)
	}

	buttons = append(buttons, []inlineKeyboardButton{
		{Text: "🔙 Назад к выбору тренера", CallbackData: "backToTrainerSelection"},
	})

	return inlineKeyboardMarkup{InlineKeyboard: buttons}
}
