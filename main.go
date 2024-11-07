package main

import (
	"context"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/service"
)

func main() {
	db.InitMongo()
	//http.HandleFunc("/", service.IndexHandler)
	//http.HandleFunc("/api/count", service.CounterHandler)
	http.HandleFunc("/home", service.HomeHandler)
	http.HandleFunc("/home/get", service.GetHottestHandler)
	http.HandleFunc("/questions/get", service.GetQuestionsByPageHandler)
	http.HandleFunc("/questions/getById", service.GetQuestionsByIdHandler)
	http.HandleFunc("/questions/upsert", service.UpsertQuestions)
	http.HandleFunc("/collections/items/get", service.GetBookmarkItems)
	http.HandleFunc("/collections/items/delete", service.DeleteBookmarkItem)
	http.HandleFunc("/collections/item/add", service.AddBookmarkItem)
	http.HandleFunc("/collections/collection/add", service.AddBookmarkCollection)
	http.HandleFunc("/collections/collection/delete", service.DelBookmarkCollection)
	http.HandleFunc("/collections/collection/get", service.GetBookmarkCollections)
	cors := handlers.CORS(handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"get", "post", "put", "patch", "delete"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}), // Add other headers as needed
	)

	log.Fatal(http.ListenAndServe(":5050", cors(http.DefaultServeMux)))
	var err error
	if err = db.MongoClient.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}
