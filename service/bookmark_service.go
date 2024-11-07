package service

import (
	"encoding/json"
	"strconv"

	"net/http"
	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
)

var bookmarkCollectionDAO dao.CollectionInterface = &dao.CollectionInterfaceImpl{}

func GetBookmarkCollections(w http.ResponseWriter, r *http.Request) {
	collections, err := bookmarkCollectionDAO.GetCollections()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(collections)
}

func GetBookmarkItems(w http.ResponseWriter, r *http.Request) {
	// collectionID := r.URL.Query().Get("collection_id")
	var err error
	collectionID, err := strconv.Atoi(r.URL.Query().Get("collection_id"))
	userId := r.URL.Query().Get("user_id")
	if err != nil || userId == "" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	items, err := bookmarkCollectionDAO.GetItemsInCollection(userId, collectionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func DeleteBookmarkItem(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")
	questionId, err := strconv.Atoi(r.URL.Query().Get("question_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	collectionId, err := strconv.Atoi(r.URL.Query().Get("collection_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = bookmarkCollectionDAO.DeleteBookMarkItem(userId, collectionId, questionId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // No content to return
}

func AddBookmarkItem(w http.ResponseWriter, r *http.Request) {
	var item model.BookmarkItemModel
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := bookmarkCollectionDAO.AddBookMarkItem(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated) // Resource created
}

func AddBookmarkCollection(w http.ResponseWriter, r *http.Request) {
	var collection model.BookmarkCollectionModel
	if err := json.NewDecoder(r.Body).Decode(&collection); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := bookmarkCollectionDAO.AddBookMarkCollection(&collection)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated) // Resource created
}

func DelBookmarkCollection(w http.ResponseWriter, r *http.Request) {
	collectionID := r.URL.Query().Get("collection_id")
	userId := r.URL.Query().Get("user_id")
	num, err := strconv.Atoi(collectionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = bookmarkCollectionDAO.DeleteBookMarkCollection(userId, num)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // No content to return
}
