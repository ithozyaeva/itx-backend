package handler

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/service"
	"ithozyeva/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// MembersHandler обработчик для работы с участниками
type MembersHandler struct {
	svc *service.MemberService
}

// NewMembersHandler создает новый экземпляр обработчика участников
func NewMembersHandler() *MembersHandler {
	svc := service.NewMemberService()
	return &MembersHandler{
		svc: svc,
	}
}

// Search выполняет поиск участников с пагинацией
func (h *MembersHandler) Search(c *fiber.Ctx) error {
	req := new(models.SearchRequest)
	if err := c.QueryParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	result, err := h.svc.Search(req.Limit, req.Offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

// GetById получает участника по ID
func (h *MembersHandler) GetById(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	result, err := h.svc.GetById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

// Реализуем все необходимые методы напрямую
func (h *MembersHandler) Create(c *fiber.Ctx) error {
	request := new(models.Member)
	err := c.BodyParser(request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	result, err := h.svc.Create(request)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

type UpdateRequest struct {
	Id        int64             `json:"id"`
	FirstName string            `json:"firstName"`
	LastName  string            `json:"lastName"`
	Birthday  *string           `json:"birthday"`
	Role      models.MemberRole `json:"role"`
}

func (h *MembersHandler) Update(c *fiber.Ctx) error {
	request := new(UpdateRequest)
	err := c.BodyParser(request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	member, err := h.svc.GetById(request.Id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Участник не найден"})
	}

	member.FirstName = request.FirstName
	member.LastName = request.LastName
	member.Role = request.Role

	parsedDate, err := utils.ParseDate(request.Birthday)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	member.Birthday = parsedDate

	result, err := h.svc.Update(member)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

type UpdateProfileRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Birthday  string `json:"birthday"`
}

func (h *MembersHandler) UpdateProfile(c *fiber.Ctx) error {
	request := new(UpdateRequest)
	err := c.BodyParser(request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	member := c.Locals("member").(*models.Member)
	member.FirstName = request.FirstName
	member.LastName = request.LastName

	parsedDate, err := utils.ParseDate(request.Birthday)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	member.Birthday = parsedDate

	result, err := h.svc.Update(member)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

func (h *MembersHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	entity := new(models.Member)

	entity.Id = int64(id)

	if err := h.svc.Delete(entity); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *MembersHandler) Me(c *fiber.Ctx) error {
	member := c.Locals("member").(*models.Member)
	return c.JSON(member)
}
