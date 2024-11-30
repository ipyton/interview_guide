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
)

type SearchRequest struct {
	Keyword string `json:"keyword"`
}

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
	res, err := request.Do(context.Background(), db.Client)
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
