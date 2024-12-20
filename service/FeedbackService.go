package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
)

var feedBack = dao.FeedbackDao{}

func GetFeedback(w http.ResponseWriter, r *http.Request) {
	feedbacks, err := feedBack.GetFeedback()
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Error getting questions", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(feedbacks)
}

func SendFeedback(w http.ResponseWriter, r *http.Request) {
	feedback := model.Feedback{}
	all, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Error reading body", http.StatusInternalServerError)
	}
	err = json.Unmarshal(all, &feedback)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = feedBack.SaveFeedback(feedback)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Error adding questions", http.StatusInternalServerError)
	}
}

func ReplyFeedback(w http.ResponseWriter, r *http.Request) {

}

type DeleteRequest struct {
	ID int64 `json:"id"`
}

func DeleteFeedback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// 读取请求体
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// 反序列化到 DeleteRequest
	var req DeleteRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	feedBack.DeleteFeedback(req.ID)
	// 打印反序列化的结果（仅供测试使用）
	fmt.Printf("Received ID: %d\n", req.ID)

	// 这里可以根据反序列化的内容执行其他逻辑
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Delete request processed successfully"))
}
