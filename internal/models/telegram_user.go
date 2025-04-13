package models

import (
	"time"
)

type TelegramUser struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	TelegramID int64     `json:"telegram_id" gorm:"uniqueIndex"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
} 