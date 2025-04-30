package repository

import (
	"fmt"
	"ithozyeva/database"
	"ithozyeva/internal/models"
)

// Изменяем с type alias на новый тип
type MemberRepositoryInterface interface {
	BaseRepository[models.Member]

	GetByTelegramID(telegramID int64) (*models.Member, error)
}
type MemberRepository struct {
	BaseRepository[models.Member]
}

func NewMemberRepository() *MemberRepository {
	return &MemberRepository{
		BaseRepository: NewBaseRepository(database.DB, &models.Member{}),
	}
}

func (r *MemberRepository) GetMemberByTelegram(telegram string) (*models.Member, error) {
	var member models.Member
	result := database.DB.Where("username = ?", telegram).First(&member)
	if result.Error != nil {
		return nil, result.Error
	}
	return &member, nil
}

func (r *MemberRepository) GetByTelegramID(telegramID int64) (*models.Member, error) {
	entity := new(models.Member)
	if err := database.DB.Where("telegram_id = ?", telegramID).First(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

// GetById получает участника по ID с проверкой на статус ментора
func (r *MemberRepository) GetById(id int64) (*models.Member, error) {
	var member models.Member
	if err := database.DB.First(&member, id).Error; err != nil {
		return nil, err
	}

	result := &models.Member{
		Id:         member.Id,
		Username:   member.Username,
		FirstName:  member.FirstName,
		TelegramID: member.TelegramID,
		LastName:   member.LastName,
		Role:       member.Role,
		Birthday:   member.Birthday,
	}

	return result, nil
}

func (r *MemberRepository) Update(member *models.Member) (*models.Member, error) {
	result := database.DB.Model(&models.Member{}).
		Where("id = ?", member.Id).
		Update("birthday", member.Birthday).
		Update("first_name", member.FirstName).
		Update("last_name", member.LastName).
		Update("role", member.Role).
		Update("username", member.Username)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("member not found")
	}

	member, err := r.GetMemberByTelegram(member.Username)
	if err != nil {
		return nil, err
	}

	return member, nil
}

func (r *MemberRepository) GetTodayBirthdays() ([]string, error) {
	query := `
		SELECT 
			username
		FROM members
		WHERE
			role = ?
			AND
    		DATE_PART('day', birthday) = date_part('day', CURRENT_DATE)
			AND
    		DATE_PART('month', birthday) = date_part('month', CURRENT_DATE)
	`

	rows, err := database.DB.Raw(query, models.MemberRoleSubscriber).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usernames []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, err
		}
		usernames = append(usernames, username)
	}
	return usernames, nil
}
