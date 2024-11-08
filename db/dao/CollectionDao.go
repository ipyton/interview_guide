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

func (dao *CollectionQuestionInterfaceImpl) IsResourceCollected(userId string, resourceId int) (bool, error) {
	return true, errors.New("not implement")
}

func (dao *CollectionQuestionInterfaceImpl) GetCollections(openId string) (*[]model.BookmarkCollectionModel, error) {
	collection := db.MongoClient.Database("interview_guide").Collection("collections")

	// 定义查询过滤器
	filter := bson.M{"openId": openId}

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

func (dao *CollectionQuestionInterfaceImpl) GetItemsInCollection(openId string, collectionID int) ([]*model.BookmarkQuestionModel, error) {
	var items []*model.BookmarkQuestionModel

	collection := db.MongoClient.Database("interview_guide").Collection("collection_items")

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
		items = append(items, &item)
	}

	// 检查游标是否遇到错误
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over cursor: %v", err)
	}

	return items, nil
}

func (dao *CollectionQuestionInterfaceImpl) DeleteBookMarkQuestion(userId string, collectionId int, questionId int) error {
	// 获取 MongoDB 客户端连接

	// 选择特定的数据库和集合
	collection := db.MongoClient.Database("interview_guide").Collection("collection_items")

	// 创建查询条件
	filter := bson.M{
		"user_id":       userId,       // 用户 ID
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

// 这个做成已有的才能插入
func (dao *CollectionQuestionInterfaceImpl) AddBookMarkQuestion(openId string, collectionID string, questionId string) error {
	collection := db.MongoClient.Database("interview_guide").Collection("collection_items")
	questionCollection := db.MongoClient.Database("interview_guide").Collection("questions")

	filter := bson.M{"questionId": openId}
	one := questionCollection.FindOne(context.TODO(), filter)
	if errors.Is(one.Err(), mongo.ErrNoDocuments) {
		return fmt.Errorf("no document found with given openId")
	}
	var question = model.QuestionModel{}

	err := one.Decode(&question)
	if err != nil {
		return err
	}

	// 创建插入文档的数据
	document := bson.M{
		"user_id":       openId,         // 用户ID
		"collection_id": collectionID,   // 集合ID
		"question_id":   questionId,     // 问题ID
		"title":         question.Title, // 书签标题
		"created_at":    time.Now(),     // 创建时间
	}

	// 执行插入操作
	_, err = collection.InsertOne(context.TODO(), document)
	if err != nil {
		return fmt.Errorf("failed to insert bookmark item: %v", err)
	}

	return nil
}

func (dao *CollectionQuestionInterfaceImpl) AddQuestionCollection(userId string, collection *model.BookmarkCollectionModel) error {
	collections := db.MongoClient.Database("interview_guide").Collection("collection")

	// 创建插入文档的数据
	document := bson.M{
		"user_id":       userId,                    // 用户 ID
		"collection_id": collection.CollectionID,   // 集合 ID
		"name":          collection.CollectionName, // 集合名称
		"description":   collection.Description,    // 集合描述
		"created_at":    collection.CreateAt,       // 创建时间
		"updated_at":    collection.UpdateAt,       // 更新时间
	}

	// 执行插入操作
	_, err := collections.InsertOne(context.TODO(), document)
	if err != nil {
		return fmt.Errorf("failed to insert bookmark collection: %v", err)
	}

	return nil
}

func (dao *CollectionQuestionInterfaceImpl) DeleteBookMarkCollection(userID string, collectionID int) error {
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
func (dao *CollectionQuestionInterfaceImpl) GetCollectionItemsByTime(openId string, pageNumber int) (*[]model.BookmarkQuestionModel, error) {
	const pageSize = 10

	// 获取 MongoDB 客户端连接

	// 选择特定的数据库和集合
	collection := db.MongoClient.Database("interview_guide").Collection("collection_items")

	// 创建查询条件
	filter := bson.M{"user_id": openId}

	// 设置分页和排序选项
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"created_at", -1}})         // 按 created_at 降序排序
	findOptions.SetSkip(int64((pageNumber - 1) * pageSize)) // 跳过前 (pageNumber-1) 页的数据
	findOptions.SetLimit(int64(pageSize))                   // 每页限制 pageSize 个文档

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

func (dao *CollectionQuestionInterfaceImpl) GetCollectionItemsByCategory(openId string, category string, pageNumber int) (*[]model.BookmarkQuestionModel, error) {
	// 每页的项目数量
	const pageSize = 10

	collection := db.MongoClient.Database("interview_guide").Collection("collection_items")

	// 创建查询条件
	filter := bson.M{
		"user_id":  openId,   // 用户 ID
		"category": category, // 项目类别
	}

	// 设置分页选项
	findOptions := options.Find()
	findOptions.SetSkip(int64((pageNumber - 1) * pageSize)) // 跳过前 (pageNumber - 1) 页的数据
	findOptions.SetLimit(int64(pageSize))                   // 每页限制 pageSize 个文档

	// 查询结果
	cursor, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch collection items by category: %v", err)
	}
	defer cursor.Close(context.TODO())

	// 遍历查询结果并解码
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
