package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"wxcloudrun-golang/db/model"
)

func UpsertQuestions(w http.ResponseWriter, r *http.Request) {
	res := &JsonResult{}
	res.Data = "success"
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	var err error
	var question model.QuestionModel

	if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	err = questionImp.UpsertQuestion(&question)
	if err != nil {
		res.Code = -1
		res.ErrorMsg = "Failed to get page number"
		fmt.Println(err.Error())
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, "Failed to encode posts", http.StatusInternalServerError)
		return
	}
}

func GetQuestionsByIdHandler(w http.ResponseWriter, r *http.Request) {
	res := &JsonResult{}
	w.Header().Set("Content-Type", "application/json")

	var question_id = -1
	var err error
	if r.URL.Query().Get("question_id") != "" {
		question_id, err = strconv.Atoi(r.URL.Query().Get("question_id"))
	}

	if err != nil {
		res.Code = -1
		res.ErrorMsg = "Failed to get page number"
	}
	if question_id == -1 {
		res.Code = -1
		res.ErrorMsg = "Invalid question id"
	} else {
		posts, _ := questionImp.GetQuestionById(question_id)
		res.Code = 1
		res.Data = posts
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, "Failed to encode posts", http.StatusInternalServerError)
		return
	}

}
