package main

import (
	"context"
	"fmt"
	"ithozyeva/config"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"ithozyeva/database"
	"ithozyeva/internal/repository"
)

type Bot struct {
	api        *tgbotapi.BotAPI
	memberRepo *repository.MemberRepository
	mu         sync.Mutex
}

func NewBot(token string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	return &Bot{
		api:        api,
		memberRepo: repository.NewMemberRepository(),
		mu:         sync.Mutex{},
	}, nil
}

func (b *Bot) Start(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			return nil
		case update := <-updates:
			if update.Message != nil {
				if update.Message.NewChatMembers != nil {
					for _, user := range update.Message.NewChatMembers {
						go b.HandleNewChatMember(&user)
					}
				}
				if update.Message.LeftChatMember != nil {
					go b.HandleLeftChatMember(update.Message.LeftChatMember)
				}
			}
		}
	}
}

func (b *Bot) HandleNewChatMember(user *tgbotapi.User) {
	b.mu.Lock()
	defer b.mu.Unlock()

	member, err := b.memberRepo.CreateMemberIfNotExists(strconv.FormatInt(user.ID, 10), user.UserName)
	if err != nil {
		log.Printf("Error creating/updating user: %v", err)
	}

	if !member.IsActive {
		if err := b.memberRepo.UpdateMemberStatus(strconv.FormatInt(user.ID, 10), true); err != nil {
			log.Printf("Error updating user status: %v", err)
		}
	}
}

func (b *Bot) HandleLeftChatMember(user *tgbotapi.User) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if err := b.memberRepo.UpdateMemberStatus(strconv.FormatInt(user.ID, 10), false); err != nil {
		log.Printf("Error updating user status: %v", err)
	}

}

func main() {
	// Load environment variables
	config.LoadSyncUsersConfig()

	if err := database.SetupSyncUsersDatabase(); err != nil {
		log.Fatalf("Database setup failed: %v", err)
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}

	// Create and start bot
	bot, err := NewBot(token)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	log.Println("Bot started")
	if err := bot.Start(context.Background()); err != nil {
		log.Fatalf("Bot stopped with error: %v", err)
	}
}
