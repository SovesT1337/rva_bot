package database

// Интерфейс для работы с контентом
type ContentRepositoryInterface interface {
	CreateTrainer(content *Trainer) (uint, error)
	GetTrainerByID(ID uint) (*Trainer, error)
	GetTrainers() ([]Trainer, error)
	UpdateTrainerTgId(id uint, tgid string) (uint, error)
	UpdateTrainerName(id uint, name string) (uint, error)
	UpdateTrainerInfo(id uint, info string) (uint, error)
	UpdateTrainer(id uint, trainer *Trainer) error
	DeleteTrainer(id uint) error

	// Методы для работы с администраторами
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
	UpdateTrackName(id uint, name string) (uint, error)
	UpdateTrackInfo(id uint, info string) (uint, error)
	DeleteTrack(id uint) error

	// Методы для работы с пользователями
	CreateUser(user *User) (uint, error)
	GetUserByID(id uint) (*User, error)
	GetUserByChatId(chatId int) (*User, error)
	GetUsers() ([]User, error)
	UpdateUser(id uint, user *User) error
	DeleteUser(id uint) error

	// Методы для работы с тренировками
	CreateTraining(content *Training) (uint, error)
	GetTrainingById(id uint) (*Training, error)
	GetTrainings() ([]Training, error)
	GetActiveTrainings() ([]Training, error)
	UpdateTraining(id uint, training *Training) error
	DeleteTraining(id uint) error

	// Методы для работы с регистрациями на тренировки
	CreateTrainingRegistration(registration *TrainingRegistration) (uint, error)
	GetTrainingRegistrationByID(id uint) (*TrainingRegistration, error)
	GetTrainingRegistrationsByTrainingID(trainingId uint) ([]TrainingRegistration, error)
	GetTrainingRegistrationsByUserID(userId uint) ([]TrainingRegistration, error)
	UpdateTrainingRegistration(id uint, registration *TrainingRegistration) error
	DeleteTrainingRegistration(id uint) error
	GetTrainingRegistrationByUserAndTraining(userId uint, trainingId uint) (*TrainingRegistration, error)

	// Методы для работы со спортивными тестами
	GetActiveSportsTests() ([]SportsTest, error)
	InitDefaultSportsTests() error

	// Методы для пошаговой записи на тренировки
	GetActiveTrainingsByTrackAndTrainer(trackId, trainerId uint) ([]Training, error)
	GetTrainersByTrack(trackId uint) ([]Trainer, error)
	GetTracksWithActiveTrainings() ([]Track, error)
}

// Создание экземпляра репозитория
func NewContentRepository() ContentRepositoryInterface {
	return &ContentRepository{}
}
