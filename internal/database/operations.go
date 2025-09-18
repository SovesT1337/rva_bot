package database

import (
	"context"
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

// ContentRepository - репозиторий для работы с контентом
type ContentRepository struct{}

// ============================================================================
// МЕТОДЫ ДЛЯ РАБОТЫ С ТРЕНИРОВКАМИ
// ============================================================================

// CreateTraining создает новую тренировку
func (r *ContentRepository) CreateTraining(content *Training) (uint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("Creating training: TrainerID=%d, TrackID=%d", content.TrainerID, content.TrackID)

	result := db.WithContext(ctx).Create(content)
	if result.Error != nil {
		log.Printf("ERROR: Failed to create training: %v", result.Error)
		return 0, result.Error
	}

	log.Printf("Training created successfully: ID=%d", content.ID)
	return content.ID, nil
}

// GetTrainingById получает тренировку по ID
func (r *ContentRepository) GetTrainingById(id uint) (*Training, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var training Training
	result := db.WithContext(ctx).First(&training, id)
	if result.Error != nil {
		log.Printf("ERROR: Failed to get training %d: %v", id, result.Error)
		return nil, result.Error
	}

	return &training, nil
}

// GetTrainings получает все тренировки
func (r *ContentRepository) GetTrainings() ([]Training, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var trainings []Training
	result := db.WithContext(ctx).Find(&trainings)
	if result.Error != nil {
		log.Printf("ERROR: Failed to get trainings: %v", result.Error)
		return nil, result.Error
	}

	return trainings, nil
}

// GetActiveTrainings получает только активные тренировки
func (r *ContentRepository) GetActiveTrainings() ([]Training, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var trainings []Training
	result := db.WithContext(ctx).Where("is_active = ?", true).Find(&trainings)
	if result.Error != nil {
		log.Printf("ERROR: Failed to get active trainings: %v", result.Error)
		return nil, result.Error
	}

	return trainings, nil
}

// UpdateTraining обновляет тренировку
func (r *ContentRepository) UpdateTraining(id uint, training *Training) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := db.WithContext(ctx).Model(&Training{}).Where("id = ?", id).Updates(training)
	if result.Error != nil {
		log.Printf("ERROR: Failed to update training %d: %v", id, result.Error)
		return result.Error
	}

	log.Printf("Training updated successfully: ID=%d", id)
	return nil
}

// DeleteTraining удаляет тренировку
func (r *ContentRepository) DeleteTraining(id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := db.WithContext(ctx).Delete(&Training{}, id)
	if result.Error != nil {
		log.Printf("ERROR: Failed to delete training %d: %v", id, result.Error)
		return result.Error
	}

	log.Printf("Training deleted successfully: ID=%d", id)
	return nil
}

// ============================================================================
// МЕТОДЫ ДЛЯ РАБОТЫ С ТРЕНЕРАМИ
// ============================================================================

// CreateTrainer создает нового тренера
func (r *ContentRepository) CreateTrainer(content *Trainer) (uint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("Creating trainer: %s (TgId: %s, ChatId: %d)", content.Name, content.TgId, content.ChatId)

	result := db.WithContext(ctx).Create(content)
	if result.Error != nil {
		log.Printf("ERROR: Failed to create trainer %s: %v", content.Name, result.Error)
		return 0, result.Error
	}

	log.Printf("Trainer created successfully: %s (ID: %d)", content.Name, content.ID)
	return content.ID, nil
}

// GetTrainerByID получает тренера по ID
func (r *ContentRepository) GetTrainerByID(ID uint) (*Trainer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var trainer Trainer
	result := db.WithContext(ctx).First(&trainer, ID)
	if result.Error != nil {
		log.Printf("ERROR: Failed to get trainer %d: %v", ID, result.Error)
		return nil, result.Error
	}

	return &trainer, nil
}

// GetTrainers получает всех тренеров
func (r *ContentRepository) GetTrainers() ([]Trainer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var trainers []Trainer
	result := db.WithContext(ctx).Find(&trainers)
	if result.Error != nil {
		log.Printf("ERROR: Failed to get trainers: %v", result.Error)
		return nil, result.Error
	}

	return trainers, nil
}

// UpdateTrainer обновляет данные тренера
func (r *ContentRepository) UpdateTrainer(id uint, trainer *Trainer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := db.WithContext(ctx).Model(&Trainer{}).Where("id = ?", id).Updates(trainer)
	if result.Error != nil {
		log.Printf("ERROR: Failed to update trainer %d: %v", id, result.Error)
		return result.Error
	}

	log.Printf("Trainer updated successfully: ID=%d", id)
	return nil
}

// UpdateTrainerTgId обновляет Telegram ID тренера
func (r *ContentRepository) UpdateTrainerTgId(id uint, tgid string) (uint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := db.WithContext(ctx).Model(&Trainer{}).Where("id = ?", id).Update("tg_id", tgid)
	if result.Error != nil {
		log.Printf("ERROR: Failed to update trainer TgId %d: %v", id, result.Error)
		return 0, result.Error
	}

	log.Printf("Trainer TgId updated successfully: ID=%d", id)
	return id, nil
}

// UpdateTrainerName обновляет имя тренера
func (r *ContentRepository) UpdateTrainerName(id uint, name string) (uint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := db.WithContext(ctx).Model(&Trainer{}).Where("id = ?", id).Update("name", name)
	if result.Error != nil {
		log.Printf("ERROR: Failed to update trainer name %d: %v", id, result.Error)
		return 0, result.Error
	}

	log.Printf("Trainer name updated successfully: ID=%d", id)
	return id, nil
}

// UpdateTrainerInfo обновляет информацию о тренере
func (r *ContentRepository) UpdateTrainerInfo(id uint, info string) (uint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := db.WithContext(ctx).Model(&Trainer{}).Where("id = ?", id).Update("info", info)
	if result.Error != nil {
		log.Printf("ERROR: Failed to update trainer info %d: %v", id, result.Error)
		return 0, result.Error
	}

	log.Printf("Trainer info updated successfully: ID=%d", id)
	return id, nil
}

// DeleteTrainer удаляет тренера
func (r *ContentRepository) DeleteTrainer(id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := db.WithContext(ctx).Delete(&Trainer{}, id)
	if result.Error != nil {
		log.Printf("ERROR: Failed to delete trainer %d: %v", id, result.Error)
		return result.Error
	}

	log.Printf("Trainer deleted successfully: ID=%d", id)
	return nil
}

// ============================================================================
// МЕТОДЫ ДЛЯ РАБОТЫ С ТРАССАМИ
// ============================================================================

// CreateTrack создает новую трассу
func (r *ContentRepository) CreateTrack(content *Track) (uint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("Creating track: %s", content.Name)

	result := db.WithContext(ctx).Create(content)
	if result.Error != nil {
		log.Printf("ERROR: Failed to create track %s: %v", content.Name, result.Error)
		return 0, result.Error
	}

	log.Printf("Track created successfully: %s (ID: %d)", content.Name, content.ID)
	return content.ID, nil
}

// GetTrackByID получает трассу по ID
func (r *ContentRepository) GetTrackByID(id uint) (*Track, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var track Track
	result := db.WithContext(ctx).First(&track, id)
	if result.Error != nil {
		log.Printf("ERROR: Failed to get track %d: %v", id, result.Error)
		return nil, result.Error
	}

	return &track, nil
}

// GetTracks получает все трассы
func (r *ContentRepository) GetTracks() ([]Track, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var tracks []Track
	result := db.WithContext(ctx).Find(&tracks)
	if result.Error != nil {
		log.Printf("ERROR: Failed to get tracks: %v", result.Error)
		return nil, result.Error
	}

	return tracks, nil
}

// UpdateTrack обновляет данные трассы
func (r *ContentRepository) UpdateTrack(id uint, track *Track) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := db.WithContext(ctx).Model(&Track{}).Where("id = ?", id).Updates(track)
	if result.Error != nil {
		log.Printf("ERROR: Failed to update track %d: %v", id, result.Error)
		return result.Error
	}

	log.Printf("Track updated successfully: ID=%d", id)
	return nil
}

// UpdateTrackName обновляет название трассы
func (r *ContentRepository) UpdateTrackName(id uint, name string) (uint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := db.WithContext(ctx).Model(&Track{}).Where("id = ?", id).Update("name", name)
	if result.Error != nil {
		log.Printf("ERROR: Failed to update track name %d: %v", id, result.Error)
		return 0, result.Error
	}

	log.Printf("Track name updated successfully: ID=%d", id)
	return id, nil
}

// UpdateTrackInfo обновляет информацию о трассе
func (r *ContentRepository) UpdateTrackInfo(id uint, info string) (uint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := db.WithContext(ctx).Model(&Track{}).Where("id = ?", id).Update("info", info)
	if result.Error != nil {
		log.Printf("ERROR: Failed to update track info %d: %v", id, result.Error)
		return 0, result.Error
	}

	log.Printf("Track info updated successfully: ID=%d", id)
	return id, nil
}

// DeleteTrack удаляет трассу
func (r *ContentRepository) DeleteTrack(id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := db.WithContext(ctx).Delete(&Track{}, id)
	if result.Error != nil {
		log.Printf("ERROR: Failed to delete track %d: %v", id, result.Error)
		return result.Error
	}

	log.Printf("Track deleted successfully: ID=%d", id)
	return nil
}

// ============================================================================
// МЕТОДЫ ДЛЯ РАБОТЫ С ПОЛЬЗОВАТЕЛЯМИ
// ============================================================================

// CreateUser создает нового пользователя
func (r *ContentRepository) CreateUser(user *User) (uint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("Creating user: %s (ChatId: %d)", user.Name, user.ChatId)

	result := db.WithContext(ctx).Create(user)
	if result.Error != nil {
		log.Printf("ERROR: Failed to create user %s: %v", user.Name, result.Error)
		return 0, result.Error
	}

	log.Printf("User created successfully: %s (ID: %d)", user.Name, user.ID)
	return user.ID, nil
}

// GetUserByID получает пользователя по ID
func (r *ContentRepository) GetUserByID(id uint) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user User
	result := db.WithContext(ctx).First(&user, id)
	if result.Error != nil {
		log.Printf("ERROR: Failed to get user %d: %v", id, result.Error)
		return nil, result.Error
	}

	return &user, nil
}

// GetUserByChatId получает пользователя по Chat ID
func (r *ContentRepository) GetUserByChatId(chatId int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user User
	result := db.WithContext(ctx).Where("chat_id = ? AND is_active = ?", chatId, true).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("User with chat_id %d not found", chatId)
		} else {
			log.Printf("ERROR: Failed to get user by chat_id %d: %v", chatId, result.Error)
		}
		return nil, result.Error
	}

	return &user, nil
}

// GetUsers получает всех пользователей
func (r *ContentRepository) GetUsers() ([]User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var users []User
	result := db.WithContext(ctx).Find(&users)
	if result.Error != nil {
		log.Printf("ERROR: Failed to get users: %v", result.Error)
		return nil, result.Error
	}

	return users, nil
}

// UpdateUser обновляет данные пользователя
func (r *ContentRepository) UpdateUser(id uint, user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Updates(user)
	if result.Error != nil {
		log.Printf("ERROR: Failed to update user %d: %v", id, result.Error)
		return result.Error
	}

	log.Printf("User updated successfully: ID=%d", id)
	return nil
}

// DeleteUser удаляет пользователя
func (r *ContentRepository) DeleteUser(id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := db.WithContext(ctx).Delete(&User{}, id)
	if result.Error != nil {
		log.Printf("ERROR: Failed to delete user %d: %v", id, result.Error)
		return result.Error
	}

	log.Printf("User deleted successfully: ID=%d", id)
	return nil
}

// UpdateUserELORating обновляет рейтинг ELO пользователя
func (r *ContentRepository) UpdateUserELORating(id uint, newRating int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Update("elo_rating", newRating)
	if result.Error != nil {
		log.Printf("ERROR: Failed to update user ELO rating %d: %v", id, result.Error)
		return result.Error
	}

	log.Printf("User ELO rating updated successfully: ID=%d, New Rating=%d", id, newRating)
	return nil
}

// ============================================================================
// МЕТОДЫ ДЛЯ РАБОТЫ С РЕГИСТРАЦИЯМИ НА ТРЕНИРОВКИ
// ============================================================================

// CreateTrainingRegistration создает новую регистрацию на тренировку
func (r *ContentRepository) CreateTrainingRegistration(registration *TrainingRegistration) (uint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("Creating training registration: TrainingID=%d, UserID=%d", registration.TrainingID, registration.UserID)

	result := db.WithContext(ctx).Create(registration)
	if result.Error != nil {
		log.Printf("ERROR: Failed to create training registration: %v", result.Error)
		return 0, result.Error
	}

	log.Printf("Training registration created successfully: ID=%d", registration.ID)
	return registration.ID, nil
}

// GetTrainingRegistrationByID получает регистрацию по ID
func (r *ContentRepository) GetTrainingRegistrationByID(id uint) (*TrainingRegistration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var registration TrainingRegistration
	result := db.WithContext(ctx).First(&registration, id)
	if result.Error != nil {
		log.Printf("ERROR: Failed to get training registration %d: %v", id, result.Error)
		return nil, result.Error
	}

	return &registration, nil
}

// GetTrainingRegistrationsByTrainingID получает все регистрации на тренировку
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

// GetTrainingRegistrationsByUserID получает все регистрации пользователя
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

// UpdateTrainingRegistration обновляет регистрацию на тренировку
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

// DeleteTrainingRegistration удаляет регистрацию на тренировку
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

// GetTrainingRegistrationByUserAndTraining получает регистрацию пользователя на конкретную тренировку
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

// ============================================================================
// МЕТОДЫ ДЛЯ РАБОТЫ С АДМИНИСТРАТОРАМИ
// ============================================================================

// CreateAdmin создает нового администратора
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

// GetAdminByID получает администратора по ID
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

// GetAdminByChatId получает администратора по Chat ID
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

// GetAdmins получает всех администраторов
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

// UpdateAdmin обновляет данные администратора
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

// DeleteAdmin удаляет администратора
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

// ============================================================================
// МЕТОДЫ ДЛЯ РАБОТЫ СО СПОРТИВНЫМИ ТЕСТАМИ
// ============================================================================

// GetActiveSportsTests получает только активные спортивные тесты
func (r *ContentRepository) GetActiveSportsTests() ([]SportsTest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var tests []SportsTest
	result := db.WithContext(ctx).Where("is_active = ?", true).Find(&tests)
	if result.Error != nil {
		log.Printf("ERROR: Failed to get active sports tests: %v", result.Error)
		return nil, result.Error
	}

	return tests, nil
}

// InitDefaultSportsTests инициализирует базовые спортивные тесты
func (r *ContentRepository) InitDefaultSportsTests() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Проверяем, есть ли уже тесты
	var count int64
	db.WithContext(ctx).Model(&SportsTest{}).Count(&count)
	if count > 0 {
		return nil // Тесты уже инициализированы
	}

	// Создаем базовые тесты
	defaultTests := []SportsTest{
		{
			Name:        "Тест на выносливость",
			Description: "Бег на 5 км на время. Проверяет общую физическую подготовку и выносливость.",
			MaxScore:    100,
			IsActive:    true,
		},
		{
			Name:        "Тест на скорость",
			Description: "Спринт 100 метров. Проверяет скоростные качества и реакцию.",
			MaxScore:    100,
			IsActive:    true,
		},
		{
			Name:        "Тест на координацию",
			Description: "Полоса препятствий. Проверяет координацию движений и ловкость.",
			MaxScore:    100,
			IsActive:    true,
		},
		{
			Name:        "Тест на силу",
			Description: "Подтягивания, отжимания, приседания. Проверяет силовые показатели.",
			MaxScore:    100,
			IsActive:    true,
		},
	}

	for _, test := range defaultTests {
		result := db.WithContext(ctx).Create(&test)
		if result.Error != nil {
			log.Printf("ERROR: Failed to create default sports test %s: %v", test.Name, result.Error)
			return result.Error
		}
		log.Printf("Default sports test created: %s", test.Name)
	}

	return nil
}

// GetActiveTrainingsByTrackAndTrainer получает активные тренировки по трассе и тренеру
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

// GetTrainersByTrack получает тренеров, у которых есть активные тренировки на указанной трассе
func (r *ContentRepository) GetTrainersByTrack(trackId uint) ([]Trainer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var trainers []Trainer
	result := db.WithContext(ctx).
		Table("trainers").
		Select("DISTINCT trainers.*").
		Joins("INNER JOIN trainings ON trainers.id = trainings.trainer_id").
		Where("trainings.track_id = ? AND trainings.is_active = ?", trackId, true).
		Find(&trainers)

	if result.Error != nil {
		log.Printf("ERROR: Failed to get trainers by track: %v", result.Error)
		return nil, result.Error
	}

	return trainers, nil
}

// GetTracksWithActiveTrainings получает трассы, на которых есть активные тренировки
func (r *ContentRepository) GetTracksWithActiveTrainings() ([]Track, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var tracks []Track
	result := db.WithContext(ctx).
		Table("tracks").
		Select("DISTINCT tracks.*").
		Joins("INNER JOIN trainings ON tracks.id = trainings.track_id").
		Where("trainings.is_active = ?", true).
		Find(&tracks)

	if result.Error != nil {
		log.Printf("ERROR: Failed to get tracks with active trainings: %v", result.Error)
		return nil, result.Error
	}

	return tracks, nil
}

// ============================================================================
// МЕТОДЫ ДЛЯ РАБОТЫ С РАСПИСАНИЕМ
// ============================================================================
