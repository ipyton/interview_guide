package service

import (
	"fmt"
	"net/http"
	"time"
	"wxcloudrun-golang/db/model"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}

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

	// 设置响应的内容类型为JSON
	w.Header().Set("Content-Type", "application/json")

	// 生成10条记录
	posts := generateHotPosts()

	// 将生成的记录编码为JSON并写入响应
	err := json.NewEncoder(w).Encode(posts)
	if err != nil {
		http.Error(w, "Failed to encode posts", http.StatusInternalServerError)
		return
	}
}
