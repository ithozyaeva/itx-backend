package handler

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"ithozyeva/internal/service"

	"github.com/gofiber/fiber/v2"
)

type TelegramAuthHandler struct {
	telegramService *service.TelegramService
	userRepo        repository.BaseRepository[models.TelegramUser]
}

func NewTelegramAuthHandler(telegramService *service.TelegramService, userRepo repository.BaseRepository[models.TelegramUser]) *TelegramAuthHandler {
	return &TelegramAuthHandler{
		telegramService: telegramService,
		userRepo:        userRepo,
	}
}

type AuthRequest struct {
	Token string `json:"token"`
}

func (h *TelegramAuthHandler) Authenticate(c *fiber.Ctx) error {
	var req AuthRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Получаем информацию о пользователе из Telegram
	userInfo, err := h.telegramService.GetUserInfo(c.Context(), req.Token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Failed to authenticate with Telegram",
		})
	}

	// Проверяем, существует ли пользователь
	existingUser, err := h.userRepo.GetByTelegramID(userInfo.ID)
	if err != nil {
		// Если пользователь не существует, создаем нового
		newUser := &models.TelegramUser{
			TelegramID: userInfo.ID,
			Username:   userInfo.Username,
			FirstName:  userInfo.FirstName,
			LastName:   userInfo.LastName,
		}
		createdUser, err := h.userRepo.Create(newUser)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create user",
			})
		}
		existingUser = createdUser
	}

	// Генерируем новый токен для пользователя
	authToken, err := h.telegramService.GenerateAuthToken(userInfo.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate auth token",
		})
	}

	// Возвращаем информацию о пользователе и токен
	return c.JSON(fiber.Map{
		"user":  existingUser,
		"token": authToken,
	})
} 