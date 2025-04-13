package middleware

import (
	"ithozyeva/internal/repository"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AuthMiddleware struct {
	userRepo *repository.TelegramUserRepository
}

func NewAuthMiddleware(db *gorm.DB) *AuthMiddleware {
	return &AuthMiddleware{
		userRepo: repository.NewTelegramUserRepository(db),
	}
}

func (m *AuthMiddleware) RequireAuth(c *fiber.Ctx) error {
	telegramIDStr := c.Get("X-Telegram-User-ID")
	if telegramIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	telegramID, err := strconv.ParseInt(telegramIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Telegram User ID",
		})
	}

	user, err := m.userRepo.GetByTelegramID(telegramID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Добавляем информацию о пользователе в контекст
	c.Locals("user", user)

	return c.Next()
} 