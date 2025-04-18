package models

import (
	"time"
)

const (
	MemberRoleUnsubscriber MemberRole = "UNSUBSCRIBER"
	MemberRoleSubscriber   MemberRole = "SUBSCRIBER"
	MemberRoleMentor       MemberRole = "MENTOR"
	MemberRoleAdmin        MemberRole = "ADMIN"
)

type MemberRole string

type Member struct {
	Id         int64      `json:"id" gorm:"primaryKey"`
	Username   string     `json:"tg" gorm:"column:username"`
	TelegramID int64      `json:"telegramID" gorm:"column:telegram_id"`
	FirstName  string     `json:"firstName" gorm:"column:first_name"`
	LastName   string     `json:"lastName" gorm:"column:last_name"`
	Role       MemberRole `json:"role" gorm:"column:role"`
	Birthday   *time.Time `json:"birthday" gorm:"column:birthday"`
}

type ReviewOnCommunity struct {
	Id       int    `json:"id"`
	AuthorId uint   `json:"authorId" gorm:"column:authorId"`
	Author   Member `json:"author" gorm:"foreignKey:authorId"`
	Text     string `json:"text"`
	Date     string `json:"date"`
	Status   string `json:"status"`
}

type ReviewOnCommunityWithAuthor struct {
	Id         int    `json:"id"`
	AuthorId   int    `json:"authorId"`
	AuthorName string `json:"authorName"`
	AuthorTg   string `json:"authorTg"`
	Text       string `json:"text"`
	Date       string `json:"date"`
}

// ReviewOnCommunityRequest представляет запрос на создание отзыва о сообществе
type ReviewOnCommunityRequest struct {
	Text     string  `json:"text" binding:"required"`
	Date     *string `json:"date"`
	AuthorTg string  `json:"authorTg" binding:"required"`
}

// MemberModel представляет модель участника с дополнительной информацией
type MemberModel struct {
	Id         int64      `json:"id"`
	TelegramID int64      `json:"telegramID"`
	Username   string     `json:"tg"`
	FirstName  string     `json:"firstName"`
	LastName   string     `json:"lastName"`
	IsMentor   bool       `json:"isMentor"`
	Birthday   *time.Time `json:"birthday"`
	Role       MemberRole `json:"role"`
}

func (ReviewOnCommunity) TableName() string {
	return "reviewOnCommunity"
}
