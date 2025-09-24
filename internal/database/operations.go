package database

import (
	"context"
	"errors"
	"time"

	"x.localhost/rvabot/internal/logger"

	"gorm.io/gorm"
)

// ContentRepository уже определен в repository.go

func (r *ContentRepository) CreateTraining(content *Training) (uint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.DatabaseInfo("Создание тренировки: TrainerID=%d, TrackID=%d", content.TrainerID, content.TrackID)

	result := r.db.WithContext(ctx).Create(content)
	if result.Error != nil {
		logger.DatabaseError("Создание тренировки: %v", result.Error)
		return 0, result.Error
	}

	logger.DatabaseInfo("Тренировка создана: %d", content.ID)
	return content.ID, nil
}

func (r *ContentRepository) GetTrainingById(id uint) (*Training, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var training Training
	result := r.db.WithContext(ctx).First(&training, id)
	if result.Error != nil {
		logger.DatabaseError("Получение тренировки %d: %v", id, result.Error)
		return nil, result.Error
	}

	return &training, nil
}

func (r *ContentRepository) GetTrainings() ([]Training, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var trainings []Training
	result := r.db.WithContext(ctx).Find(&trainings)
	if result.Error != nil {
		logger.DatabaseError("Получение тренировок: %v", result.Error)
		return nil, result.Error
	}

	return trainings, nil
}

func (r *ContentRepository) GetActiveTrainings() ([]Training, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var trainings []Training
	result := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&trainings)
	if result.Error != nil {
		logger.DatabaseError("Получение активных тренировок: %v", result.Error)
		return nil, result.Error
	}

	return trainings, nil
}

func (r *ContentRepository) UpdateTraining(id uint, training *Training) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := r.db.WithContext(ctx).Model(&Training{}).Where("id = ?", id).Updates(training)
	if result.Error != nil {
		logger.DatabaseError("Обновление тренировки %d: %v", id, result.Error)
		return result.Error
	}

	logger.DatabaseInfo("Тренировка обновлена: %d", id)
	return nil
}

func (r *ContentRepository) DeleteTraining(id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := r.db.WithContext(ctx).Delete(&Training{}, id)
	if result.Error != nil {
		logger.DatabaseError("Удаление тренировки %d: %v", id, result.Error)
		return result.Error
	}

	logger.DatabaseInfo("Тренировка удалена: %d", id)
	return nil
}

func (r *ContentRepository) CreateTrainer(content *Trainer) (uint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.DatabaseInfo("Создание тренера: %s (TgId: %s, ChatId: %d)", content.Name, content.TgId, content.ChatId)

	result := r.db.WithContext(ctx).Create(content)
	if result.Error != nil {
		logger.DatabaseError("Создание тренера %s: %v", content.Name, result.Error)
		return 0, result.Error
	}

	logger.DatabaseInfo("Тренер создан: %s (%d)", content.Name, content.ID)
	return content.ID, nil
}

func (r *ContentRepository) GetTrainerByID(ID uint) (*Trainer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var trainer Trainer
	result := r.db.WithContext(ctx).First(&trainer, ID)
	if result.Error != nil {
		logger.DatabaseError("Не удалось получить тренера %d: %v", ID, result.Error)
		return nil, result.Error
	}

	return &trainer, nil
}

func (r *ContentRepository) GetTrainers() ([]Trainer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var trainers []Trainer
	result := r.db.WithContext(ctx).Find(&trainers)
	if result.Error != nil {
		logger.DatabaseError("Не удалось получить тренеров: %v", result.Error)
		return nil, result.Error
	}

	return trainers, nil
}

func (r *ContentRepository) UpdateTrainer(id uint, trainer *Trainer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := r.db.WithContext(ctx).Model(&Trainer{}).Where("id = ?", id).Updates(trainer)
	if result.Error != nil {
		logger.DatabaseError("Не удалось обновить тренера %d: %v", id, result.Error)
		return result.Error
	}

	logger.DatabaseInfo("Тренер обновлен: ID=%d", id)
	return nil
}

func (r *ContentRepository) DeleteTrainer(id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := r.db.WithContext(ctx).Delete(&Trainer{}, id)
	if result.Error != nil {
		logger.DatabaseError("Не удалось удалить тренера %d: %v", id, result.Error)
		return result.Error
	}

	logger.DatabaseInfo("Тренер удален: ID=%d", id)
	return nil
}

func (r *ContentRepository) CreateTrack(content *Track) (uint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.DatabaseInfo("Создание трека: %s", content.Name)

	result := r.db.WithContext(ctx).Create(content)
	if result.Error != nil {
		logger.DatabaseError("Не удалось создать трек %s: %v", content.Name, result.Error)
		return 0, result.Error
	}

	logger.DatabaseInfo("Трек создан: %s (ID: %d)", content.Name, content.ID)
	return content.ID, nil
}

func (r *ContentRepository) GetTrackByID(id uint) (*Track, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var track Track
	result := r.db.WithContext(ctx).First(&track, id)
	if result.Error != nil {
		logger.DatabaseError("Не удалось получить трек %d: %v", id, result.Error)
		return nil, result.Error
	}

	return &track, nil
}

func (r *ContentRepository) GetTracks() ([]Track, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var tracks []Track
	result := r.db.WithContext(ctx).Find(&tracks)
	if result.Error != nil {
		logger.DatabaseError("Не удалось получить треки: %v", result.Error)
		return nil, result.Error
	}

	return tracks, nil
}

func (r *ContentRepository) UpdateTrack(id uint, track *Track) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := r.db.WithContext(ctx).Model(&Track{}).Where("id = ?", id).Updates(track)
	if result.Error != nil {
		logger.DatabaseError("Не удалось обновить трек %d: %v", id, result.Error)
		return result.Error
	}

	logger.DatabaseInfo("Трек обновлен: ID=%d", id)
	return nil
}

func (r *ContentRepository) DeleteTrack(id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := r.db.WithContext(ctx).Delete(&Track{}, id)
	if result.Error != nil {
		logger.DatabaseError("Не удалось удалить трек %d: %v", id, result.Error)
		return result.Error
	}

	logger.DatabaseInfo("Трек удален: ID=%d", id)
	return nil
}

func (r *ContentRepository) CreateUser(user *User) (uint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.DatabaseInfo("Создание пользователя: %s (ChatId: %d)", user.Name, user.ChatId)

	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		logger.DatabaseError("Не удалось создать пользователя %s: %v", user.Name, result.Error)
		return 0, result.Error
	}

	logger.DatabaseInfo("Пользователь создан: %s (ID: %d)", user.Name, user.ID)
	return user.ID, nil
}

func (r *ContentRepository) GetUserByID(id uint) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user User
	result := r.db.WithContext(ctx).First(&user, id)
	if result.Error != nil {
		logger.DatabaseError("Не удалось получить пользователя %d: %v", id, result.Error)
		return nil, result.Error
	}

	return &user, nil
}

func (r *ContentRepository) GetUserByChatId(chatId int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user User
	result := r.db.WithContext(ctx).Where("chat_id = ? AND is_active = ?", chatId, true).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logger.DatabaseInfo("User with chat_id %d not found", chatId)
		} else {
			logger.DatabaseError("Не удалось получить пользователя по chat_id %d: %v", chatId, result.Error)
		}
		return nil, result.Error
	}

	return &user, nil
}

func (r *ContentRepository) GetUsers() ([]User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var users []User
	result := r.db.WithContext(ctx).Find(&users)
	if result.Error != nil {
		logger.DatabaseError("Не удалось получить пользователей: %v", result.Error)
		return nil, result.Error
	}

	return users, nil
}

func (r *ContentRepository) UpdateUser(id uint, user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := r.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Updates(user)
	if result.Error != nil {
		logger.DatabaseError("Не удалось обновить пользователя %d: %v", id, result.Error)
		return result.Error
	}

	logger.DatabaseInfo("User updated successfully: ID=%d", id)
	return nil
}

func (r *ContentRepository) DeleteUser(id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := r.db.WithContext(ctx).Delete(&User{}, id)
	if result.Error != nil {
		logger.DatabaseError("Не удалось удалить пользователя %d: %v", id, result.Error)
		return result.Error
	}

	logger.DatabaseInfo("User deleted successfully: ID=%d", id)
	return nil
}

func (r *ContentRepository) UpdateUserELORating(id uint, newRating int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := r.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Update("elo_rating", newRating)
	if result.Error != nil {
		logger.DatabaseError("Не удалось обновить ELO рейтинг пользователя %d: %v", id, result.Error)
		return result.Error
	}

	logger.DatabaseInfo("User ELO rating updated successfully: ID=%d, New Rating=%d", id, newRating)
	return nil
}

func (r *ContentRepository) CreateTrainingRegistration(registration *TrainingRegistration) (uint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.DatabaseInfo("Creating training registration: TrainingID=%d, UserID=%d", registration.TrainingID, registration.UserID)

	result := r.db.WithContext(ctx).Create(registration)
	if result.Error != nil {
		logger.DatabaseError("Не удалось создать регистрацию на тренировку: %v", result.Error)
		return 0, result.Error
	}

	logger.DatabaseInfo("Training registration created successfully: ID=%d", registration.ID)
	return registration.ID, nil
}

func (r *ContentRepository) GetTrainingRegistrationByID(id uint) (*TrainingRegistration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var registration TrainingRegistration
	result := r.db.WithContext(ctx).First(&registration, id)
	if result.Error != nil {
		logger.DatabaseError("Не удалось получить регистрацию на тренировку %d: %v", id, result.Error)
		return nil, result.Error
	}

	return &registration, nil
}

func (r *ContentRepository) GetTrainingRegistrationsByTrainingID(trainingId uint) ([]TrainingRegistration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var registrations []TrainingRegistration
	result := r.db.WithContext(ctx).Where("training_id = ?", trainingId).Find(&registrations)
	if result.Error != nil {
		logger.DatabaseError("Failed to get training registrations for training %d: %v", trainingId, result.Error)
		return nil, result.Error
	}

	return registrations, nil
}

func (r *ContentRepository) GetTrainingRegistrationsByUserID(userId uint) ([]TrainingRegistration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var registrations []TrainingRegistration
	result := r.db.WithContext(ctx).Where("user_id = ?", userId).Find(&registrations)
	if result.Error != nil {
		logger.DatabaseError("Failed to get training registrations for user %d: %v", userId, result.Error)
		return nil, result.Error
	}

	return registrations, nil
}

func (r *ContentRepository) UpdateTrainingRegistration(id uint, registration *TrainingRegistration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := r.db.WithContext(ctx).Model(&TrainingRegistration{}).Where("id = ?", id).Updates(registration)
	if result.Error != nil {
		logger.DatabaseError("Failed to update training registration %d: %v", id, result.Error)
		return result.Error
	}

	logger.DatabaseInfo("Training registration updated successfully: ID=%d", id)
	return nil
}

func (r *ContentRepository) DeleteTrainingRegistration(id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := r.db.WithContext(ctx).Delete(&TrainingRegistration{}, id)
	if result.Error != nil {
		logger.DatabaseError("Failed to delete training registration %d: %v", id, result.Error)
		return result.Error
	}

	logger.DatabaseInfo("Training registration deleted successfully: ID=%d", id)
	return nil
}

func (r *ContentRepository) GetTrainingRegistrationByUserAndTraining(userId uint, trainingId uint) (*TrainingRegistration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var registration TrainingRegistration
	result := r.db.WithContext(ctx).Where("user_id = ? AND training_id = ?", userId, trainingId).First(&registration)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logger.DatabaseInfo("Training registration not found for user %d and training %d", userId, trainingId)
		} else {
			logger.DatabaseError("Failed to get training registration for user %d and training %d: %v", userId, trainingId, result.Error)
		}
		return nil, result.Error
	}

	return &registration, nil
}

func (r *ContentRepository) CreateAdmin(admin *Admin) (uint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.DatabaseInfo("Creating admin: %s (TgId: %s, ChatId: %d)", admin.Name, admin.TgId, admin.ChatId)

	result := r.db.WithContext(ctx).Create(admin)
	if result.Error != nil {
		logger.DatabaseError("Failed to create admin %s: %v", admin.Name, result.Error)
		return 0, result.Error
	}

	logger.DatabaseInfo("Admin created successfully: %s (ID: %d)", admin.Name, admin.ID)
	return admin.ID, nil
}

func (r *ContentRepository) GetAdminByID(id uint) (*Admin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var admin Admin
	result := r.db.WithContext(ctx).First(&admin, id)
	if result.Error != nil {
		logger.DatabaseError("Failed to get admin %d: %v", id, result.Error)
		return nil, result.Error
	}

	return &admin, nil
}

func (r *ContentRepository) GetAdminByChatId(chatId int) (*Admin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var admin Admin
	result := r.db.WithContext(ctx).Where("chat_id = ? AND is_active = ?", chatId, true).First(&admin)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logger.DatabaseInfo("Admin with chat_id %d not found", chatId)
		} else {
			logger.DatabaseError("Failed to get admin by chat_id %d: %v", chatId, result.Error)
		}
		return nil, result.Error
	}

	return &admin, nil
}

func (r *ContentRepository) GetAdmins() ([]Admin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var admins []Admin
	result := r.db.WithContext(ctx).Find(&admins)
	if result.Error != nil {
		logger.DatabaseError("Failed to get admins: %v", result.Error)
		return nil, result.Error
	}

	return admins, nil
}

func (r *ContentRepository) UpdateAdmin(id uint, admin *Admin) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := r.db.WithContext(ctx).Model(&Admin{}).Where("id = ?", id).Updates(admin)
	if result.Error != nil {
		logger.DatabaseError("Failed to update admin %d: %v", id, result.Error)
		return result.Error
	}

	logger.DatabaseInfo("Admin updated successfully: ID=%d", id)
	return nil
}

func (r *ContentRepository) DeleteAdmin(id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := r.db.WithContext(ctx).Delete(&Admin{}, id)
	if result.Error != nil {
		logger.DatabaseError("Failed to delete admin %d: %v", id, result.Error)
		return result.Error
	}

	logger.DatabaseInfo("Admin deleted successfully: ID=%d", id)
	return nil
}

func IsAdmin(chatId int, repo ContentRepositoryInterface) bool {
	admin, err := repo.GetAdminByChatId(chatId)
	if err != nil {
		logger.DatabaseError("Admin check failed for user %d: %v", chatId, err)
		return false
	}
	if admin == nil {
		logger.DatabaseInfo("User %d not found in admins database", chatId)
		return false
	}

	logger.DatabaseInfo("User %d (%s) has admin rights", chatId, admin.Name)
	return true
}

func (r *ContentRepository) GetActiveTrainingsByTrackAndTrainer(trackId, trainerId uint) ([]Training, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var trainings []Training
	result := r.db.WithContext(ctx).Where("track_id = ? AND trainer_id = ? AND is_active = ?", trackId, trainerId, true).Find(&trainings)
	if result.Error != nil {
		logger.DatabaseError("Failed to get active trainings by track and trainer: %v", result.Error)
		return nil, result.Error
	}

	return trainings, nil
}

func (r *ContentRepository) GetTrainersByTrack(trackId uint) ([]Trainer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var trainers []Trainer
	result := r.db.WithContext(ctx).
		Table("trainers").
		Select("DISTINCT trainers.*").
		Joins("INNER JOIN trainings ON trainers.id = trainings.trainer_id").
		Where("trainings.track_id = ? AND trainings.is_active = ? AND trainings.start_time > ?", trackId, true, time.Now()).
		Find(&trainers)

	if result.Error != nil {
		logger.DatabaseError("Failed to get trainers by track: %v", result.Error)
		return nil, result.Error
	}

	return trainers, nil
}

func (r *ContentRepository) GetTracksWithActiveTrainings() ([]Track, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var tracks []Track
	result := r.db.WithContext(ctx).
		Table("tracks").
		Select("DISTINCT tracks.*").
		Joins("INNER JOIN trainings ON tracks.id = trainings.track_id").
		Where("trainings.is_active = ? AND trainings.start_time > ?", true, time.Now()).
		Find(&tracks)

	if result.Error != nil {
		logger.DatabaseError("Failed to get tracks with active trainings: %v", result.Error)
		return nil, result.Error
	}

	return tracks, nil
}

func (r *ContentRepository) GetUserTrainings(userId uint) ([]Training, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var trainings []Training
	result := r.db.WithContext(ctx).
		Table("trainings").
		Select("DISTINCT trainings.*").
		Joins("INNER JOIN training_registrations ON trainings.id = training_registrations.training_id").
		Where("training_registrations.user_id = ? AND trainings.is_active = ? AND trainings.start_time > ?", userId, true, time.Now()).
		Find(&trainings)

	if result.Error != nil {
		logger.DatabaseError("Failed to get user trainings: %v", result.Error)
		return nil, result.Error
	}

	return trainings, nil
}

// TrainingRequest methods
func (r *ContentRepository) CreateTrainingRequest(request *TrainingRequest) (uint, error) {
	result := r.db.Create(request)
	if result.Error != nil {
		logger.DatabaseError("Failed to create training request: %v", result.Error)
		return 0, result.Error
	}

	logger.DatabaseInfo("Training request created with ID: %d", request.ID)
	return request.ID, nil
}

func (r *ContentRepository) GetTrainingRequestByID(id uint) (*TrainingRequest, error) {
	var request TrainingRequest
	result := r.db.First(&request, id)
	if result.Error != nil {
		logger.DatabaseError("Failed to get training request by ID %d: %v", id, result.Error)
		return nil, result.Error
	}

	return &request, nil
}

func (r *ContentRepository) GetTrainingRequests() ([]TrainingRequest, error) {
	var requests []TrainingRequest
	result := r.db.Order("created_at DESC").Find(&requests)
	if result.Error != nil {
		logger.DatabaseError("Failed to get training requests: %v", result.Error)
		return nil, result.Error
	}

	return requests, nil
}

func (r *ContentRepository) GetUnreviewedTrainingRequests() ([]TrainingRequest, error) {
	var requests []TrainingRequest
	result := r.db.Where("is_reviewed = ?", false).Order("created_at ASC").Find(&requests)
	if result.Error != nil {
		logger.DatabaseError("Failed to get unreviewed training requests: %v", result.Error)
		return nil, result.Error
	}

	return requests, nil
}

func (r *ContentRepository) UpdateTrainingRequest(id uint, request *TrainingRequest) error {
	result := r.db.Model(&TrainingRequest{}).Where("id = ?", id).Updates(request)
	if result.Error != nil {
		logger.DatabaseError("Failed to update training request %d: %v", id, result.Error)
		return result.Error
	}

	logger.DatabaseInfo("Training request %d updated", id)
	return nil
}

func (r *ContentRepository) DeleteTrainingRequest(id uint) error {
	result := r.db.Delete(&TrainingRequest{}, id)
	if result.Error != nil {
		logger.DatabaseError("Failed to delete training request %d: %v", id, result.Error)
		return result.Error
	}

	logger.DatabaseInfo("Training request %d deleted", id)
	return nil
}
