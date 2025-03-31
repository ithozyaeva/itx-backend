package service

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
)

type ReviewOnServiceService struct {
	BaseService[models.ReviewOnService]
	repository *repository.ReviewOnServiceRepository
}

func NewReviewOnServiceService() *ReviewOnServiceService {
	repo := repository.NewReviewOnServiceRepository()
	return &ReviewOnServiceService{
		BaseService: NewBaseService(repo),
		repository:  repo,
	}
}

func (s *ReviewOnServiceService) Search(limit *int, offset *int) (*models.RegistrySearch[models.ReviewOnService], error) {
	items, total, err := s.repository.Search(limit, offset)
	if err != nil {
		return nil, err
	}

	return &models.RegistrySearch[models.ReviewOnService]{
		Items: items,
		Total: int(total),
	}, nil
}

func (s *ReviewOnServiceService) GetReviewsWithMentorInfo(limit *int, offset *int) (*models.RegistrySearch[models.ReviewOnServiceWithMentor], error) {
	items, total, err := s.repository.GetReviewsWithMentorInfo(limit, offset)
	if err != nil {
		return nil, err
	}

	return &models.RegistrySearch[models.ReviewOnServiceWithMentor]{
		Items: items,
		Total: int(total),
	}, nil
}
