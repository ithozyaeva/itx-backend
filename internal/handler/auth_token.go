package handler

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/service"

	"github.com/gofiber/fiber/v2"
)

type TelegramAuthHandler struct {
	telegramService *service.TelegramService
	authService     *service.AuthTokenService
}

func NewTelegramAuthHandler() *TelegramAuthHandler {
	tgService, err := service.NewTelegramService()

	if err != nil {
		panic(err)
	}

	return &TelegramAuthHandler{
		telegramService: tgService,
		authService:     service.NewAuthTokenService(),
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

	// Проверяем, существует ли токен
	existingToken, existingUser, err := h.authService.GetByToken(req.Token)
	if err != nil {
		// Если токена не существует, создаем нового
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Добавляем заголовок
	c.Response().Header.Add("X-Telegram-User-Token", existingToken.Token)

	// Возвращаем информацию о пользователе и токен
	return c.JSON(fiber.Map{
		"user":  existingUser,
		"token": existingToken.Token,
	})
}

type HandleBotMessageReq struct {
	Token     string `json:"token"`
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (h *TelegramAuthHandler) HandleBotMessage(c *fiber.Ctx) error {
	var req HandleBotMessageReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Проверяем, существует ли пользователь
	existingUser, err := h.authService.GetByTelegramID(req.UserID)
	if err != nil {
		// Если пользователь не существует, создаем нового
		newUser := &models.Member{
			TelegramID: req.UserID,
			Username:   req.Username,
			FirstName:  req.FirstName,
			LastName:   req.LastName,
		}

		createdUser, err := h.authService.CreateNewMember(newUser, req.Token)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create user",
			})
		}
		existingUser = createdUser
	} else {
		_, err := h.authService.CreateOrUpdateToken(req.UserID, req.Token)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get auth token",
			})
		}
	}

	return c.JSON(existingUser)
}
