package bot

import (
	"fmt"
	"log"
	"os"
	"strings"

	"ithozyeva/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

type TelegramBot struct {
	bot     *tgbotapi.BotAPI
	service *service.TelegramService
}

func NewTelegramBot() (*TelegramBot, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
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
	redirectUrl := args[0]
	log.Printf("Redirect URL before processing: %s", redirectUrl)
	
	if !strings.HasPrefix(redirectUrl, "http://") && !strings.HasPrefix(redirectUrl, "https://") {
		redirectUrl = "http://" + redirectUrl
	}
	log.Printf("Final redirect URL: %s", redirectUrl)

	// Генерируем токен для пользователя
	token, err := b.service.GenerateAuthToken(message.From.ID)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		b.sendMessage(message.Chat.ID, "Произошла ошибка при авторизации")
		return
	}
	log.Printf("Generated token for user %d: %s", message.From.ID, token)

	// Формируем URL для перенаправления с токеном
	authUrl := fmt.Sprintf("%s?token=%s", redirectUrl, token)
	log.Printf("Final auth URL: %s", authUrl)

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