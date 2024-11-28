package model

import "time"

type BookmarkCollectionModel struct {
	OpenId         string    `gorm:"column:openid" json:"openid" bson:"openid"`
	CollectionID   int64     `gorm:"column:collection_id" json:"collection_id" bson:"collection_id"`
	CollectionName string    `gorm:"column:collection_name" json:"collection_name" bson:"collection_name"`
	Description    string    `gorm:"column:description" json:"description" bson:"description"`
	Count          int64     `gorm:"column:count" json:"count" bson:"count"`
	CreateAt       time.Time `gorm:"column:created_at" json:"created_at" bson:"created_at"`
	Extra1         string    `gorm:"column:extra1;omitempty" json:"extra1,omitempty" bson:"extra1,omitempty"` // 额外字段1
	Extra2         string    `gorm:"column:extra2;omitempty" json:"extra2,omitempty" bson:"extra2,omitempty"` // 额外字段2
	Extra3         string    `gorm:"column:extra3;omitempty" json:"extra3,omitempty" bson:"extra3,omitempty"` // 额外字段3
}

// compare with question model this model is more simple with less information
type BookmarkQuestionModel struct {
	CollectionID int64     `gorm:"column:collection_id" json:"collection_id" bson:"collection_id"`
	QuestionId   int64     `gorm:"column:question_id" json:"question_id" bson:"question_id"`
	Title        string    `gorm:"column:title" json:"title" bson:"title"` //question, post 两周类型
	Content      string    `gorm:"column:content" json:"content" bson:"content"`
	CreateAt     time.Time `gorm:"column:created_at" json:"created_at" bson:"created_at"`
	Extra1       string    `gorm:"column:extra1;omitempty" json:"extra1,omitempty" bson:"extra1,omitempty"` // 额外字段1
	Extra2       string    `gorm:"column:extra2;omitempty" json:"extra2,omitempty" bson:"extra2,omitempty"` // 额外字段2
	Extra3       string    `gorm:"column:extra3;omitempty" json:"extra3,omitempty" bson:"extra3,omitempty"` // 额外字段3
}
type BookmarkResourceModel struct {
	CollectionID  int64  `gorm:"column:collection_id" json:"collection_id" bson:"collection_id"`
	ResourceID    int64  `gorm:"column:resource_id" json:"resource_id" bson:"resource_id"`
	Type          string `gorm:"column:type" json:"type" bson:"type"` //question, post 两周类型
	ResourceTitle string `gorm:"column:resource_title" json:"resource_title" bson:"resource_title"`
	Content       string `gorm:"column:content" json:"content" bson:"content"`
	Extra1        string `gorm:"column:extra1;omitempty" json:"extra1,omitempty" bson:"extra1,omitempty"` // 额外字段1
	Extra2        string `gorm:"column:extra2;omitempty" json:"extra2,omitempty" bson:"extra2,omitempty"` // 额外字段2
	Extra3        string `gorm:"column:extra3;omitempty" json:"extra3,omitempty" bson:"extra3,omitempty"` // 额外字段3
}
