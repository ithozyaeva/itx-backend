package middleware

import (
	"ithozyeva/config"
	"ithozyeva/internal/repository"
	"ithozyeva/internal/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthMiddleware struct {
	userRepo   *repository.AuthTokenRepository
	memberRepo *repository.MemberRepository
}

func NewAuthMiddleware(db *gorm.DB) *AuthMiddleware {
	return &AuthMiddleware{
		userRepo:   repository.NewAuthTokenRepository(),
		memberRepo: repository.NewMemberRepository(),
	}
}

func (m *AuthMiddleware) RequireJWTAuth(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	if tokenStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.CFG.JwtSecret), nil
	})

	if err == nil && token.Valid {
		return c.Next()
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "Unauthorized",
	})
}

func (m *AuthMiddleware) RequireTGAuth(c *fiber.Ctx) error {
	// Если JWT токен не валиден или отсутствует, проверяем Telegram токен
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

	member, err := m.memberRepo.GetByTelegramID(authToken.TelegramID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Member not found",
		})
	}

	// Добавляем информацию о пользователе в контекст
	c.Locals("member", member)

	return c.Next()
}
