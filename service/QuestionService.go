package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"wxcloudrun-golang/db/model"
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
	if questions == nil || len(*questions) == 0 {
		questions = &[]model.QuestionModel{} // Initialize as an empty slice if nil
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

func Rate(writer http.ResponseWriter, request *http.Request) {
	// Parse the JSON body from the incoming request
	var userRate model.UserRate
	openid := request.URL.Query().Get("openid")
	userRate.OpenId = openid
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&userRate); err != nil {
		http.Error(writer, fmt.Sprintf("Failed to decode request body: %s", err), http.StatusBadRequest)
		return
	}

	err := questionImp.RateQuestion(userRate)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Send response
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	response := map[string]string{"status": "success"}
	json.NewEncoder(writer).Encode(response)
}

func GetRatings(writer http.ResponseWriter, request *http.Request) {
	questionIdStr := request.URL.Query().Get("questionId")
	if questionIdStr == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	question_id, err := strconv.ParseInt(questionIdStr, 10, 64)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
	}
	ratings, err := questionImp.GetRatings(question_id)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(writer).Encode(ratings)

}

func SeeBefore(writer http.ResponseWriter, request *http.Request) {
	var seeBefore model.SeeBeforeCount
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&seeBefore); err != nil {
		http.Error(writer, fmt.Sprintf("Failed to decode request body: %s", err), http.StatusBadRequest)
		return
	}
	questionImp.SeeBefore(seeBefore)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	response := map[string]string{"status": "success"}
	json.NewEncoder(writer).Encode(response)
}
