package model

type ClassModel struct {
	ClassId    int          `bson:"class_id" json:"class_id"`
	ClassName  string       `bson:"class_name" json:"class_name"`
	Count      int          `bson:"count" json:"count"`
	SubClasses []ClassModel `bson:"sub_class" json:"sub_class"`
	Extra1     string       `bson:"extra1" json:"extra1"`
	Extra2     string       `bson:"extra2" json:"extra2"`
	Extra3     string       `bson:"extra3" json:"extra3"`
}
