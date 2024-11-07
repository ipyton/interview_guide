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
	GetCollections() (*[]model.BookmarkCollectionModel, error)
	GetItemsInCollection(userId string, questionId int) ([]*model.BookmarkItemModel, error)
	DeleteBookMarkItem(userId string, collectionId int, questionId int) error
	AddBookMarkItem(item *model.BookmarkItemModel) error
	AddBookMarkCollection(collection *model.BookmarkCollectionModel) error
	DeleteBookMarkCollection(userId string, collectionID int) error
	IsResourceCollected(userId string, questionId int) (bool, error)
}

type QuestionInterface interface {
	UpsertQuestion(question *model.QuestionModel) error // Add a new question
	DelQuestion(id int) error                           // Delete a question by ID
	GetQuestionById(id int) (model.QuestionModel, error)
	QueryQuestions(page int) ([]model.QuestionModel, error)
}

type ClassInterface interface {
	UpsertClass(classIds []string, class model.ClassModel) error
	GetRootClasses() ([]model.ClassModel, error)
	GetSubClasses(classIds []string, class model.ClassModel) (model.ClassModel, error)
	DeleteClass(classIds []string) error
}

type CounterInterface interface {
	GetAndIncrease(increaseDoc string) (int, error)
}
