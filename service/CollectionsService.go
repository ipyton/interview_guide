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
	userId := r.URL.Query().Get("user_id")
	if err != nil || userId == "" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	items, err := bookmarkCollectionDAO.GetItemsInCollection(userId, int64(collectionID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
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
		CollectionID string `json:"collection_id"`
		QuestionID   string `json:"question_id"`
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

	// 验证参数
	if req.OpenID == "" || req.CollectionID == "" || req.QuestionID == "" {
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
	userId := r.URL.Query().Get("user_id")
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
	req, err := processGetItemsRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// 验证参数
	if req.OpenID == "" || req.Category == "" || req.PageNumber < 1 {
		http.Error(w, "Missing or invalid parameters", http.StatusBadRequest)
		return
	}
	items, err := bookmarkCollectionDAO.GetCollectionItemsByTime(req.OpenID, req.PageNumber)
	marshal, err := json.Marshal(JsonResult{Code: 1, Data: *items})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(marshal)
	w.WriteHeader(http.StatusOK)
	// 输出解析结果（这里可以根据业务需求继续处理 req）
}

func GetCollectionItemsByCategory(w http.ResponseWriter, r *http.Request) {

}
