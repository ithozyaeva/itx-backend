package middleware

import (
	"ithozyeva/internal/repository"
	"ithozyeva/internal/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AuthMiddleware struct {
	userRepo *repository.AuthTokenRepository
}

func NewAuthMiddleware(db *gorm.DB) *AuthMiddleware {
	return &AuthMiddleware{
		userRepo: repository.NewAuthTokenRepository(),
	}
}

func (m *AuthMiddleware) RequireAuth(c *fiber.Ctx) error {
	telegramToken := c.Get("X-Telegram-User-Token")
	if telegramToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	authToken, err := m.userRepo.GetByToken(telegramToken)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Telegram User ID",
		})
	}

	if utils.CheckExpirationDate(authToken.ExpiredAt) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	user, err := m.userRepo.GetByTelegramID(authToken.TelegramID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Добавляем информацию о пользователе в контекст
	c.Locals("user", user)

	return c.Next()
}
