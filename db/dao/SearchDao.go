package dao

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"io"
	"log"
	"strconv"
	"strings"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

type SearchDaoImpl struct {
	SearchDaoInterface
}

func (this *SearchDaoImpl) CreateQuestionIndex(question model.QuestionModel) error {

	reqBody, err := json.Marshal(question)
	fmt.Printf("Request Body: %s\n", reqBody)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(question.ID)
	formatInt := strconv.FormatInt(question.ID, 10)
	request := esapi.IndexRequest{
		Index:      "new-questions", // 索引名称
		DocumentID: formatInt,
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
	return nil
}

type SearchResult struct {
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		Hits []struct {
			Source model.QuestionModel `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func (this *SearchDaoImpl) SearchQuestions(keyword string) (*SearchResult, error) {

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  keyword,
				"fields": []string{"title", "details", "author_name", "content"},
			},
		},
	}

	searchQuery, err := json.Marshal(query)
	if err != nil {
		log.Fatal("Error marshaling query: ", err)
	}

	res, err := db.Client.Search(
		db.Client.Search.WithContext(context.Background()),
		db.Client.Search.WithIndex("new-questions"), // 替换为你的索引名称
		db.Client.Search.WithBody(strings.NewReader(string(searchQuery))),
		db.Client.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error executing the search: %s", err)
	}
	defer res.Body.Close()

	// 输出搜索结果
	if res.IsError() {
		log.Printf("Error: %s\n", res.String())
		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal("Error reading response body:", err)
		}
		log.Printf("Response body: %s", body)
	} else {
		all, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		var result SearchResult
		err = json.Unmarshal(all, &result)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		// 这里处理返回的结果
		fmt.Println("Total hits:", result.Hits.Total.Value)
		for _, hit := range result.Hits.Hits {
			fmt.Printf("Found document: %+v\n", hit.Source)
		}
		return &result, nil

	}
	return nil, nil

}
