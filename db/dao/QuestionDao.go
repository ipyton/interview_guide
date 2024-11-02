package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

type QuestionInterfaceImpl struct {
}

func (impl *QuestionInterfaceImpl) AddQuestion(question *model.QuestionModel) error {
	var collection = db.MongoClient.Database("interview_guide").Collection("question")

	_, err := collection.InsertOne(context.TODO(), question)
	return err
}

func (impl *QuestionInterfaceImpl) UpdateQuestion(question *model.QuestionModel) error {
	var collection = db.MongoClient.Database("interview_guide").Collection("question")

	filter := bson.M{"_id": question.ID}
	update := bson.M{"$set": question}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (impl *QuestionInterfaceImpl) GetQuestionById(id int) (model.QuestionModel, error) {
	var collection = db.MongoClient.Database("interview_guide").Collection("question")

	filter := bson.M{"_id": id}
	var question model.QuestionModel
	return question, collection.FindOne(context.TODO(), filter).Decode(&question)
}
func (impl *QuestionInterfaceImpl) QueryQuestions(page int) ([]model.QuestionModel, error) {
	var collection = db.MongoClient.Database("interview_guide").Collection("question")

	if page == 0 {
		page = 1
	}
	pageSize := 10                // 每页显示的条数
	skip := (page - 1) * pageSize // 计算跳过的记录数

	// 查询数据
	cur, err := collection.Find(context.TODO(), bson.M{}, options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize)))
	if err != nil {
		log.Fatal(err)
	}
	var results []model.QuestionModel
	for cur.Next(context.TODO()) {
		var result model.QuestionModel
		if err := cur.Decode(&result); err != nil {
			log.Fatal(err)
		}
		results = append(results, result)
	}
	return results, cur.Close(context.TODO())
}

// DelQuestion deletes a question by ID
func (impl *QuestionInterfaceImpl) DelQuestion(id int) error {
	var collection = db.MongoClient.Database("interview_guide").Collection("question")

	_, err := collection.DeleteOne(context.TODO(), bson.M{"id": id})
	return err
}
