package service

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
)

type EventsService = BaseService[models.Event]

func NewEventsService() EventsService {
	repo := repository.NewEventRepository()
	return NewBaseService(repo)
}
