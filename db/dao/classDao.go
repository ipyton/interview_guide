package dao

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

type ClassInterfaceImpl struct {
	ClassInterface
}

func (classInterfaceImpl ClassInterfaceImpl) InsertClass(class model.ClassModel) error {
	session, err := db.MongoClient.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(context.Background())

	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {
		var collection = db.MongoClient.Database("interview_guide").Collection("classes")
		var err error
		fmt.Println(class)
		class.ClassId, err = CounterImpl{}.GetAndIncrease("classes")
		filter := bson.D{{"class_id", class.ParentClassId}}
		_, err = collection.UpdateOne(context.TODO(), filter, bson.D{{"$set", bson.D{
			{"isLeaf", false},
		}}})
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		class.IsLeaf = true
		_, err = collection.InsertOne(context.TODO(), class)

		return nil, err
	}
	_, err = session.WithTransaction(context.Background(), transactionFunc)

	return err
}

func (classInterfaceImpl ClassInterfaceImpl) UpdateClass(class model.ClassModel) error {

	var collection = db.MongoClient.Database("interview_guide").Collection("classes")
	var err error
	//新插入
	update := bson.D{
		{"$set", bson.D{
			{"class_name", class.ClassName},
			//{"count", class.Count},
			//{"parent_class_id", class.ParentClassId},
			{"extra1", class.Extra1},
		}},
	}
	opts := options.Update().SetUpsert(true)
	_, err = collection.UpdateOne(context.TODO(), bson.M{"class_id": class.ClassId}, update, opts)

	return err

}

func (classInterfaceImpl ClassInterfaceImpl) GetClasses(parentClassId int64) ([]model.ClassModel, error) {
	var collection = db.MongoClient.Database("interview_guide").Collection("classes")
	var err error
	var result []model.ClassModel
	cur, err := collection.Find(context.TODO(), bson.M{"parent_class_id": parentClassId})
	defer cur.Close(context.TODO())
	if err = cur.All(context.TODO(), &result); err != nil {
		log.Fatal(err.Error())
	}
	for _, class := range result {
		fmt.Printf("Name: %s, ClassName: %s\n", class.ClassId, class.ClassName)
	}
	return result, err
}
func (classInterfaceImpl ClassInterfaceImpl) DeleteClass(classId int64) error {
	var collection = db.MongoClient.Database("interview_guide").Collection("classes")
	_, err := collection.DeleteOne(context.TODO(), bson.M{"class_id": classId})
	if err != nil {
		return err
	}
	return nil
}
