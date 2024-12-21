package dao

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

type CollectionQuestionInterfaceImpl struct {
	CollectionQuestionInterface
}

//var Imp CounterInterface = &CounterInterfaceImp{}

func (dao *CollectionQuestionInterfaceImpl) IsResourceCollected(userId string, resourceId int64) (bool, error) {
	return true, errors.New("not implement")
}

func (dao *CollectionQuestionInterfaceImpl) GetCollections(openId string) (*[]model.BookmarkCollectionModel, error) {
	collection := db.MongoClient.Database("interview_guide").Collection("collections")

	// 定义查询过滤器
	filter := bson.M{"openid": openId}

	// 执行查询
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to execute find: %v", err)
	}
	defer cursor.Close(context.TODO())
	// 存储查询结果
	var results []model.BookmarkCollectionModel
	for cursor.Next(context.TODO()) {
		var result model.BookmarkCollectionModel
		if err := cursor.Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode result: %v", err)
		}
		results = append(results, result)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over cursor: %v", err)
	}

	return &results, nil
}

func (dao *CollectionQuestionInterfaceImpl) GetItemsInCollection(openId string, collectionID int64, isDescending bool) (*[]model.BookmarkQuestionModel, error) {
	var items = &[]model.BookmarkQuestionModel{}

	collection := db.MongoClient.Database("interview_guide").Collection("collection_items")
	fmt.Println(collectionID)
	fmt.Println(openId)
	// 构建查询条件
	filter := bson.M{
		"collection_id": collectionID,
		"openid":        openId,
	}
	// 执行查询
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to execute find: %v", err)
	}

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var item model.BookmarkQuestionModel
		if err := cursor.Decode(&item); err != nil {
			return nil, fmt.Errorf("failed to decode result: %v", err)
		}
		*items = append(*items, item)
	}

	// 检查游标是否遇到错误
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over cursor: %v", err)
	}

	return items, nil
}

func (dao *CollectionQuestionInterfaceImpl) DeleteBookMarkQuestion(userId string, collectionId int64, questionId int64) error {
	// 获取 MongoDB 客户端连接

	// 选择特定的数据库和集合
	collection := db.MongoClient.Database("interview_guide").Collection("collection_items")

	// 创建查询条件
	filter := bson.M{
		"openid":        userId,       // 用户 ID
		"collection_id": collectionId, // 集合 ID
		"question_id":   questionId,   // 问题 ID
	}

	// 执行删除操作
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("failed to delete item: %v", err)
	}

	// 检查是否找到文档并删除
	if result.DeletedCount == 0 {
		return fmt.Errorf("no document found with given filters")
	}

	return nil
}

func (dao *CollectionQuestionInterfaceImpl) AddBookMarkQuestion(openId string, collectionID int64, questionId int64) error {
	collection := db.MongoClient.Database("interview_guide").Collection("collection_items")
	questionCollection := db.MongoClient.Database("interview_guide").Collection("question")
	fmt.Println(questionId)
	filter := bson.M{"question_id": questionId}
	one := questionCollection.FindOne(context.TODO(), filter)
	if errors.Is(one.Err(), mongo.ErrNoDocuments) {
		return fmt.Errorf("no document found with given questionId")
	}
	var question = model.QuestionModel{}

	err := one.Decode(&question)
	if err != nil {
		return err
	}

	// 创建插入文档的数据
	document := bson.M{
		"openid":        openId,         // 用户ID
		"collection_id": collectionID,   // 集合ID
		"question_id":   questionId,     // 问题ID
		"title":         question.Title, // 书签标题
		"content":       question.Content,
		"created_at":    time.Now(), // 创建时间
	}

	// 执行插入操作
	_, err = collection.InsertOne(context.TODO(), document)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("failed to insert bookmark item: %v", err)
	}

	return nil
}

func (dao *CollectionQuestionInterfaceImpl) AddQuestionCollection(collection *model.BookmarkCollectionModel) error {
	print("ooxx")

	collections := db.MongoClient.Database("interview_guide").Collection("collections")
	impl := CounterImpl{}
	print("ooxx")
	id, err2 := impl.GetAndIncrease("collections")
	fmt.Println(id)
	if err2 != nil {
		print(err2.Error())
		return err2
	}
	print("ooxx")

	// 创建插入文档的数据
	document := bson.M{
		"openid":          collection.OpenId,         // 用户 ID
		"collection_id":   id,                        // 集合 ID
		"collection_name": collection.CollectionName, // 集合名称
		"description":     "",                        // 集合描述
		"created_at":      time.Now(),                // 创建时间
	}

	// 执行插入操作
	_, err := collections.InsertOne(context.TODO(), document)
	if err != nil {
		print(err.Error())
		return err
	}

	return nil
}

func (dao *CollectionQuestionInterfaceImpl) DeleteBookMarkCollection(userID string, collectionID int64) error {
	collection := db.MongoClient.Database("interview_guide").Collection("collections")

	// 创建删除条件
	filter := bson.M{
		"user_id":       userID,       // 用户 ID
		"collection_id": collectionID, // 集合 ID
	}

	// 执行删除操作
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("failed to delete bookmark collection: %v", err)
	}

	// 检查是否找到文档并删除
	if result.DeletedCount == 0 {
		return fmt.Errorf("no collection found with the given user_id and collection_id")
	}

	return nil
}
func (dao *CollectionQuestionInterfaceImpl) GetCollectionItemsByCollectionAndTime(openId string, lastQuestionId int64, isDescending bool, collectionId int64) (*[]model.BookmarkQuestionModel, error) {
	const pageSize = 10

	// 获取 MongoDB 客户端连接

	// 选择特定的数据库和集合
	collection := db.MongoClient.Database("interview_guide").Collection("collection_items")
	filter := bson.M{}
	if collectionId == -1 {
		filter = bson.M{
			"question_id": bson.M{"$gt": lastQuestionId}, // 查询 question_id 大于上一页的最后一个 question_id
			"openid":      bson.M{"$eq": openId},
		}
	} else {
		filter = bson.M{
			"question_id":   bson.M{"$gt": lastQuestionId}, // 查询 question_id 大于上一页的最后一个 question_id
			"openid":        bson.M{"$eq": openId},
			"collection_id": collectionId,
		}

	}

	// 设置分页和排序选项
	findOptions := options.Find()
	fmt.Println(isDescending)
	if isDescending {
		findOptions.SetSort(bson.D{{"created_at", -1}})
	} else {
		findOptions.SetSort(bson.D{{"created_at", 1}})
	}
	findOptions.SetLimit(int64(pageSize)) // 每页限制 pageSize 个文档

	// 查询结果
	cursor, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch collection items: %v", err)
	}
	defer cursor.Close(context.TODO())

	// 遍历查询结果
	var items []model.BookmarkQuestionModel
	for cursor.Next(context.TODO()) {
		var item model.BookmarkQuestionModel
		if err := cursor.Decode(&item); err != nil {
			log.Printf("failed to decode item: %v", err)
			continue
		}
		items = append(items, item)
	}

	// 检查游标是否出错
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return &items, nil
}

func (dao *CollectionQuestionInterfaceImpl) GetCollectionItemsByTag(openId string, questionId int64, tag int64, collectionId int64) (*[]model.BookmarkQuestionModel, error) {
	// 每页的项目数量
	const pageSize = 10

	collection := db.MongoClient.Database("interview_guide").Collection("collection_items")

	filter := bson.M{
		"question_id": bson.M{"$gt": questionId}, // 查询 question_id 大于上一页的最后一个 question_id
		"openid":      bson.M{"$eq": openId},
		"tag":         tag, // 项目类别

	}
	// 设置分页选项
	findOptions := options.Find()

	findOptions.SetLimit(int64(pageSize)) // 每页限制 pageSize 个文档

	// 查询结果
	cursor, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch collection items: %v", err)
	}
	defer cursor.Close(context.TODO())

	var items []model.BookmarkQuestionModel
	for cursor.Next(context.TODO()) {
		var item model.BookmarkQuestionModel
		if err := cursor.Decode(&item); err != nil {
			log.Printf("failed to decode item: %v", err)
			continue
		}
		items = append(items, item)
	}

	// 检查游标是否出错
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}
	return &items, nil
}
