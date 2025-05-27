package handler

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"ithozyeva/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type EventsHandler struct {
	BaseHandler[models.Event]
}

func NewEventsHandler() *EventsHandler {
	return &EventsHandler{
		BaseHandler: *NewBaseHandler(service.NewEventsService()),
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
