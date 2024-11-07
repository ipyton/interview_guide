package model

type QuestionModel struct {
	ID         int      `gorm:"column:question_id;primaryKey" json:"question_id" bson:"question_id"` // Unique identifier for the question
	Title      string   `gorm:"column:title;size:50" json:"title" bson:"title"`                      // Title of the question
	ClassId    int      `gorm:"column:class_id" json:"class_id" bson:"class_id"`
	Type       string   `gorm:"column:type" json:"type" bson:"type"`                          // question and post
	Content    string   `gorm:"column:content;type:mediumtext" json:"content" bson:"content"` // Content of the question
	Details    string   `gorm:"column:details;type:mediumtext" json:"details" bson:"details"`
	AuthorID   int      `gorm:"column:author_id" json:"author_id" bson:"author_id"` // ID of the author
	AuthorName string   `gorm:"column:author_name" json:"author_name" bson:"author_name"`
	Avatar     string   `gorm:"column:avatar;size:100" json:"avatar" bson:"avatar"` // Avatar of the author
	Likes      int      `gorm:"column:likes" json:"likes" bson:"likes"`             // Number of likes
	Views      int      `gorm:"column:views" json:"views" bson:"views"`             // Number of views
	Tags       []string `gorm:"column:tags" json:"tags" bson:"tags"`                // Tag 1

}
