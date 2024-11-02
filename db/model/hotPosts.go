package model

import "time"

type HotPost struct {
	InfoID      int       `json:"info_id" bson:"info_id"`
	QuestionID  int       `json:"question_id" bson:"question_id"`
	Intro       string    `json:"intro" bson:"intro"`
	AuthorID    int       `json:"author_id" bson:"author_id"`
	Title       string    `json:"title" bson:"title"`
	Likes       int       `json:"likes" bson:"likes"`
	PublishDate time.Time `json:"publish_date" bson:"publish_date"`
	Extra1      string    `json:"extra1,omitempty" bson:"extra1,omitempty"` // 额外字段1
	Extra2      string    `json:"extra2,omitempty" bson:"extra2,omitempty"` // 额外字段2
	Extra3      string    `json:"extra3,omitempty" bson:"extra3,omitempty"` // 额外字段3
}
