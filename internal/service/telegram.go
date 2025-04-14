package service

import (
	"fmt"
	"ithozyeva/config"
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
func (s *TelegramService) GenerateAuthToken(userID int64) (string, error) {
	// Здесь можно использовать JWT или другой механизм генерации токена
	// Для простоты используем хеш от ID пользователя
	token := fmt.Sprintf("tg_%d", userID)
	return token, nil
}

// GetBotUrl возвращает URL бота с параметрами для авторизации
func (s *TelegramService) GetBotUrl(redirectUrl string) string {
	return fmt.Sprintf("https://t.me/itx_welcome_bot?start=%s", redirectUrl)
}
