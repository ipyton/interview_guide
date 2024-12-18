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

func (this *SearchDaoImpl) GetSuggestions(keyword string) (SuggestResponse, error) {
	reqBody := fmt.Sprintf(`{
		"suggest": {
			"text": "%s", 
			"completion": {
				"field": "suggest",
				"size": 5
			}
		}
	}`, keyword)
	response, err := db.ElasticClient.Search(
		db.ElasticClient.Search.WithContext(context.Background()),
		db.ElasticClient.Search.WithIndex("interview-questions"),
		db.ElasticClient.Search.WithBody(strings.NewReader(reqBody)),
		db.ElasticClient.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		log.Fatalf("Error executing the query: %s", err)
	}
	var result SuggestResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Fatalf("Error decoding the response body: %s", err)
		return result, err
	}
	return result, err
}

func removeDuplicates(input []string) []string {
	seen := make(map[string]struct{})
	var result []string

	for _, item := range input {
		if _, ok := seen[item]; !ok {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

// 使用 Elasticsearch 自带的 _analyze API 分析文本并返回分词结果
func analyzeText(text string) ([]string, error) {
	// 构造 _analyze 请求
	analyzeRequest := esapi.IndicesAnalyzeRequest{
		Index:  "interview-questions",                                                             // 索引名称
		Body:   strings.NewReader(fmt.Sprintf(`{"text": "%s", "analyzer": "ik_max_word"}`, text)), // 要分析的文本
		Pretty: true,                                                                              // 使用 ik_max_word 分词器
	}

	// 执行 _analyze 请求
	res, err := analyzeRequest.Do(context.Background(), db.ElasticClient)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// 解析响应体
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	// 获取分词结果
	tokens, ok := result["tokens"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	// 提取所有 token（分词结果）
	var words []string
	for _, token := range tokens {
		tokenMap, ok := token.(map[string]interface{})
		if !ok {
			continue
		}
		if term, ok := tokenMap["token"].(string); ok {
			words = append(words, term)
		}
	}
	return words, nil
}

func (this *SearchDaoImpl) UpsertQuestionIndex(question model.QuestionModel) error {
	// 初始化 jieba 分词器

	// 对 title, content, details 字段进行分词
	titleTokens, err := analyzeText(question.Title) // 使用全模式分词
	contentTokens, err := analyzeText(question.Content)
	detailsTokens, err := analyzeText(question.Details)
	if err != nil {

		return err
	}
	// 将分词结果转换为数组形式，用于 suggest 字段
	suggestInput := append(append([]string{}, titleTokens...), contentTokens...)
	suggestInput = append(suggestInput, detailsTokens...)
	suggestInput = removeDuplicates(suggestInput)
	// 准备 Elasticsearch 请求体
	reqBody, err := json.Marshal(map[string]interface{}{
		"doc": map[string]interface{}{
			"title":       question.Title,
			"class_id":    question.ClassId,
			"type":        question.Type,
			"content":     question.Content,
			"details":     question.Details,
			"author_id":   question.AuthorID,
			"author_name": question.AuthorName,
			"avatar":      question.Avatar,
			"likes":       question.Likes,
			"views":       question.Views,
			"difficulty":  question.Difficulty,
			"tags":        question.Tags,
			"suggest": map[string]interface{}{
				"input": suggestInput, // 使用分词结果作为补全输入
			},
		},
		"upsert": map[string]interface{}{
			"title":       question.Title,
			"class_id":    question.ClassId,
			"type":        question.Type,
			"content":     question.Content,
			"details":     question.Details,
			"author_id":   question.AuthorID,
			"author_name": question.AuthorName,
			"avatar":      question.Avatar,
			"likes":       question.Likes,
			"views":       question.Views,
			"difficulty":  question.Difficulty,
			"tags":        question.Tags,
			"suggest": map[string]interface{}{
				"input": suggestInput, // 使用分词结果作为补全输入
			},
		},
	})
	if err != nil {
		fmt.Println("Error marshalling request body:", err)
		return err
	}
	fmt.Printf("Request Body: %s\n", reqBody)

	// 获取文档 ID 并格式化为字符串
	formatInt := strconv.FormatInt(question.ID, 10)
	request := esapi.UpdateRequest{
		Index:      "interview-questions", // 索引名称
		DocumentID: formatInt,
		Body:       bytes.NewReader(reqBody),
		Refresh:    "true",
	}

	// 执行更新或插入请求
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

	res, err := db.ElasticClient.Search(
		db.ElasticClient.Search.WithContext(context.Background()),
		db.ElasticClient.Search.WithIndex("interview-questions"), // 替换为你的索引名称
		db.ElasticClient.Search.WithBody(strings.NewReader(string(searchQuery))),
		db.ElasticClient.Search.WithPretty(),
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
