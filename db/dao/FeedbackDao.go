package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

type FeedbackDao struct {
}

var counter = CounterImpl{}

func (FeedbackDao) GetFeedback() ([]model.Feedback, error) {
	var collection = db.MongoClient.Database("interview_guide").Collection("feed")
	filter := bson.M{}
	cur, err := collection.Find(context.Background(), filter)
	var results []model.Feedback

	if err != nil {
		return results, err
	}
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		var result model.Feedback
		if err := cur.Decode(&result); err != nil {
			log.Fatal(err)
			return results, err
		}
		results = append(results, result)
	}
	return results, cur.Err()
}

func (FeedbackDao) SaveFeedback(feedback model.Feedback) error {
	var collection = db.MongoClient.Database("interview_guide").Collection("feed")
	increase, _ := counter.GetAndIncrease("feed")
	feedback.Id = increase
	//var questions_collection = db.MongoClient.Database("interview_guide").Collection("question")
	_, err := collection.InsertOne(context.Background(), feedback)
	return err
}

func (FeedbackDao) DeleteFeedback(id int64) {
	var collection = db.MongoClient.Database("interview_guide").Collection("feed")
	filter := bson.M{
		"id": id,
	}
	collection.DeleteOne(context.TODO(), filter)
}
