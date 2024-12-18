package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"io"
	"log"
	"net/http"
	"time"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/dao"
)

type SearchRequest struct {
	Keyword string `json:"keyword"`
}

var search dao.SearchDaoImpl

func GetResults(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		request := SearchRequest{}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(body, &request)
		fmt.Println(request.Keyword)
		questions, err := search.SearchQuestions(request.Keyword)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		marshal, err := json.Marshal(questions)
		if err != nil {
		}
		w.Write(marshal)
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

type InterviewQuestion struct {
	ID       string    `json:"id"`
	Question string    `json:"question"`
	Category string    `json:"category"`
	Level    string    `json:"level"`
	Created  time.Time `json:"created"`
}

func Testing(w http.ResponseWriter, r *http.Request) {
	question := InterviewQuestion{
		ID:       "1",
		Question: "What is the difference between a goroutine and a thread?",
		Category: "Golang",
		Level:    "Medium",
		Created:  time.Now(),
	}
	reqBody, err := json.Marshal(question)
	fmt.Printf("Request Body: %s\n", reqBody)

	if err != nil {
		fmt.Println(err)
	}

	request := esapi.IndexRequest{
		Index:      "interview-questions", // 索引名称
		DocumentID: question.ID,
		Body:       bytes.NewReader(reqBody),
		Refresh:    "true",
	}
	res, err := request.Do(context.Background(), db.ElasticClient)
	if err != nil {
		log.Fatalf("Error indexing document: %s", err.Error())
	}
	defer res.Body.Close()

	// 输出响应结果
	if res.IsError() {
		fmt.Printf("[%s] Error indexing document ID=%s\n", res.Status(), question.ID)
	} else {
		// 打印出返回的ID
		fmt.Printf("[%s] Document indexed successfully ID=%s\n", res.Status(), question.ID)
	}
}

type SuggestionOption struct {
	Text    string      `json:"text"`
	Score   float64     `json:"score"`
	Payload interface{} `json:"payload"`
}

type SuggestResponse struct {
	Suggest map[string][]struct {
		Text    string             `json:"text"`
		Offset  int                `json:"offset"`
		Length  int                `json:"length"`
		Options []SuggestionOption `json:"options"`
	} `json:"suggest"`
}

func GetSuggestions(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("keyword")
	suggestions, err := search.GetSuggestions(keyword)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(suggestions)

}
