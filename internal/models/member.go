package models

type Member struct {
	Id         int64  `json:"id"`
	Username   string `json:"tg"`
	TelegramID int64  `json:"telegramID"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Role       string `json:"role"`
}

type ReviewOnCommunity struct {
	Id       int    `json:"id"`
	AuthorId uint   `json:"authorId" gorm:"column:authorId"`
	Author   Member `json:"author" gorm:"foreignKey:authorId"`
	Text     string `json:"text"`
	Date     string `json:"date"`
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
	Text     string `json:"text" binding:"required"`
	Date     string `json: "date"`
	AuthorTg string `json:"authorTg" binding:"required"`
}

// MemberModel представляет модель участника с дополнительной информацией
type MemberModel struct {
	Id         int64  `json:"id"`
	TelegramID int64  `json:"telegramID"`
	Username   string `json:"tg"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	IsMentor   bool   `json:"isMentor"`
}

func (ReviewOnCommunity) TableName() string {
	return "reviewOnCommunity"
}
