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

func (FeedbackDao) GetFeedback() ([]model.Feedback, error) {
	var collection = db.MongoClient.Database("interview_guide").Collection("feedback")
	filter := bson.M{}
	cur, err := collection.Find(context.Background(), filter)
	var results []model.Feedback

	if err != nil {
		return results, err
	}
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		var result model.AdvisedQuestions
		if err := cur.Decode(&result); err != nil {
			log.Fatal(err)
			return results, err
		}
	}
	return results, cur.Err()
}

func (FeedbackDao) SaveFeedback(feedback model.Feedback) error {
	var collection = db.MongoClient.Database("interview_guide").Collection("feed")
	//var questions_collection = db.MongoClient.Database("interview_guide").Collection("question")
	_, err := collection.InsertOne(context.Background(), feedback)
	return err
}
