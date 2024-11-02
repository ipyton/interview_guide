package model

type QuestionModel struct {
	ID       int    `gorm:"column:question_id;primaryKey" json:"question_id" bson:"question_id"` // Unique identifier for the question
	Title    string `gorm:"column:title;size:50" json:"title" bson:"title"`                      // Title of the question
	Content  string `gorm:"column:content;type:mediumtext" json:"content" bson:"content"`        // Content of the question
	AuthorID int    `gorm:"column:author_id" json:"author_id" bson:"author_id"`                  // ID of the author
	Avatar   string `gorm:"column:avatar;size:100" json:"avatar" bson:"avatar"`                  // Avatar of the author
	Likes    int    `gorm:"column:likes" json:"likes" bson:"likes"`                              // Number of likes
	Views    int    `gorm:"column:views" json:"views" bson:"views"`                              // Number of views
	Tag1     string `gorm:"column:tag1;size:20" json:"tag1" bson:"tag1"`                         // Tag 1
	Tag2     string `gorm:"column:tag2;size:20" json:"tag2" bson:"tag2"`                         // Tag 2
	Tag3     string `gorm:"column:tag3;size:20" json:"tag3" bson:"tag3"`                         // Tag 3
	Tag4     string `gorm:"column:tag4;size:20" json:"tag4" bson:"tag4"`                         // Tag 4
	Tag5     string `gorm:"column:tag5;size:20" json:"tag5" bson:"tag5"`                         // Tag 5
	Extra1   string `gorm:"column:extra1;size:100" json:"extra1" bson:"extra1"`                  // Extra field 1
	Extra2   string `gorm:"column:extra2;size:100" json:"extra2" bson:"extra2"`                  // Extra field 2
}
