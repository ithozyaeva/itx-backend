package handler

import (
	"fmt"
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"ithozyeva/internal/service"
	"ithozyeva/internal/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type EventsHandler struct {
	BaseHandler[models.Event]
	svc *service.EventsService
}

func NewEventsHandler() *EventsHandler {
	svc := service.NewEventsService()
	return &EventsHandler{
		BaseHandler: *NewBaseHandler(svc),
		svc:         svc,
	}
}

var EventsSearchFields = map[string]string{
	"dateFrom": "date >= ?",
	"dateTo":   "date < ?",
}

type EventsSearchRequest struct {
	Limit    *int    `query:"limit"`
	Offset   *int    `query:"offset"`
	DateFrom *string `query:"dateFrom"`
	DateTo   *string `query:"dateTo"`
}

func (h *EventsHandler) Search(c *fiber.Ctx) error {
	req := new(EventsSearchRequest)
	if err := c.QueryParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	filter := make(repository.SearchFilter)
	// *filter = make(map[string]interface{})

	if req.DateFrom != nil {
		filter[EventsSearchFields["dateFrom"]] = *req.DateFrom
	}
	if req.DateTo != nil {
		filter[EventsSearchFields["dateTo"]] = *req.DateTo
	}

	result, err := h.service.Search(req.Limit, req.Offset, &filter, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

func (h *EventsHandler) GetOld(c *fiber.Ctx) error {

	result, err := h.service.Search(nil, nil, &repository.SearchFilter{
		"date < ?": gorm.Expr("CURRENT_TIMESTAMP"),
	}, &repository.Order{
		ColumnBy: "date",
		Order:    "DESC",
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

func (h *EventsHandler) GetNext(c *fiber.Ctx) error {
	result, err := h.service.Search(nil, nil, &repository.SearchFilter{
		"date >= ?": gorm.Expr("CURRENT_TIMESTAMP"),
	}, &repository.Order{
		ColumnBy: "date",
		Order:    "ASC",
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

func (h *EventsHandler) AddMember(c *fiber.Ctx) error {
	req := new(WorkWithEventRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	member := c.Locals("member").(*models.Member)

	result, err := h.svc.AddMember(req.EventId, int(member.Id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

type WorkWithEventRequest struct {
	EventId int `json:"eventId" query:"eventId"`
}

func (h *EventsHandler) RemoveMember(c *fiber.Ctx) error {
	req := new(WorkWithEventRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	member := c.Locals("member").(*models.Member)

	result, err := h.svc.RemoveMember(req.EventId, int(member.Id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

func (h *EventsHandler) GetICSFile(c *fiber.Ctx) error {
	req := new(WorkWithEventRequest)
	if err := c.QueryParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	event, err := h.svc.GetById(int64(req.EventId))
	if err != nil {

	}

	ics := utils.GenerateICS(event)

	c.Set("Content-Type", "text/calendar")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=event_%d.ics", event.Id))
	return c.SendString(ics)
}
