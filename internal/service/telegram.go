package service

import (
	"fmt"
	"ithozyeva/config"
	"ithozyeva/internal/utils"
)

type TelegramService struct {
	botToken string
}

func NewTelegramService() (*TelegramService, error) {
	botToken := config.CFG.TelegramToken
	if botToken == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN is not set")
	}

	return &TelegramService{
		botToken: botToken,
	}, nil
}

type TelegramUser struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type TelegramAuthResponse struct {
	OK     bool         `json:"ok"`
	Result TelegramUser `json:"result"`
}

// GenerateAuthToken генерирует токен для авторизации
func (s *TelegramService) GenerateAuthToken(userID int64) string {
	token := utils.HashToken(fmt.Sprintf("%d", userID))

	return token
}
