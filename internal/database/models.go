package database

import (
	"time"
)

type Trainer struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	TgId      string
	ChatId    int `gorm:"uniqueIndex"`
	Info      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Admin struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	TgId      string
	ChatId    int `gorm:"uniqueIndex"`
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	TgId      string
	ChatId    int `gorm:"uniqueIndex"`
	IsActive  bool
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
	StartTime       time.Time
	EndTime         time.Time
	MaxParticipants int
	IsActive        bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type TrainingRegistration struct {
	ID         uint `gorm:"primaryKey"`
	TrainingID uint
	UserID     uint
	Status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type TrainingRequest struct {
	ID         uint `gorm:"primaryKey"`
	UserID     uint
	Message    string
	IsReviewed bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
