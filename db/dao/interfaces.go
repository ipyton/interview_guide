package dao

import (
	"mime/multipart"
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

type CollectionQuestionInterface interface {
	GetCollections(openId string) (*[]model.BookmarkCollectionModel, error)
	DeleteBookMarkQuestion(openid string, collectionId int64, questionId int64) error
	AddBookMarkQuestion(openId string, collectionID int64, questionId int64) error
	AddQuestionCollection(collection *model.BookmarkCollectionModel) error
	DeleteBookMarkCollection(openid string, collectionID int64) error
	IsResourceCollected(openid string, questionId int64) (bool, error)
	GetCollectionItemsByCollectionAndTime(openId string, lastQuestionId int64, isDescending bool, collectionId int64) (*[]model.BookmarkQuestionModel, error)
	// IsExist(openid string, questionId int64, )
	GetCollectionItemsByTag(openId string, questionId int64, tag int64, collectionId int64) (*[]model.BookmarkQuestionModel, error)
}

type QuestionInterface interface {
	UpsertQuestion(question *model.QuestionModel) error // Add a new question
	DelQuestion(id int64) error                         // Delete a question by ID
	GetQuestionById(id int64) (model.QuestionModel, error)
	QueryQuestions(page int64) ([]model.QuestionModel, error)
	BatchAdd(questions *[]model.QuestionModel) error
	GetQuestionsByPaging(lastId int64, classId int64) (*[]model.QuestionModel, error)
	AdviceQuestion(question model.AdvisedQuestions) error
	GetAdvisedQuestions() ([]model.AdvisedQuestions, error)
	ApproveAQuestion(questionId int64) error
	RejectQuestion(questionId int64) error
	RateQuestion(userRate model.UserRate) error
	GetRatings(questionId int64) (model.RatingsResponse, error)
	SeeBefore(seeBefore model.SeeBeforeCount) error
}

type ClassInterface interface {
	UpdateClass(class model.ClassModel) error
	InsertClass(class model.ClassModel) error
	GetClasses(parentClassId int64) ([]model.ClassModel, error)
	DeleteClass(classId int64) error
	GetQuestionsById(lastId int64) (*[]model.QuestionModel, error)
}

type CounterInterface interface {
	GetAndIncrease(increaseDoc string) (int64, error)
}

type UserStatusInterface interface {
	SaveLoginStatus(user model.User) error
	DeleteLoginStatus(userCode string) error
	CancelRegistration(user model.User) error
	IsUserExists(openid string) bool
	Registration(userStatus model.UserStatus, user model.User) error
	UpsertLoginStatus(userStatus model.UserStatus, ip string) error
}
type UserInformationInterface interface {
	UploadAvatar(openid string, file multipart.File) error
	UpdateUserInfo(user model.User) error
	ChangeMembershipStatus(openid string, status bool) error
	AddPoints(openid string, points int64) error
	UpdateUserClass(openid string, classId int64, className string) error
}
type VoiceInterface interface {
	GetVoice(questionId string) error
	GenerateVoice(questionId string) error
}
type FileManager interface {
	GetFile(questionId string) (string, error)
	UploadFile(fileName string, fileType string, file []byte) (string, error)
	DeleteFile(fileName string, fileType string) error
	UploadFileByMultipartFile(fileName string, fileType string, file multipart.File) error
}

type SearchDaoInterface interface {
	UpsertQuestionIndex(question model.QuestionModel) error
	SearchQuestions(keyword string, pageSize int64, page int64) (model.QuestionModel, error)
	GetSuggestions(keyword string) (SuggestResponse, error)
}
