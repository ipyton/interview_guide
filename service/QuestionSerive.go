package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func GetQuestionsByPagingHandler(writer http.ResponseWriter, request *http.Request) {
	var question_id int64
	var class_id int64
	var err error
	if request.URL.Query().Get("question_id") != "" {
		question_id, err = strconv.ParseInt(request.URL.Query().Get("question_id"), 10, 64)
		if request.URL.Query().Get("class_id") != "" {
			class_id, err = strconv.ParseInt(request.URL.Query().Get("class_id"), 10, 64)
		}
	}
	if err != nil {
		fmt.Println(err.Error())
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	questions, err := questionImp.GetQuestionsByPaging(question_id, class_id)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	marshal, err := json.Marshal(JsonResult{Code: 1, Data: *questions})
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(marshal)

}

func GetQuestionsByTag(writer http.ResponseWriter, request *http.Request) {

}

func GetHottestQuestionHandler(writer http.ResponseWriter, request *http.Request) {

}
