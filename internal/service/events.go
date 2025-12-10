package service

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"time"
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

// GetFutureEvents получает только будущие события
func (s *EventsService) GetFutureEvents(now time.Time) ([]models.Event, error) {
	filter := &repository.SearchFilter{
		"date >= ?": now,
	}
	result, err := s.Search(nil, nil, filter, nil)
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}
