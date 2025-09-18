package database

import (
	"time"
)

type Trainer struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	TgId      string
	ChatId    int `gorm:"uniqueIndex"` // Telegram Chat ID для уведомлений тренера
	Info      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Admin struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	TgId      string
	ChatId    int  `gorm:"uniqueIndex"` // Telegram Chat ID для проверки прав администратора
	IsActive  bool // Активен ли администратор
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	ChatId    int  `gorm:"uniqueIndex"` // Telegram Chat ID пользователя
	IsActive  bool // Активен ли пользователь
	ELORating int  `gorm:"default:1200"` // Рейтинг ELO (по умолчанию 1200)
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Track struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Info      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Training struct {
	ID              uint `gorm:"primaryKey"`
	TrainerID       uint
	TrackID         uint
	Time            time.Time
	MaxParticipants int  // Максимальное количество участников
	IsActive        bool // Активна ли тренировка для регистрации
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type TrainingRegistration struct {
	ID         uint `gorm:"primaryKey"`
	TrainingID uint
	UserID     uint
	Status     string // "pending", "confirmed", "rejected"
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
