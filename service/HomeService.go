package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}

var questionImp dao.QuestionInterface = &dao.QuestionInterfaceImpl{}

func generateHotPosts() []model.HotPost {
	var posts []model.HotPost
	for i := 1; i <= 10; i++ {
		post := model.HotPost{
			InfoID:      i,
			QuestionID:  1000 + i,
			Intro:       fmt.Sprintf("This is an introduction for post %d", i),
			AuthorID:    200 + i,
			Title:       fmt.Sprintf("Hot Post Title %d", i),
			Likes:       100 * i,
			PublishDate: time.Now().AddDate(0, 0, -i), // 发布日期是最近的 10 天
			Extra1:      fmt.Sprintf("Extra field 1 for post %d", i),
			Extra2:      fmt.Sprintf("Extra field 2 for post %d", i),
			Extra3:      fmt.Sprintf("Extra field 3 for post %d", i),
		}
		posts = append(posts, post)
	}
	return posts
}

func GetHottestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	posts := generateHotPosts()

	err := json.NewEncoder(w).Encode(posts)
	if err != nil {
		http.Error(w, "Failed to encode posts", http.StatusInternalServerError)
		return
	}
}

func GetQuestionsByPageHandler(w http.ResponseWriter, r *http.Request) {
	res := &JsonResult{}
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)

	body := make(map[string]interface{})

	if err := decoder.Decode(&body); err != nil {
		http.Error(w, "Failed to decode body", http.StatusBadRequest)
	}
	posts, _ := questionImp.QueryQuestions(body["page_number"].(int))
	res.Code = 1
	res.Data = posts
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, "Failed to encode posts", http.StatusInternalServerError)
		return
	}

}
