package model

type Feedback struct {
	Id      int64  `json:"id" bson:"id"`
	Content string `json:"content" bson:"content"`
	Author  string `json:"author" bson:"author"`
	Contact string `json:"contact" bson:"contact"`
}
