package service

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
)

// MemberRepoAdapter адаптер для репозитория участников
type MemberRepoAdapter struct {
	*repository.MemberRepository
}

// MemberService представляет сервис для работы с участниками
type MemberService struct {
	BaseService[models.Member]
	repo       *repository.MemberRepository
	mentorRepo *repository.MentorRepository
}

// NewMemberService создает новый экземпляр сервиса участников
func NewMemberService() *MemberService {
	repo := repository.NewMemberRepository()
	adapter := &MemberRepoAdapter{repo}

	return &MemberService{
		BaseService: NewBaseService[models.Member](adapter),
		repo:        repo,
		mentorRepo:  repository.NewMentorRepository(),
	}
}

func (s *MemberService) GetTodayBirthdays() ([]string, error) {
	return s.repo.GetTodayBirthdays()
}

func (s *MemberService) GetMentor(memberId int64) (*models.MentorModel, error) {
	mentorDb, err := s.mentorRepo.GetByMemberID(memberId)

	if err != nil {
		return nil, err
	}

	mentor := mentorDb.ToModel()

	return &mentor, nil
}
