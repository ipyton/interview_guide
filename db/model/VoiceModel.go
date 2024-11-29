package model

type VoiceModel struct {
	QuestionId int64  `bson:"QuestionId"`
	Path       string `bson:"Path"`
}
