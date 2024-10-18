package model

import "time"

type HotPost struct {
	InfoID      int       `json:"info_id"`
	QuestionID  int       `json:"question_id"`
	Intro       string    `json:"intro"`
	AuthorID    int       `json:"author_id"`
	Title       string    `json:"title"`
	Likes       int       `json:"likes"`
	PublishDate time.Time `json:"publish_date"`
	Extra1      string    `json:"extra1,omitempty"` // 额外字段1
	Extra2      string    `json:"extra2,omitempty"` // 额外字段2
	Extra3      string    `json:"extra3,omitempty"` // 额外字段3
}
