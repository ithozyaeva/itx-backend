package repository

import (
	"fmt"
	"ithozyeva/database"
	"ithozyeva/internal/models"
	"time"
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

// Search выполняет поиск участников с пагинацией и проверкой на статус ментора
func (r *MemberRepository) Search(limit *int, offset *int) ([]models.MemberModel, int, error) {
	var count int64

	// Сначала считаем общее количество всех записей
	if err := database.DB.Model(&models.Member{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Создаем SQL-запрос с LEFT JOIN для получения всех участников и информации о том, являются ли они менторами
	query := `
		SELECT 
			m.id,
			m.username,
			m.first_name,
			m.last_name,
			m.telegram_id,
			m.birthday,
			CASE WHEN mt.id IS NOT NULL THEN true ELSE false END as is_mentor
		FROM members m
		LEFT JOIN mentors mt ON m.id = mt."memberId"
	`

	// Добавляем LIMIT и OFFSET, если они переданы
	if limit != nil {
		query += fmt.Sprintf(" LIMIT %d", *limit)
	}
	if offset != nil {
		query += fmt.Sprintf(" OFFSET %d", *offset)
	}

	// Выполняем запрос
	rows, err := database.DB.Raw(query).Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Обрабатываем результаты
	var result []models.MemberModel
	for rows.Next() {
		var id, tgId int64
		var username, firstName, lastName string
		var birthday time.Time
		var isMentor bool
		if err := rows.Scan(&id, &username, &firstName, &lastName, &tgId, &birthday, &isMentor); err != nil {
			return nil, 0, err
		}
		result = append(result, models.MemberModel{
			Id:         id,
			Username:   username,
			TelegramID: tgId,
			FirstName:  firstName,
			LastName:   lastName,
			IsMentor:   isMentor,
			Birthday:   &birthday,
		})
	}

	return result, int(count), nil
}

// GetById получает участника по ID с проверкой на статус ментора
func (r *MemberRepository) GetById(id int64) (*models.MemberModel, error) {
	var member models.Member
	if err := database.DB.First(&member, id).Error; err != nil {
		return nil, err
	}

	// Проверяем, является ли участник ментором
	var count int64
	if err := database.DB.Model(&models.MentorDbShortModel{}).Where("memberId = ?", id).Count(&count).Error; err != nil {
		return nil, err
	}

	result := &models.MemberModel{
		Id:        member.Id,
		Username:  member.Username,
		FirstName: member.FirstName,
		LastName:  member.LastName,
		IsMentor:  count > 0, // Если есть хотя бы одна запись, то участник является ментором
		Birthday:  member.Birthday,
	}

	return result, nil
}

// UpdateBirthday обновляет дату рождения участника
func (r *MemberRepository) UpdateBirthday(memberID int64, birthday time.Time) error {
	result := database.DB.Model(&models.Member{}).
		Where("id = ?", memberID).
		Update("birthday", birthday)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("member not found")
	}

	return nil
}

func (r *MemberRepository) GetTodayBirthdays() ([]string, error) {
	query := `
		SELECT 
			username
		FROM members
		WHERE
    		DATE_PART('day', birthday) = date_part('day', CURRENT_DATE)
			AND
    		DATE_PART('month', birthday) = date_part('month', CURRENT_DATE)
	`

	rows, err := database.DB.Raw(query).Rows()
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
