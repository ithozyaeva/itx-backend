package models

import "time"

type PlaceType string

const (
	EventOnline  PlaceType = "ONLINE"
	EventOffline PlaceType = "OFFLINE"
	EventHybrid  PlaceType = "HYBRID"
)

type Event struct {
	Id              int64     `json:"id" gorm:"primaryKey"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Date            time.Time `json:"date" time_format:"2006-01-02T15:04" time_location:"UTC"`
	PlaceType       PlaceType `json:"placeType"`
	Place           string    `json:"place"`
	CustomPlaceType string    `json:"customPlaceType"`
	EventType       string    `json:"eventType"`
	Open            bool      `json:"open"`
	VideoLink       string    `json:"videoLink" gorm:"column:video_link"`
	EventTags       []EventTag `json:"eventTags" gorm:"many2many:event_event_tags;foreignKey:id;joinForeignKey:event_id;References:id;joinReferences:event_tag_id;replace:true"`
	Hosts           []Member  `json:"hosts" gorm:"many2many:event_hosts;foreignKey:id;joinForeignKey:event_id;References:id;joinReferences:member_id;replace:true"`
	Members         []Member  `json:"members" gorm:"many2many:event_members;foreignKey:id;joinForeignKey:event_id;References:id;joinReferences:member_id;replace:true"`
}

type EventTag struct {
	Id   int64  `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"unique"`
}