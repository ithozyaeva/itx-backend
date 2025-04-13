package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type TelegramService struct {
	botToken string
}

func NewTelegramService() (*TelegramService, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
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

func (s *TelegramService) GetUserInfo(ctx context.Context, token string) (*TelegramUser, error) {
	// Проверяем токен через Telegram API
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getUserInfo", s.botToken)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Добавляем токен в заголовок
	req.Header.Add("X-Telegram-Token", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	var result TelegramAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	if !result.OK {
		return nil, fmt.Errorf("telegram API error")
	}

	return &result.Result, nil
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