package service

import (
	"encoding/json"
	"net/http"
	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
)

var bookmarkCollectionDAO dao.CollectionDAO = &dao.CollectionDAOImpl{}

func GetBookmarkCollections(w http.ResponseWriter, r *http.Request) {
	collections, err := bookmarkCollectionDAO.GetBookCollections()
	if err != nil {
		http.Error(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(collections)
}

func GetBookmarkItems(w http.ResponseWriter, r *http.Request) {
	collectionID := r.URL.Query().Get("collection_id")
	items, err := bookCollectionDAO.GetBookMarkItems(collectionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func DeleteBookmarkItem(w http.ResponseWriter, r *http.Request) {
	resourceID := r.URL.Query().Get("resource_id")
	err := bookmarkCollectionDAO.DeleteBookMarkItem(resourceID)
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

	err := bookCollectionDAO.AddBookMarkItem(&item)
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

	err := bookCollectionDAO.AddBookMarkCollection(&collection)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated) // Resource created
}

func DelBookmarkCollection(w http.ResponseWriter, r *http.Request) {
	collectionID := r.URL.Query().Get("collection_id")
	err := bookCollectionDAO.DeleteBookMarkCollection(collectionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // No content to return
}
