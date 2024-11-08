package dao

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

type CounterImpl struct {
	CounterInterface
}

func (counterImpl CounterImpl) GetAndIncrease(collectionName string) (int, error) {
	collection := db.MongoClient.Database("interview_guide").Collection("counters")
	filter := bson.M{"_id": collectionName}
	update := bson.M{"$inc": bson.M{"value": 1}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedDoc model.CounterModel

	err := collection.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&updatedDoc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// Counter does not exist, initialize it
			_, err := collection.InsertOne(context.TODO(), bson.M{"_id": collectionName, "seq": 1})
			if err != nil {
				return 0, fmt.Errorf("failed to initialize counter: %v", err)
			}
			return 1, nil // First value is 1
		}
		return 0, fmt.Errorf("failed to update counter: %v", err)
	}

	return updatedDoc.Count, nil
}
