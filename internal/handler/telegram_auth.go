package handler

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"ithozyeva/internal/service"

	"github.com/gofiber/fiber/v2"
)

type TelegramAuthHandler struct {
	telegramService *service.TelegramService
	userRepo        *repository.TelegramUserRepository
}

func NewTelegramAuthHandler(telegramService *service.TelegramService, userRepo *repository.TelegramUserRepository) *TelegramAuthHandler {
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

	// Проверяем, существует ли пользователь
	existingUser, err := h.userRepo.GetByToken(req.Token)
	if err != nil {
		// Если пользователь не существует, создаем нового
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	} else if !existingUser.IsAuthenticated {
		// Если пользователь существует - подтверждаем аутентификацию
		h.userRepo.Update(&models.TelegramUser{
			TelegramID:      existingUser.TelegramID,
			ID:              existingUser.ID,
			Username:        existingUser.Username,
			FirstName:       existingUser.FirstName,
			LastName:        existingUser.LastName,
			IsAuthenticated: true,
			Token:           existingUser.Token,
		})
	}

	// Генерируем новый токен для пользователя
	authToken, err := h.telegramService.GenerateAuthToken(existingUser.TelegramID)
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

type CreateUserReq struct {
	Token     string `json:"token"`
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (h *TelegramAuthHandler) CreateUser(c *fiber.Ctx) error {
	var req CreateUserReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Проверяем, существует ли пользователь
	existingUser, err := h.userRepo.GetByToken(req.Token)
	if err != nil {
		// Если пользователь не существует, создаем нового
		newUser := &models.TelegramUser{
			TelegramID:      req.UserID,
			Username:        req.Username,
			FirstName:       req.FirstName,
			LastName:        req.LastName,
			IsAuthenticated: false,
			Token:           req.Token,
		}

		createdUser, err := h.userRepo.Create(newUser)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create user",
			})
		}

		existingUser = createdUser
	}

	return c.JSON(existingUser)
}
