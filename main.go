package main

import (
	"fmt"
	"log"
	"net/http"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/service"
)

func main() {
	if err := db.Init(); err != nil {
		panic(fmt.Sprintf("mysql init failed with %+v", err))
	}

	http.HandleFunc("/", service.IndexHandler)
	http.HandleFunc("/api/count", service.CounterHandler)
	http.HandleFunc("/home", service.HomeHandler)
	http.HandleFunc("/home/get", service.GetHottestHandler)
	http.HandleFunc("/collections/items/get", service.GetBookmarkItems)
	http.HandleFunc("/collections/items/delete", service.DeleteBookmarkItem)
	http.HandleFunc("/collections/item/add", service.AddBookmarkItem)
	http.HandleFunc("/collections/collection/add", service.AddBookmarkCollection)
	http.HandleFunc("/collections/collection/delete", service.DelBookmarkCollection)
	http.HandleFunc("/collections/collection/get", service.GetBookmarkCollections)

	log.Fatal(http.ListenAndServe(":80", nil))
}
