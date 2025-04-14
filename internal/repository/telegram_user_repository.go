package repository

import (
	"ithozyeva/internal/models"

	"gorm.io/gorm"
)

type TelegramUserRepository struct {
	db *gorm.DB
}

func NewTelegramUserRepository(db *gorm.DB) *TelegramUserRepository {
	return &TelegramUserRepository{db: db}
}

func (r *TelegramUserRepository) Create(user *models.TelegramUser) (*models.TelegramUser, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *TelegramUserRepository) GetByTelegramID(telegramID int64) (*models.TelegramUser, error) {
	var user models.TelegramUser
	err := r.db.Where("telegram_id = ?", telegramID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *TelegramUserRepository) Update(user *models.TelegramUser) (*models.TelegramUser, error) {
	if err := r.db.Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *TelegramUserRepository) Delete(user *models.TelegramUser) error {
	return r.db.Delete(user).Error
}

func (r *TelegramUserRepository) GetById(id int64) (*models.TelegramUser, error) {
	var user models.TelegramUser
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *TelegramUserRepository) Search(limit *int, offset *int) ([]models.TelegramUser, int64, error) {
	var users []models.TelegramUser
	var count int64

	if err := r.db.Model(&models.TelegramUser{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	query := r.db.Model(&models.TelegramUser{})
	if limit != nil {
		query = query.Limit(*limit)
	}
	if offset != nil {
		query = query.Offset(*offset)
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, count, nil
}

func (r *TelegramUserRepository) GetByToken(token string) (*models.TelegramUser, error) {
	var user models.TelegramUser
	err := r.db.Where("token = ?", token).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
