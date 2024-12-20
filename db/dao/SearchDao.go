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
	"unicode"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

type SearchDaoImpl struct {
	SearchDaoInterface
}

type SuggestionOption struct {
	Text    string              `json:"text"`
	Score   float64             `json:"score"`
	Payload interface{}         `json:"payload"`
	Source  model.QuestionModel `json:"_source"`
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
	fmt.Println("Keyword:", keyword)

	// 构建请求体
	reqBody := map[string]interface{}{
		"suggest": map[string]interface{}{
			"song-suggest": map[string]interface{}{
				"prefix": keyword,
				"completion": map[string]interface{}{
					"field": "suggest",
				},
			},
		},
	}

	// 将请求体转换为 JSON 字符串
	body, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatalf("Error marshalling the request body: %s", err)
		return SuggestResponse{}, err
	}

	// 调用 Elasticsearch 的搜索 API
	response, err := db.ElasticClient.Search(
		db.ElasticClient.Search.WithContext(context.Background()),
		db.ElasticClient.Search.WithIndex("interview-questions"),
		db.ElasticClient.Search.WithBody(bytes.NewReader(body)), // 使用 JSON 字符串或字节切片
		db.ElasticClient.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		log.Fatalf("Error executing the query: %s", err)
		return SuggestResponse{}, err
	}
	defer response.Body.Close() // 记得关闭响应体

	// 打印响应的状态码，便于调试
	fmt.Println("Response Status:", response.Status)

	// 打印响应的内容
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %s", err)
	}
	fmt.Println("Response Body:", string(bodyBytes))

	// 解析 JSON 响应
	var result SuggestResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		log.Fatalf("Error decoding the response body: %s", err)
		return SuggestResponse{}, err
	}

	// 返回解析后的结果
	return result, nil
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
func escapeIllegalChars(input string) string {
	var result bytes.Buffer

	for _, r := range input {
		// 如果字符是控制字符（包括换行符、回车符等），我们会将其转义
		if unicode.IsControl(r) {
			switch r {
			case '\n':
				result.WriteString("\\n")
			case '\r':
				result.WriteString("\\r")
			case '\t':
				result.WriteString("\\t")
			default:
				result.WriteString("\\u" + strings.ToUpper(string(r)))
			}
		} else {
			// 非控制字符直接写入结果
			result.WriteRune(r)
		}
	}

	return result.String()
}

// 使用 Elasticsearch 自带的 _analyze API 分析文本并返回分词结果
func analyzeText(text string) ([]string, error) {
	// 构造 _analyze 请求
	text = escapeIllegalChars(text)
	analyzeRequest := esapi.IndicesAnalyzeRequest{
		Index:  "interview-questions",                                                          // 索引名称
		Body:   strings.NewReader(fmt.Sprintf(`{"text": "%s", "analyzer": "ik_smart"}`, text)), // 要分析的文本
		Pretty: true,                                                                           // 使用 ik_max_word 分词器
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
		fmt.Println(result)
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

func concatAdjacentStrings(strs []string) []string {
	var result []string

	// 确保有至少两个字符串可供拼接
	for i := 0; i < len(strs)-1; i++ {
		// 拼接当前字符串和下一个字符串
		concatenated := strs[i] + strs[i+1]
		result = append(result, concatenated)
	}

	return result
}

func (this *SearchDaoImpl) UpsertQuestionIndex(question model.QuestionModel) error {
	// 初始化 jieba 分词器

	// 对 title, content, details 字段进行分词
	titleTokens, err := analyzeText(question.Title) // 使用全模式分词
	//contentTokens, err := analyzeText(question.Content)
	detailsTokens, err := analyzeText(question.Details)
	if err != nil {

		return err
	}
	// 将分词结果转换为数组形式，用于 suggest 字段
	suggestInput := append([]string{}, titleTokens...)
	suggestInput = append(suggestInput, detailsTokens...)
	suggestInput = concatAdjacentStrings(suggestInput)
	suggestInput = removeDuplicates(suggestInput)
	// 准备 Elasticsearch 请求体
	reqBody, err := json.Marshal(map[string]interface{}{
		"doc": map[string]interface{}{
			"title":       question.Title,
			"question_id": question.ID,
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
			"question_id": question.ID,
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

func (this *SearchDaoImpl) SearchQuestions(keyword string, pageSize int64, page int64) (*SearchResult, error) {
	from := (page - 1) * pageSize

	// Construct the search query with pagination (from and size)
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  keyword,
				"fields": []string{"title", "details", "author_name", "content"},
			},
		},
		"size": pageSize, // Number of results per page
		"from": from,     // The starting point of the next page of results
	}

	// Marshal the query to JSON format
	searchQuery, err := json.Marshal(query)
	if err != nil {
		log.Fatal("Error marshaling query: ", err)
	}

	// Execute the search query
	res, err := db.ElasticClient.Search(
		db.ElasticClient.Search.WithContext(context.Background()),
		db.ElasticClient.Search.WithIndex("interview-questions"), // Replace with your index name
		db.ElasticClient.Search.WithBody(strings.NewReader(string(searchQuery))),
		db.ElasticClient.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error executing the search: %s", err)
	}
	defer res.Body.Close()

	// Output the search results
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

		// Process the returned results
		fmt.Println("Total hits:", result.Hits.Total.Value)
		for _, hit := range result.Hits.Hits {
			fmt.Printf("Found document: %+v\n", hit.Source)
		}
		return &result, nil
	}

	return nil, nil

}
