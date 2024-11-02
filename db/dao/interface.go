package dao

import (
	"wxcloudrun-golang/db/model"
)

//type CounterInterface interface {
//	GetCounter(id int) (*model.CounterModel, error)
//	UpsertCounter(counter *model.CounterModel) error
//	ClearCounter(id int) error
//}
//
//type CounterInterfaceImp struct{}
//
//var Imp CounterInterface = &CounterInterfaceImp{}

type CollectionInterface interface {
	GetBookCollections() (*[]model.BookmarkCollectionModel, error)
	GetBookMarkItems(collectionID int) ([]*model.BookmarkItemModel, error)
	DeleteBookMarkItem(resourceID string) error
	AddBookMarkItem(item *model.BookmarkItemModel) error
	AddBookMarkCollection(collection *model.BookmarkCollectionModel) error
	DeleteBookMarkCollection(collectionID int) error
}

type QuestionInterface interface {
	AddQuestion(question *model.QuestionModel) error // Add a new question
	DelQuestion(id int) error                        // Delete a question by ID
	GetQuestionById(id int) (model.QuestionModel, error)
	QueryQuestions(page int) ([]model.QuestionModel, error)
}
