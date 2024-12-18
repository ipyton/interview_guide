package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"wxcloudrun-golang/db/model"
)

func GetFeedback(w http.ResponseWriter, r *http.Request) {
	questions, err := questionImp.GetAdvisedQuestions()
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Error getting questions", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(questions)
}

func SendFeedback(w http.ResponseWriter, r *http.Request) {
	questionModel := model.AdvisedQuestions{}
	all, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Error reading body", http.StatusInternalServerError)
	}
	err = json.Unmarshal(all, &questionModel)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = questionImp.AdviceQuestion(questionModel)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Error adding questions", http.StatusInternalServerError)
	}
}

func ReplyFeedback(w http.ResponseWriter, r *http.Request) {

}
