package dao

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

type ClassInterfaceImpl struct {
	ClassInterface
}

func (classInterfaceImpl ClassInterfaceImpl) UpsertClass(classIds []string, class model.ClassModel) error {
	//var collection = db.MongoClient.Database("interview_guide").Collection("classes")
	var err error
	//if len(classIds) == 0 {
	//	if class.ClassId == -1 {
	//		class.ClassId, err = CounterImpl{}.GetAndIncrease("classes")
	//		if err != nil {
	//			return err
	//		}
	//	}
	//
	//	filter := bson.M{"class_id": class.ClassId}
	//	update := bson.D{
	//		{"$set", bson.D{
	//			{"name", class.ClassName},
	//			{"count", class.Count},
	//			{"sub_class", class.SubClasses},
	//			{"extra1", class.Extra1},
	//		}},
	//	}
	//
	//	opts := options.Update().SetUpsert(true)
	//	_, err = collection.UpdateOne(context.TODO(), filter, update, opts)
	//	if err != nil {
	//		return err
	//	}
	//}
	return err

}

func (classInterfaceImpl ClassInterfaceImpl) GetSubClasses(classIds []string, class model.ClassModel) ([]model.ClassModel, error) {
	var collection = db.MongoClient.Database("interview_guide").Collection("classes")
	if len(classIds) == 0 {
		projection := bson.M{"class_id": 1, "class_name": 1, "count": 1, "extra1": 1, "extra2": 1, "extra3": 1}
		cur, err := collection.Find(context.TODO(), bson.M{}, options.Find().SetProjection(projection))
		if err != nil {
			log.Fatal(err)
		}
		defer cur.Close(context.TODO())
		var results []model.ClassModel
		if err = cur.All(context.TODO(), &results); err != nil {
			log.Fatal(err)
		}

		for _, class := range results {
			fmt.Printf("Name: %s, Email: %s\n", class.ClassId, class.ClassName)
		}
	}
	filter := bson.M{}
	subClassPath := "sub_class"
	for _, classId := range classIds {
		filter[subClassPath+".class_id"] = classId
		subClassPath += ".$[].sub_class"
	}

	var result model.ClassModel
	filter["class_id"] = classIds[0]
	projection := bson.M{
		"sub_class": 1,
	}
	err := collection.FindOne(context.TODO(), filter, options.FindOne().SetProjection(projection)).Decode(&result)
	return result.SubClasses, err
}
func (classInterfaceImpl ClassInterfaceImpl) DeleteClass(classIds []string) error {
	//var collection = db.MongoClient.Database("interview_guide").Collection("classes")
	//subClassPath := bson.D{{"class_id", classIds[0]}, {"class_name", classIds[1]}}

	return nil
}
