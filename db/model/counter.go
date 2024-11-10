package model

type CounterModel struct {
	Id    string `gorm:"column:id" json:"id" bson:"_id"`
	Value int64  `gorm:"column:value" json:"value" bson:"value"`
}
