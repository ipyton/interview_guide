package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"wxcloudrun-golang/db"
)

func getNextSequenceValue(sequenceName string) (int64, error) {
	collection := db.MongoClient.Database("interview_guide").Collection("counters")

	filter := bson.M{"_id": sequenceName}
	update := bson.M{"$inc": bson.M{"value": 1}}
	options := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var result struct {
		SequenceValue int64 `bson:"value"`
	}

	err := collection.FindOneAndUpdate(context.TODO(), filter, update, options).Decode(&result)
	if err != nil {
		return 0, err
	}

	return result.SequenceValue, nil
}
