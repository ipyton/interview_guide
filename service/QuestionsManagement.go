package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
)

var searchDao = dao.SearchDaoImpl{}

func UpsertQuestions(w http.ResponseWriter, r *http.Request) {
	res := &JsonResult{}
	res.Data = "success"
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	var err error
	var question model.QuestionModel

	if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		fmt.Println(err.Error())
		http.Error(w, "Failed to encode posts", http.StatusInternalServerError)
		return
	}
}

func GetQuestionByIdHandler(w http.ResponseWriter, r *http.Request) {
	res := &JsonResult{}
	w.Header().Set("Content-Type", "application/json")

	var question_id = -1
	var err error
	//var openid = r.Header.Get("openid")
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
		posts, _ := questionImp.GetQuestionById(int64(question_id))

		res.Code = 1
		res.Data = posts
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, "Failed to encode posts", http.StatusInternalServerError)
		return
	}
}

func truncateByRune(s string, limit int) string {
	runes := []rune(s)
	if len(runes) > limit {
		return string(runes[:limit]) + "..." // 添加省略号以表示截断
	}
	return s
}

func UpsertQuestionsByFile(w http.ResponseWriter, r *http.Request) {
	// 解析 multipart 表单数据，限制最大上传文件大小
	err := r.ParseMultipartForm(10 << 20) // 限制为 10 MB

	if err != nil {
		fmt.Println(err.Error())

		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	// 获取上传的文件
	file, _, err := r.FormFile("file")
	if err != nil {
		fmt.Println(err.Error())

		http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// 读取文件内容到内存中
	fileContents, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Error reading the file", http.StatusInternalServerError)
		return
	}
	type QuestionRequest struct {
		Question string   `json:"question"`
		Answer   string   `json:"answer"`
		Tags     []string `json:"tags"`
		ClassId  int64    `json:"class_id"`
	}
	// 解析 JSON 文件内容到切片
	var dataList []QuestionRequest
	err = json.Unmarshal(fileContents, &dataList)
	if err != nil {
		fmt.Println(err.Error())

		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}
	var list []model.QuestionModel
	for _, data := range dataList {
		increase, err := dao.CounterImpl{}.GetAndIncrease("questions")
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, "Error increasing questions", http.StatusInternalServerError)
			return
		}
		list = append(list, model.QuestionModel{ID: increase, Title: data.Question, Tags: data.Tags,
			Details: data.Answer, Content: truncateByRune(data.Answer, 20), ClassId: data.ClassId})
	}
	err = questionImp.BatchAdd(&list)
	if err != nil {
		print(err.Error())
		http.Error(w, "Error adding questions", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func AdviceAQuestion(w http.ResponseWriter, r *http.Request) {
	questionModel := model.AdvisedQuestions{}
	all, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Error reading body", http.StatusInternalServerError)
	}
	err = json.Unmarshal(all, &questionModel)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = questionImp.AdviceQuestion(questionModel)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Error adding questions", http.StatusInternalServerError)
	}

}

func GetAdvisedQuestions(w http.ResponseWriter, r *http.Request) {
	questions, err := questionImp.GetAdvisedQuestions()
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Error getting questions", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(questions)

}

type ApproveRequest struct {
	QuestionId int64 `json:"question_id"`
}

func ApproveAQuestion(w http.ResponseWriter, r *http.Request) {
	all, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Error reading body", http.StatusInternalServerError)
	}
	approved := ApproveRequest{}
	err = json.Unmarshal(all, &approved)
	if err != nil && approved.QuestionId == 0 {
		fmt.Println(err.Error())
		http.Error(w, "Error reading body", http.StatusInternalServerError)
	}
	questionImp.ApproveAQuestion(approved.QuestionId)
}
