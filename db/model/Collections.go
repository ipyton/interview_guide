package model

import "time"

type BookmarkCollectionModel struct {
	UserId         string    `gorm:"column:user_id" json:"user_id" bson:"user_id"`
	CollectionID   int       `gorm:"column:collection_id" json:"collection_id" bson:"collection_id"`
	CollectionName string    `gorm:"column:collection_name" json:"collection_name" bson:"collection_name"`
	Description    string    `gorm:"column:description" json:"description" bson:"description"`
	Count          int       `gorm:"column:count" json:"count" bson:"count"`
	CreateAt       time.Time `gorm:"column:create_at" json:"create_at" bson:"create_at"`
	UpdateAt       time.Time `gorm:"column:update_at" json:"update_at" bson:"update_at"`
	Extra1         string    `gorm:"column:extra1;omitempty" json:"extra1,omitempty" bson:"extra1,omitempty"` // 额外字段1
	Extra2         string    `gorm:"column:extra2;omitempty" json:"extra2,omitempty" bson:"extra2,omitempty"` // 额外字段2
	Extra3         string    `gorm:"column:extra3;omitempty" json:"extra3,omitempty" bson:"extra3,omitempty"` // 额外字段3
}

type BookmarkQuestionModel struct {
	CollectionID int       `gorm:"column:collection_id" json:"collection_id" bson:"collection_id"`
	QuestionId   int       `gorm:"column:resource_id" json:"resource_id" bson:"resource_id"`
	Title        string    `gorm:"column:type" json:"type" bson:"type"` //question, post 两周类型
	Introduction string    `gorm:"column:resource_title" json:"resource_title" bson:"resource_title"`
	AddTime      time.Time `gorm:"column:add_time" json:"add_time" bson:"add_time"`
	Extra1       string    `gorm:"column:extra1;omitempty" json:"extra1,omitempty" bson:"extra1,omitempty"` // 额外字段1
	Extra2       string    `gorm:"column:extra2;omitempty" json:"extra2,omitempty" bson:"extra2,omitempty"` // 额外字段2
	Extra3       string    `gorm:"column:extra3;omitempty" json:"extra3,omitempty" bson:"extra3,omitempty"` // 额外字段3
}
type BookmarkResourceModel struct {
	CollectionID  int    `gorm:"column:collection_id" json:"collection_id" bson:"collection_id"`
	ResourceID    int    `gorm:"column:resource_id" json:"resource_id" bson:"resource_id"`
	Type          string `gorm:"column:type" json:"type" bson:"type"` //question, post 两周类型
	ResourceTitle string `gorm:"column:resource_title" json:"resource_title" bson:"resource_title"`
	Content       string `gorm:"column:content" json:"content" bson:"content"`
	Extra1        string `gorm:"column:extra1;omitempty" json:"extra1,omitempty" bson:"extra1,omitempty"` // 额外字段1
	Extra2        string `gorm:"column:extra2;omitempty" json:"extra2,omitempty" bson:"extra2,omitempty"` // 额外字段2
	Extra3        string `gorm:"column:extra3;omitempty" json:"extra3,omitempty" bson:"extra3,omitempty"` // 额外字段3
}
