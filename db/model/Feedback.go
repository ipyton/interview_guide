package model

type Feedback struct {
	Id      int64  `json:"id" bson:"id"`
	Title   string `json:"title" bson:"title"`
	Content string `json:"content" bson:"content"`
	Author  string `json:"author" bson:"author"`
}
