package main

import (
	"context"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
	"strings"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/service"
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		impl := dao.UserStatusDaoImpl{}

		// 从请求头中提取 Authorization 字段
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// 提取并验证令牌
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if !impl.IsTokenValid(token, r.RequestURI) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// 如果验证成功，调用下一个处理器
		next.ServeHTTP(w, r)
	})
}

func main() {
	db.InitMongo()
	mux := http.NewServeMux()

	// 定义路由及对应的处理器
	mux.Handle("/home", http.HandlerFunc(service.HomeHandler))
	mux.Handle("/home/get", http.HandlerFunc(service.GetHottestHandler))
	mux.Handle("/questions/get", http.HandlerFunc(service.GetQuestionsByPageHandler))
	mux.Handle("/questions/getById", http.HandlerFunc(service.GetQuestionsByIdHandler))
	mux.Handle("/questions/upsert", http.HandlerFunc(service.UpsertQuestions))
	mux.Handle("/collections/items/get", http.HandlerFunc(service.GetBookmarkItems))
	mux.Handle("/collections/items/delete", http.HandlerFunc(service.DeleteBookmarkItem))
	mux.Handle("/collections/item/add", http.HandlerFunc(service.AddBookmarkItem))
	mux.Handle("/collections/collection/add", http.HandlerFunc(service.AddBookmarkCollection))
	mux.Handle("/collections/collection/delete", http.HandlerFunc(service.DelBookmarkCollection))
	mux.Handle("/collections/collection/get", http.HandlerFunc(service.GetBookmarkCollections))
	mux.Handle("/classes/upsert", http.HandlerFunc(service.UpsertClassHandler))
	mux.Handle("/classes/delete", http.HandlerFunc(service.DeleteClassHandler))
	mux.Handle("/classes/get", http.HandlerFunc(service.GetClassHandler))
	mux.Handle("/user/login", http.HandlerFunc(service.LoginHandler))
	//mux.Handle("/user/validate", http.HandlerFunc(service.ValidateTokenHandler))
	mux.Handle("/user/logout", http.HandlerFunc(service.LogoutHandler))

	cors := handlers.CORS(handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"get", "post", "put", "patch", "delete"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}), // Add other headers as needed
	)

	log.Fatal(http.ListenAndServe(":5050", cors(http.DefaultServeMux)))
	var err error
	if err = db.MongoClient.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}
