package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

const applyurl = "mongodb://192.168.31.75:27017/?replicaSet=rs0"

func InitMongo() {

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	//applyurl := fmt.Sprintf("mongodb://%s:%s/", "192.168.31.75", "27017")

	//credential := options.Credential{
	//	AuthMechanism: "SCRAM-SHA-256",
	//	// AuthMechanism: "SCRAM-SHA-1",
	//	Username:   "admin",
	//	Password:   "password",
	//	AuthSource: "admin",
	//}

	opts := options.Client().ApplyURI(applyurl).SetServerAPIOptions(serverAPI)
	//.SetAuth(credential)

	// Create a new client and connect to the server
	var err error
	MongoClient, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	// Send a ping to confirm a successful connection
	var result bson.M
	if err := MongoClient.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
}
