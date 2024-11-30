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
	formatInt := strconv.FormatInt(question.ID, 10)
	request := esapi.IndexRequest{
		Index:      "interview-questions", // 索引名称
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

func (this *SearchDaoImpl) SearchQuestions(keyword string) (*model.QuestionModel, error) {

	searchQuery := fmt.Sprintf(`{
		"query": {
			"match": {
				"content": "%s"  // 假设你要搜索的字段是 content
			}
		}
	}`, keyword)
	res, err := db.Client.Search(
		db.Client.Search.WithContext(context.Background()),
		db.Client.Search.WithIndex("interview-questions"), // 替换为你的索引名称
		db.Client.Search.WithBody(strings.NewReader(searchQuery)),
		db.Client.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error executing the search: %s", err)
	}
	defer res.Body.Close()

	// 输出搜索结果
	if res.IsError() {
		log.Printf("Error: %s", res.String())
	} else {
		// 解析并输出响应数据
		fmt.Println("Search results:")
		// 这里可以进一步解析 res.Body 获取返回的文档
	}
	all, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	question := model.QuestionModel{}
	err = json.Unmarshal(all, &question)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &question, nil

}
