package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
)

var bookmarkCollectionDAO dao.CollectionQuestionInterface = &dao.CollectionQuestionInterfaceImpl{}

type CollectionRequest struct {
	OpenID     string `json:"open_id"`
	Category   string `json:"category"`
	PageNumber int64  `json:"page_number"`
}

func processGetItemsRequest(r *http.Request) (*CollectionRequest, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read request body")
	}
	defer r.Body.Close()

	// 解析 JSON 到 CollectionRequest 结构体
	var req CollectionRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, fmt.Errorf("Invalid JSON format")

	}
	return &req, nil
}

func GetBookmarkCollections(w http.ResponseWriter, r *http.Request) {
	openId := r.Header.Get("openid")
	collections, err := bookmarkCollectionDAO.GetCollections(openId)
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
	openid := r.Header.Get("openid")
	if openid == "" {
		http.Error(w, "openid required", http.StatusBadRequest)
	}
	if err != nil {
		http.Error(w, "Error request", http.StatusBadRequest)
		return
	}
	items, err := bookmarkCollectionDAO.GetItemsInCollection(openid, int64(collectionID))
	fmt.Println(items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(*items)

}

func DeleteBookmarkItem(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")
	questionId, err := strconv.ParseInt(r.URL.Query().Get("question_id"), 10, 64)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	collectionId, err := strconv.ParseInt(r.URL.Query().Get("collection_id"), 10, 64)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = bookmarkCollectionDAO.DeleteBookMarkQuestion(userId, collectionId, questionId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // No content to return
}

func AddBookmarkItem(w http.ResponseWriter, r *http.Request) {
	type BookmarkRequest struct {
		OpenID       string `json:"open_id"`
		CollectionID int64  `json:"collection_id"`
		QuestionID   int64  `json:"question_id"`
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 解析 JSON 到 BookmarkRequest 结构体
	var req BookmarkRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	req.OpenID = r.Header.Get("openid")
	// 验证参数
	if req.OpenID == "" {
		http.Error(w, "Missing or invalid parameters", http.StatusBadRequest)
		return
	}
	err = bookmarkCollectionDAO.AddBookMarkQuestion(req.OpenID, req.CollectionID, req.QuestionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated) // Resource created
}

func AddBookmarkCollection(w http.ResponseWriter, r *http.Request) {
	var collection model.BookmarkCollectionModel
	println(r.Body)
	if err := json.NewDecoder(r.Body).Decode(&collection); err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(collection.CollectionName)
	collection.OpenId = r.Header.Get("openid")
	err := bookmarkCollectionDAO.AddQuestionCollection(&collection)
	if err != nil {
		fmt.Println("adasdasdasd")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated) // Resource created
}

func DelBookmarkCollection(w http.ResponseWriter, r *http.Request) {
	collectionID := r.URL.Query().Get("collection_id")
	userId := r.Header.Get("openid")
	num, err := strconv.ParseInt(collectionID, 10, 64)
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

func GetCollectionItemsByTime(w http.ResponseWriter, r *http.Request) {
	openid := r.Header.Get("openid")
	question_id, err := strconv.ParseInt(r.URL.Query().Get("question_id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	isDescending, err := strconv.ParseBool(r.URL.Query().Get("isDescending"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// 验证参数
	if openid == "" {
		http.Error(w, "Missing or invalid parameters", http.StatusBadRequest)
		return
	}
	fmt.Println(openid)
	fmt.Println(question_id)
	fmt.Println(isDescending)
	items, err := bookmarkCollectionDAO.GetCollectionItemsByTime(openid, question_id, isDescending)
	marshal, err := json.Marshal(JsonResult{Code: 1, Data: *items})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(marshal)

}

func GetCollectionItemsByCategory(w http.ResponseWriter, r *http.Request) {

}
