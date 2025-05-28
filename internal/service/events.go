package service

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
)

type EventsService struct {
	BaseService[models.Event]
	repo repository.EventRepository
}

func NewEventsService() *EventsService {
	repo := repository.NewEventRepository()
	return &EventsService{
		BaseService: NewBaseService(repo),
		repo:        *repo,
	}
}

func (s *EventsService) AddMember(eventId int, memberId int) (*models.Event, error) {
	return s.repo.AddMember(eventId, memberId)
}

func (s *EventsService) RemoveMember(eventId int, memberId int) (*models.Event, error) {
	return s.repo.RemoveMember(eventId, memberId)
}
