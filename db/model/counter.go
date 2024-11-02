package model

import "time"

type CounterModel struct {
	Id        int       `gorm:"column:id" json:"id" bson:"_id"`
	Count     int       `gorm:"column:count" json:"count" bson:"count"`
	CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt" bson:"updatedAt"`
}
