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
	UpsertClass(class model.ClassModel) error
	GetClasses(parentClassId int) ([]model.ClassModel, error)
	DeleteClass(classId int) error
}

type CounterInterface interface {
	GetAndIncrease(increaseDoc string) (int, error)
}

type UserInterface interface {
	SaveLoginStatus(user model.User) error
	DeleteLoginStatus(userCode string) error
	UpdateUserInfo(user model.User) error
	CancelRegistration(user model.User) error
	IsUserExists(openid string) bool
	ChangeMembershipStatus(openid string, status bool) error
	AddPoints(openid string, points int) error
	Registration(userStatus model.UserStatus, user model.User) error
	UpsertLoginStatus(userStatus model.UserStatus, ip string) error
}
