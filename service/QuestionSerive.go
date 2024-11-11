package service

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func GetQuestionsByIdHandler(writer http.ResponseWriter, request *http.Request) {
	var question_id = -1
	var err error
	if request.URL.Query().Get("question_id") != "" {
		question_id, err = strconv.Atoi(request.URL.Query().Get("question_id"))
	}
	if err != nil {
	}
	questions, err := questionImp.GetQuestionsById(int64(question_id))
	writer.Header().Set("Content-Type", "application/json")
	marshal, err := json.Marshal(JsonResult{Code: 1, Data: *questions})
	_, err = writer.Write(marshal)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)

}

func GetQuestionsByTag(writer http.ResponseWriter, request *http.Request) {

}

func GetHottestQuestionHandler(writer http.ResponseWriter, request *http.Request) {

}
