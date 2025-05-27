package repository

import (
	"fmt"
	"ithozyeva/database"
	"ithozyeva/internal/models"
)

type EventRepository struct {
	BaseRepository[models.Event]
}

func NewEventRepository() *EventRepository {
	return &EventRepository{
		BaseRepository: NewBaseRepository(database.DB, &models.Event{}),
	}
}

func (e *EventRepository) Search(limit *int, offset *int, filter *SearchFilter, order *Order) ([]models.Event, int64, error) {
	var events []models.Event
	var count int64

	if err := database.DB.Model(&models.Event{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	query := database.DB.Model(&models.Event{}).Preload("Hosts")

	if filter != nil {
		for key, value := range *filter {
			query = query.Where(key, value)
		}
	}

	if order != nil {
		query = query.Order(fmt.Sprintf("\"%s\" %s", order.ColumnBy, order.Order))
	}

	if limit != nil {
		query = query.Limit(*limit)
	}

	if offset != nil {
		query = query.Offset(*offset)
	}

	if err := query.Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, count, nil
}

// GetById получает отзыв по ID с информацией о услуге
func (r *EventRepository) GetById(id int64) (*models.Event, error) {
	var event models.Event
	if err := database.DB.Preload("Hosts").First(&event, id).Error; err != nil {
		return nil, err
	}
	return &event, nil
}
