package database

import (
	"context"
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

type ContentRepository struct{}

func (r *ContentRepository) CreateTraining(content *Training) (uint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("Creating training: TrainerID=%d, TrackID=%d", content.TrainerID, content.TrackID)

	result := db.WithContext(ctx).Create(content)
	if result.Error != nil {
		log.Printf("ОШИБКА: Не удалось создать тренировку: %v", result.Error)
		return 0, result.Error
	}

	log.Printf("Training created successfully: ID=%d", content.ID)
	return content.ID, nil
}

func (r *ContentRepository) GetTrainingById(id uint) (*Training, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var training Training
	result := db.WithContext(ctx).First(&training, id)
	if result.Error != nil {
		log.Printf("ОШИБКА: Не удалось получить тренировку %d: %v", id, result.Error)
		return nil, result.Error
	}

	return &training, nil
}

func (r *ContentRepository) GetTrainings() ([]Training, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var trainings []Training
	result := db.WithContext(ctx).Find(&trainings)
	if result.Error != nil {
		log.Printf("ОШИБКА: Не удалось получить тренировки: %v", result.Error)
		return nil, result.Error
	}

	return trainings, nil
}

func (r *ContentRepository) GetActiveTrainings() ([]Training, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var trainings []Training
	result := db.WithContext(ctx).Where("is_active = ?", true).Find(&trainings)
	if result.Error != nil {
		log.Printf("ОШИБКА: Не удалось получить активные тренировки: %v", result.Error)
		return nil, result.Error
	}

	return trainings, nil
}

func (r *ContentRepository) UpdateTraining(id uint, training *Training) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := db.WithContext(ctx).Model(&Training{}).Where("id = ?", id).Updates(training)
	if result.Error != nil {
		log.Printf("ОШИБКА: Не удалось обновить тренировку %d: %v", id, result.Error)
		return result.Error
	}

	log.Printf("Training updated successfully: ID=%d", id)
	return nil
}

func (r *ContentRepository) DeleteTraining(id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := db.WithContext(ctx).Delete(&Training{}, id)
	if result.Error != nil {
		log.Printf("ОШИБКА: Не удалось удалить тренировку %d: %v", id, result.Error)
		return result.Error
	}

	log.Printf("Training deleted successfully: ID=%d", id)
	return nil
}

func (r *ContentRepository) CreateTrainer(content *Trainer) (uint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("Creating trainer: %s (TgId: %s, ChatId: %d)", content.Name, content.TgId, content.ChatId)

	result := db.WithContext(ctx).Create(content)
	if result.Error != nil {
		log.Printf("ОШИБКА: Не удалось создать тренера %s: %v", content.Name, result.Error)
		return 0, result.Error
	}

	log.Printf("Trainer created successfully: %s (ID: %d)", content.Name, content.ID)
	return content.ID, nil
}

func (r *ContentRepository) GetTrainerByID(ID uint) (*Trainer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var trainer Trainer
	result := db.WithContext(ctx).First(&trainer, ID)
	if result.Error != nil {
		log.Printf("ОШИБКА: Не удалось получить тренера %d: %v", ID, result.Error)
		return nil, result.Error
	}

	return &trainer, nil
}

func (r *ContentRepository) GetTrainers() ([]Trainer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var trainers []Trainer
	result := db.WithContext(ctx).Find(&trainers)
	if result.Error != nil {
		log.Printf("ОШИБКА: Не удалось получить тренеров: %v", result.Error)
		return nil, result.Error
	}

	return trainers, nil
}

func (r *ContentRepository) UpdateTrainer(id uint, trainer *Trainer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := db.WithContext(ctx).Model(&Trainer{}).Where("id = ?", id).Updates(trainer)
	if result.Error != nil {
		log.Printf("ОШИБКА: Не удалось обновить тренера %d: %v", id, result.Error)
		return result.Error
	}

	log.Printf("Trainer updated successfully: ID=%d", id)
	return nil
}

func (r *ContentRepository) DeleteTrainer(id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := db.WithContext(ctx).Delete(&Trainer{}, id)
	if result.Error != nil {
		log.Printf("ОШИБКА: Не удалось удалить тренера %d: %v", id, result.Error)
		return result.Error
	}

	log.Printf("Trainer deleted successfully: ID=%d", id)
	return nil
}

func (r *ContentRepository) CreateTrack(content *Track) (uint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("Creating track: %s", content.Name)

	result := db.WithContext(ctx).Create(content)
	if result.Error != nil {
		log.Printf("ОШИБКА: Не удалось создать трек %s: %v", content.Name, result.Error)
		return 0, result.Error
	}

	log.Printf("Track created successfully: %s (ID: %d)", content.Name, content.ID)
	return content.ID, nil
}

func (r *ContentRepository) GetTrackByID(id uint) (*Track, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var track Track
	result := db.WithContext(ctx).First(&track, id)
	if result.Error != nil {
		log.Printf("ОШИБКА: Не удалось получить трек %d: %v", id, result.Error)
		return nil, result.Error
	}

	return &track, nil
}

func (r *ContentRepository) GetTracks() ([]Track, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var tracks []Track
	result := db.WithContext(ctx).Find(&tracks)
	if result.Error != nil {
		log.Printf("ОШИБКА: Не удалось получить треки: %v", result.Error)
		return nil, result.Error
	}

	return tracks, nil
}

func (r *ContentRepository) UpdateTrack(id uint, track *Track) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := db.WithContext(ctx).Model(&Track{}).Where("id = ?", id).Updates(track)
	if result.Error != nil {
		log.Printf("ОШИБКА: Не удалось обновить трек %d: %v", id, result.Error)
		return result.Error
	}

	log.Printf("Track updated successfully: ID=%d", id)
	return nil
}

func (r *ContentRepository) DeleteTrack(id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := db.WithContext(ctx).Delete(&Track{}, id)
	if result.Error != nil {
		log.Printf("ОШИБКА: Не удалось удалить трек %d: %v", id, result.Error)
		return result.Error
	}

	log.Printf("Track deleted successfully: ID=%d", id)
	return nil
}

func (r *ContentRepository) CreateUser(user *User) (uint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("Creating user: %s (ChatId: %d)", user.Name, user.ChatId)

	result := db.WithContext(ctx).Create(user)
	if result.Error != nil {
		log.Printf("ОШИБКА: Не удалось создать пользователя %s: %v", user.Name, result.Error)
		return 0, result.Error
	}

	log.Printf("User created successfully: %s (ID: %d)", user.Name, user.ID)
	return user.ID, nil
}

func (r *ContentRepository) GetUserByID(id uint) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user User
	result := db.WithContext(ctx).First(&user, id)
	if result.Error != nil {
		log.Printf("ОШИБКА: Не удалось получить пользователя %d: %v", id, result.Error)
		return nil, result.Error
	}

	return &user, nil
}

func (r *ContentRepository) GetUserByChatId(chatId int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user User
	result := db.WithContext(ctx).Where("chat_id = ? AND is_active = ?", chatId, true).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("User with chat_id %d not found", chatId)
		} else {
			log.Printf("ОШИБКА: Не удалось получить пользователя по chat_id %d: %v", chatId, result.Error)
		}
		return nil, result.Error
	}

	return &user, nil
}

func (r *ContentRepository) GetUsers() ([]User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var users []User
	result := db.WithContext(ctx).Find(&users)
	if result.Error != nil {
		log.Printf("ОШИБКА: Не удалось получить пользователей: %v", result.Error)
		return nil, result.Error
	}

	return users, nil
}

func (r *ContentRepository) UpdateUser(id uint, user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Updates(user)
	if result.Error != nil {
		log.Printf("ОШИБКА: Не удалось обновить пользователя %d: %v", id, result.Error)
		return result.Error
	}

	log.Printf("User updated successfully: ID=%d", id)
	return nil
}

func (r *ContentRepository) DeleteUser(id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := db.WithContext(ctx).Delete(&User{}, id)
	if result.Error != nil {
		log.Printf("ОШИБКА: Не удалось удалить пользователя %d: %v", id, result.Error)
		return result.Error
	}

	log.Printf("User deleted successfully: ID=%d", id)
	return nil
}

func (r *ContentRepository) UpdateUserELORating(id uint, newRating int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Update("elo_rating", newRating)
	if result.Error != nil {
		log.Printf("ОШИБКА: Не удалось обновить ELO рейтинг пользователя %d: %v", id, result.Error)
		return result.Error
	}

	log.Printf("User ELO rating updated successfully: ID=%d, New Rating=%d", id, newRating)
	return nil
}

func (r *ContentRepository) CreateTrainingRegistration(registration *TrainingRegistration) (uint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("Creating training registration: TrainingID=%d, UserID=%d", registration.TrainingID, registration.UserID)

	result := db.WithContext(ctx).Create(registration)
	if result.Error != nil {
		log.Printf("ОШИБКА: Не удалось создать регистрацию на тренировку: %v", result.Error)
		return 0, result.Error
	}

	log.Printf("Training registration created successfully: ID=%d", registration.ID)
	return registration.ID, nil
}

func (r *ContentRepository) GetTrainingRegistrationByID(id uint) (*TrainingRegistration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var registration TrainingRegistration
	result := db.WithContext(ctx).First(&registration, id)
	if result.Error != nil {
		log.Printf("ОШИБКА: Не удалось получить регистрацию на тренировку %d: %v", id, result.Error)
		return nil, result.Error
	}

	return &registration, nil
}

func (r *ContentRepository) GetTrainingRegistrationsByTrainingID(trainingId uint) ([]TrainingRegistration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var registrations []TrainingRegistration
	result := db.WithContext(ctx).Where("training_id = ?", trainingId).Find(&registrations)
	if result.Error != nil {
		log.Printf("ERROR: Failed to get training registrations for training %d: %v", trainingId, result.Error)
		return nil, result.Error
	}

	return registrations, nil
}

func (r *ContentRepository) GetTrainingRegistrationsByUserID(userId uint) ([]TrainingRegistration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var registrations []TrainingRegistration
	result := db.WithContext(ctx).Where("user_id = ?", userId).Find(&registrations)
	if result.Error != nil {
		log.Printf("ERROR: Failed to get training registrations for user %d: %v", userId, result.Error)
		return nil, result.Error
	}

	return registrations, nil
}

func (r *ContentRepository) UpdateTrainingRegistration(id uint, registration *TrainingRegistration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := db.WithContext(ctx).Model(&TrainingRegistration{}).Where("id = ?", id).Updates(registration)
	if result.Error != nil {
		log.Printf("ERROR: Failed to update training registration %d: %v", id, result.Error)
		return result.Error
	}

	log.Printf("Training registration updated successfully: ID=%d", id)
	return nil
}

func (r *ContentRepository) DeleteTrainingRegistration(id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := db.WithContext(ctx).Delete(&TrainingRegistration{}, id)
	if result.Error != nil {
		log.Printf("ERROR: Failed to delete training registration %d: %v", id, result.Error)
		return result.Error
	}

	log.Printf("Training registration deleted successfully: ID=%d", id)
	return nil
}

func (r *ContentRepository) GetTrainingRegistrationByUserAndTraining(userId uint, trainingId uint) (*TrainingRegistration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var registration TrainingRegistration
	result := db.WithContext(ctx).Where("user_id = ? AND training_id = ?", userId, trainingId).First(&registration)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("Training registration not found for user %d and training %d", userId, trainingId)
		} else {
			log.Printf("ERROR: Failed to get training registration for user %d and training %d: %v", userId, trainingId, result.Error)
		}
		return nil, result.Error
	}

	return &registration, nil
}

func (r *ContentRepository) CreateAdmin(admin *Admin) (uint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("Creating admin: %s (TgId: %s, ChatId: %d)", admin.Name, admin.TgId, admin.ChatId)

	result := db.WithContext(ctx).Create(admin)
	if result.Error != nil {
		log.Printf("ERROR: Failed to create admin %s: %v", admin.Name, result.Error)
		return 0, result.Error
	}

	log.Printf("Admin created successfully: %s (ID: %d)", admin.Name, admin.ID)
	return admin.ID, nil
}

func (r *ContentRepository) GetAdminByID(id uint) (*Admin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var admin Admin
	result := db.WithContext(ctx).First(&admin, id)
	if result.Error != nil {
		log.Printf("ERROR: Failed to get admin %d: %v", id, result.Error)
		return nil, result.Error
	}

	return &admin, nil
}

func (r *ContentRepository) GetAdminByChatId(chatId int) (*Admin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var admin Admin
	result := db.WithContext(ctx).Where("chat_id = ? AND is_active = ?", chatId, true).First(&admin)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("Admin with chat_id %d not found", chatId)
		} else {
			log.Printf("ERROR: Failed to get admin by chat_id %d: %v", chatId, result.Error)
		}
		return nil, result.Error
	}

	return &admin, nil
}

func (r *ContentRepository) GetAdmins() ([]Admin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var admins []Admin
	result := db.WithContext(ctx).Find(&admins)
	if result.Error != nil {
		log.Printf("ERROR: Failed to get admins: %v", result.Error)
		return nil, result.Error
	}

	return admins, nil
}

func (r *ContentRepository) UpdateAdmin(id uint, admin *Admin) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := db.WithContext(ctx).Model(&Admin{}).Where("id = ?", id).Updates(admin)
	if result.Error != nil {
		log.Printf("ERROR: Failed to update admin %d: %v", id, result.Error)
		return result.Error
	}

	log.Printf("Admin updated successfully: ID=%d", id)
	return nil
}

func (r *ContentRepository) DeleteAdmin(id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := db.WithContext(ctx).Delete(&Admin{}, id)
	if result.Error != nil {
		log.Printf("ERROR: Failed to delete admin %d: %v", id, result.Error)
		return result.Error
	}

	log.Printf("Admin deleted successfully: ID=%d", id)
	return nil
}

func IsAdmin(chatId int, repo ContentRepositoryInterface) bool {
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

func (r *ContentRepository) GetActiveTrainingsByTrackAndTrainer(trackId, trainerId uint) ([]Training, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var trainings []Training
	result := db.WithContext(ctx).Where("track_id = ? AND trainer_id = ? AND is_active = ?", trackId, trainerId, true).Find(&trainings)
	if result.Error != nil {
		log.Printf("ERROR: Failed to get active trainings by track and trainer: %v", result.Error)
		return nil, result.Error
	}

	return trainings, nil
}

func (r *ContentRepository) GetTrainersByTrack(trackId uint) ([]Trainer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var trainers []Trainer
	result := db.WithContext(ctx).
		Table("trainers").
		Select("DISTINCT trainers.*").
		Joins("INNER JOIN trainings ON trainers.id = trainings.trainer_id").
		Where("trainings.track_id = ? AND trainings.is_active = ? AND trainings.start_time > ?", trackId, true, time.Now()).
		Find(&trainers)

	if result.Error != nil {
		log.Printf("ERROR: Failed to get trainers by track: %v", result.Error)
		return nil, result.Error
	}

	return trainers, nil
}

func (r *ContentRepository) GetTracksWithActiveTrainings() ([]Track, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var tracks []Track
	result := db.WithContext(ctx).
		Table("tracks").
		Select("DISTINCT tracks.*").
		Joins("INNER JOIN trainings ON tracks.id = trainings.track_id").
		Where("trainings.is_active = ? AND trainings.start_time > ?", true, time.Now()).
		Find(&tracks)

	if result.Error != nil {
		log.Printf("ERROR: Failed to get tracks with active trainings: %v", result.Error)
		return nil, result.Error
	}

	return tracks, nil
}

func (r *ContentRepository) GetUserTrainings(userId uint) ([]Training, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var trainings []Training
	result := db.WithContext(ctx).
		Table("trainings").
		Select("DISTINCT trainings.*").
		Joins("INNER JOIN training_registrations ON trainings.id = training_registrations.training_id").
		Where("training_registrations.user_id = ? AND trainings.is_active = ? AND trainings.start_time > ?", userId, true, time.Now()).
		Find(&trainings)

	if result.Error != nil {
		log.Printf("ERROR: Failed to get user trainings: %v", result.Error)
		return nil, result.Error
	}

	return trainings, nil
}
