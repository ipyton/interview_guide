package model

type ClassModel struct {
	ClassId       int    `bson:"class_id" json:"class_id"`
	ClassName     string `bson:"class_name" json:"class_name"`
	Count         int    `bson:"count" json:"count"`
	ParentClassId int    `bson:"parent_class_id" json:"parent_class_id"`
	Extra1        string `bson:"extra1" json:"extra1"`
	Extra2        string `bson:"extra2" json:"extra2"`
	Extra3        string `bson:"extra3" json:"extra3"`
}
