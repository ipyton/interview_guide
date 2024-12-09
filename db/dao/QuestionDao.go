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

type QuestionInterfaceImpl struct {
	QuestionInterface
}

var searchDao SearchDaoImpl

func (impl *QuestionInterfaceImpl) BatchAdd(questions *[]model.QuestionModel) error {
	var collection = db.MongoClient.Database("interview_guide").Collection("question")
	var documents []interface{}
	for _, question := range *questions {
		documents = append(documents, question)
	}
	_, err := collection.InsertMany(context.TODO(), documents)
	return err
}

func (impl *QuestionInterfaceImpl) UpsertQuestion(question *model.QuestionModel) error {
	var collection = db.MongoClient.Database("interview_guide").Collection("question")

	//filter := bson.M{"_id": 1}
	//update := bson.M{"$set": }
	questionId := question.ID

	if questionId == -1 {
		value, err2 := getNextSequenceValue("questions")
		if err2 != nil {
			return err2
		}

		questionId = value
	}
	question.ID = questionId
	filter := bson.M{"question_id": questionId}
	update := bson.D{
		{"$set", bson.D{
			{"title", question.Title},
			{"content", question.Content},
			{"details", question.Details},
			{"author_id", question.AuthorID},
			{"author_name", question.AuthorName},
			{"avatar", question.Avatar},
			{"likes", question.Likes},
			{"views", question.Views},
			{"tags", question.Tags},
			{"class_id", question.ClassId},
		}},
	}

	// 设置 upsert 选项
	opts := options.Update().SetUpsert(true)
	err2 := searchDao.UpsertQuestionIndex(*question)
	if err2 != nil {
		log.Fatal(err2.Error())
		return err2
	}
	_, err := collection.UpdateOne(context.TODO(), filter, update, opts)
	return err
}

func (impl *QuestionInterfaceImpl) GetQuestionById(id int64) (model.QuestionModel, error) {
	var collection = db.MongoClient.Database("interview_guide").Collection("question")

	filter := bson.M{"question_id": id}
	var question model.QuestionModel
	return question, collection.FindOne(context.TODO(), filter).Decode(&question)
}

func (impl *QuestionInterfaceImpl) GetQuestionsByPaging(lastId int64, classId int64) (*[]model.QuestionModel, error) {
	// 使用 question_id 字段作为分页游标
	limit := 10
	var collection = db.MongoClient.Database("interview_guide").Collection("question")

	filter := bson.M{
		"question_id": bson.M{"$gt": lastId}, // 查询 question_id 大于上一页的最后一个 question_id
		"class_id":    bson.M{"$eq": classId},
	}
	projection := bson.M{
		"question_id": 1, // 1表示包含字段
		"title":       1, // 1表示包含字段
		"content":     1, // 0表示排除字段（如果不需要_id字段）
	}
	findOptions := options.Find().SetProjection(projection)
	findOptions.SetSort(bson.D{{"question_id", 1}}) // 按 question_id 升序排序

	findOptions.SetLimit(int64(limit))

	// 执行查询
	cursor, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// 解析查询结果
	var questions []model.QuestionModel
	for cursor.Next(context.TODO()) {
		var question model.QuestionModel
		if err := cursor.Decode(&question); err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}

	if err := cursor.Err(); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	fmt.Println(len(questions))

	return &questions, nil
}

func (impl *QuestionInterfaceImpl) QueryQuestions(page int64) ([]model.QuestionModel, error) {
	var collection = db.MongoClient.Database("interview_guide").Collection("question")

	if page == 0 {
		page = 1
	}
	pageSize := 10
	skip := (page - 1) * int64(pageSize)

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
func (impl *QuestionInterfaceImpl) DelQuestion(id int64) error {
	var collection = db.MongoClient.Database("interview_guide").Collection("question")

	_, err := collection.DeleteOne(context.TODO(), bson.M{"id": id})
	return err
}

func (impl *QuestionInterfaceImpl) AdviceQuestion(question model.AdvisedQuestions) error {
	var collection = db.MongoClient.Database("interview_guide").Collection("question_review")
	_, err := collection.InsertOne(context.Background(), question)
	return err
}

func (impl *QuestionInterfaceImpl) GetAdvisedQuestions() ([]model.AdvisedQuestions, error) {
	var collection = db.MongoClient.Database("interview_guide").Collection("question_review")
	filter := bson.M{}
	cur, err := collection.Find(context.Background(), filter)
	var results []model.AdvisedQuestions

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

func (impl *QuestionInterfaceImpl) ApproveAQuestion(questionId int64) error {
	var collection = db.MongoClient.Database("interview_guide").Collection("question_review")
	//var questions_collection = db.MongoClient.Database("interview_guide").Collection("question")
	filter := bson.M{"question_id": questionId}
	update := bson.D{
		{"$set", bson.D{{"review_status", "done"}}},
	}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	questionResult := collection.FindOneAndDelete(context.TODO(), bson.M{"question_id": questionId})
	var result model.AdvisedQuestions
	err = questionResult.Decode(&result)
	if err != nil {
		return err
	}
	err = impl.UpsertQuestion(&result.QuestionModel)
	if err != nil {
		return err
	}
	return nil
}
