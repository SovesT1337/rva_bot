package database

import (
	"gorm.io/gorm"
)

type ContentRepositoryInterface interface {
	CreateTrainer(content *Trainer) (uint, error)
	GetTrainerByID(ID uint) (*Trainer, error)
	GetTrainers() ([]Trainer, error)
	UpdateTrainer(id uint, trainer *Trainer) error
	DeleteTrainer(id uint) error

	CreateAdmin(admin *Admin) (uint, error)
	GetAdminByID(id uint) (*Admin, error)
	GetAdminByChatId(chatId int) (*Admin, error)
	GetAdmins() ([]Admin, error)
	UpdateAdmin(id uint, admin *Admin) error
	DeleteAdmin(id uint) error

	CreateTrack(content *Track) (uint, error)
	GetTrackByID(id uint) (*Track, error)
	GetTracks() ([]Track, error)
	UpdateTrack(id uint, track *Track) error
	DeleteTrack(id uint) error

	CreateUser(user *User) (uint, error)
	GetUserByID(id uint) (*User, error)
	GetUserByChatId(chatId int) (*User, error)
	GetUsers() ([]User, error)
	UpdateUser(id uint, user *User) error
	DeleteUser(id uint) error

	CreateTraining(content *Training) (uint, error)
	GetTrainingById(id uint) (*Training, error)
	GetTrainings() ([]Training, error)
	GetActiveTrainings() ([]Training, error)
	UpdateTraining(id uint, training *Training) error
	DeleteTraining(id uint) error

	CreateTrainingRegistration(registration *TrainingRegistration) (uint, error)
	GetTrainingRegistrationByID(id uint) (*TrainingRegistration, error)
	GetTrainingRegistrationsByTrainingID(trainingId uint) ([]TrainingRegistration, error)
	GetTrainingRegistrationsByUserID(userId uint) ([]TrainingRegistration, error)
	UpdateTrainingRegistration(id uint, registration *TrainingRegistration) error
	DeleteTrainingRegistration(id uint) error
	GetTrainingRegistrationByUserAndTraining(userId uint, trainingId uint) (*TrainingRegistration, error)

	GetActiveTrainingsByTrackAndTrainer(trackId, trainerId uint) ([]Training, error)
	GetTrainersByTrack(trackId uint) ([]Trainer, error)
	GetTracksWithActiveTrainings() ([]Track, error)
	GetUserTrainings(userId uint) ([]Training, error)

	CreateTrainingRequest(request *TrainingRequest) (uint, error)
	GetTrainingRequestByID(id uint) (*TrainingRequest, error)
	GetTrainingRequests() ([]TrainingRequest, error)
	GetUnreviewedTrainingRequests() ([]TrainingRequest, error)
	UpdateTrainingRequest(id uint, request *TrainingRequest) error
	DeleteTrainingRequest(id uint) error
}

type ContentRepository struct {
	db *gorm.DB
}

func NewContentRepository(database *Database) ContentRepositoryInterface {
	return &ContentRepository{
		db: database.GetDB(),
	}
}
