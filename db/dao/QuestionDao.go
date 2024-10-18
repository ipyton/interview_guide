package dao

import (
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

type QuestionInterfaceImpl struct {
}

func (impl *QuestionInterfaceImpl) AddQuestion(question *model.QuestionModel) error {
	cli := db.Get()
	return cli.Table("questions").Create(question).Error // Replace with your actual table name
}

// DelQuestion deletes a question by ID
func (impl *QuestionInterfaceImpl) DelQuestion(id int32) error {
	cli := db.Get()
	return cli.Table("questions").Delete(&model.QuestionModel{ID: id}).Error // Replace with your actual table name
}
