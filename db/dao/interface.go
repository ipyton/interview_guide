package dao

import (
	"wxcloudrun-golang/db/model"
)

// CounterInterface 计数器数据模型接口
type CounterInterface interface {
	GetCounter(id int32) (*model.CounterModel, error)
	UpsertCounter(counter *model.CounterModel) error
	ClearCounter(id int32) error
}

// CounterInterfaceImp 计数器数据模型实现
type CounterInterfaceImp struct{}

// Imp 实现实例
var Imp CounterInterface = &CounterInterfaceImp{}

type CollectionInterface interface {
	GetBookCollections() ([]*model.BookmarkCollectionModel, error)
	GetBookMarkItems(collectionID int) ([]*model.BookmarkItemModel, error)
	DeleteBookMarkItem(resourceID string) error
	AddBookMarkItem(item *model.BookmarkItemModel) error
	AddBookMarkCollection(collection *model.BookmarkCollectionModel) error
	DeleteBookMarkCollection(collectionID int) error
}

type QuestionInterface interface {
	AddQuestion(question *model.QuestionModel) error // Add a new question
	DelQuestion(id int32) error                      // Delete a question by ID
}
