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

func (classInterfaceImpl ClassInterfaceImpl) UpsertClass(class model.ClassModel) error {
	var collection = db.MongoClient.Database("interview_guide").Collection("classes")
	var err error
	if class.ClassId == -1 {
		class.ClassId, err = CounterImpl{}.GetAndIncrease("classes")
		if err != nil {
			return err
		}
	}
	opts := options.Update().SetUpsert(true)
	fmt.Println(class)
	update := bson.D{
		{"$set", bson.D{
			{"class_name", class.ClassName},
			{"count", class.Count},
			{"parent_class_id", class.ParentClassId},
			{"extra1", class.Extra1},
		}},
	}
	_, err = collection.UpdateOne(context.TODO(), bson.M{"class_id": class.ClassId}, update, opts)

	return err

}

func (classInterfaceImpl ClassInterfaceImpl) GetClasses(parentClassId int) ([]model.ClassModel, error) {
	var collection = db.MongoClient.Database("interview_guide").Collection("classes")
	var err error
	var result []model.ClassModel
	cur, err := collection.Find(context.TODO(), bson.M{"parent_class_id": parentClassId})
	defer cur.Close(context.TODO())
	if err = cur.All(context.TODO(), &result); err != nil {
		log.Fatal(err.Error())
	}
	for _, class := range result {
		fmt.Printf("Name: %s, Email: %s\n", class.ClassId, class.ClassName)
	}
	return result, err
}
func (classInterfaceImpl ClassInterfaceImpl) DeleteClass(classId int) error {
	var collection = db.MongoClient.Database("interview_guide").Collection("classes")
	_, err := collection.DeleteOne(context.TODO(), bson.M{"class_id": classId})
	if err != nil {
		return err
	}
	return nil
}
