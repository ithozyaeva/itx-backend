package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"ithozyeva/config"
	"ithozyeva/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	bot     *tgbotapi.BotAPI
	service *service.TelegramService
}

func NewTelegramBot() (*TelegramBot, error) {

	botToken := config.CFG.TelegramToken
	if botToken == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN is not set")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, fmt.Errorf("error creating bot: %v", err)
	}

	service, err := service.NewTelegramService()
	if err != nil {
		return nil, fmt.Errorf("error creating telegram service: %v", err)
	}

	return &TelegramBot{
		bot:     bot,
		service: service,
	}, nil
}

func (b *TelegramBot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				b.handleStartCommand(update.Message)
			}
		}
	}
}

func (b *TelegramBot) handleStartCommand(message *tgbotapi.Message) {
	log.Printf("Received /start command from user %d with args: %s", message.From.ID, message.CommandArguments())

	// Получаем аргументы команды
	args := strings.Split(message.CommandArguments(), " ")
	if len(args) == 0 || args[0] == "" {
		log.Printf("No arguments provided for /start command")
		b.sendMessage(message.Chat.ID, "Пожалуйста, используйте кнопку авторизации на сайте")
		return
	}

	// Первый аргумент - URL для перенаправления
	redirectUrl := config.CFG.PublicDomain
	log.Printf("Redirect URL before processing: %s", redirectUrl)

	if !strings.HasPrefix(redirectUrl, "http://") && !strings.HasPrefix(redirectUrl, "https://") {
		redirectUrl = "http://" + redirectUrl
	}
	log.Printf("Final redirect URL: %s", redirectUrl)

	// Генерируем токен для пользователя
	token := b.service.GenerateAuthToken(message.From.ID)

	log.Printf("Generated token for user %d: %s", message.From.ID, token)

	// Формируем URL для перенаправления с токеном
	authUrl := fmt.Sprintf("%s?token=%s", redirectUrl, token)
	log.Printf("Final auth URL: %s", authUrl)

	sendAuthToBackend(token, message.From)

	// Отправляем сообщение с кнопкой для авторизации
	msg := tgbotapi.NewMessage(message.Chat.ID, "Нажмите кнопку ниже для авторизации")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Авторизоваться", authUrl),
		),
	)

	if _, err := b.bot.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	} else {
		log.Printf("Successfully sent auth button to user %d", message.From.ID)
	}
}

func (b *TelegramBot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := b.bot.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

type AuthRequest struct {
	Token     string `json:"token"`
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func sendAuthToBackend(token string, user *tgbotapi.User) {
	data := AuthRequest{
		Token:     token,
		UserID:    user.ID,
		Username:  user.UserName,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("Ошибка сериализации JSON:", err)
		return
	}

	url := fmt.Sprintf("%s/api/auth/telegramFromBot", config.CFG.BackendDomain)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Ошибка отправки запроса:", err)
		return
	}
	defer resp.Body.Close()

	log.Println("Ответ от Fiber:", resp.Status)
}
