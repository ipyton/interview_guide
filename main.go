package main

import (
	"context"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
	"strings"
	"time"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/service"
)

func AuthMiddleware(next http.Handler) http.Handler {
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
		valid, err := impl.IsTokenValid(token, r.RequestURI)
		r.Header.Set("openid", valid.Openid)
		if err != nil && valid.ExpiresAt.Before(time.Now()) {
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
	mux.Handle("/home", AuthMiddleware(http.HandlerFunc(service.HomeHandler)))
	mux.Handle("/home/get", AuthMiddleware(http.HandlerFunc(service.GetHottestHandler)))
	mux.Handle("/questions/get", AuthMiddleware(http.HandlerFunc(service.GetQuestionsByPageHandler)))
	mux.Handle("/questions/getById", AuthMiddleware(http.HandlerFunc(service.GetQuestionsByIdHandler)))
	mux.Handle("/questions/upsert", AuthMiddleware(http.HandlerFunc(service.UpsertQuestions)))
	mux.Handle("/collections/collection/get_items_by_time", AuthMiddleware(http.HandlerFunc(service.GetCollectionItemsByTime)))
	mux.Handle("/collections/items/delete", AuthMiddleware(http.HandlerFunc(service.DeleteBookmarkItem)))
	mux.Handle("/collections/item/add", AuthMiddleware(http.HandlerFunc(service.AddBookmarkItem)))
	mux.Handle("/collections/collection/get_items_by_collection", AuthMiddleware(http.HandlerFunc(service.GetBookmarkItems)))
	mux.Handle("/collections/collection/get_items_by_category", AuthMiddleware(http.HandlerFunc(service.GetCollectionItemsByCategory)))
	mux.Handle("/collections/collection/add", AuthMiddleware(http.HandlerFunc(service.AddBookmarkCollection)))
	mux.Handle("/collections/collection/delete", AuthMiddleware(http.HandlerFunc(service.DelBookmarkCollection)))
	mux.Handle("/collections/collection/get", AuthMiddleware(http.HandlerFunc(service.GetBookmarkCollections)))
	mux.Handle("/classes/upsert", AuthMiddleware(http.HandlerFunc(service.UpsertClassHandler)))
	mux.Handle("/classes/delete", AuthMiddleware(http.HandlerFunc(service.DeleteClassHandler)))
	mux.Handle("/classes/get", AuthMiddleware(http.HandlerFunc(service.GetClassHandler)))
	mux.Handle("/user/login", http.HandlerFunc(service.LoginHandler)) // 登录不需要身份验证
	//mux.Handle("/user/validate", http.HandlerFunc(service.ValidateTokenHandler)) // 验证令牌不需要身份验证
	mux.Handle("/user/get_info", AuthMiddleware(http.HandlerFunc(service.GetUserInfo)))
	mux.Handle("/user/set_class", AuthMiddleware(http.HandlerFunc(service.SetClass)))

	cors := handlers.CORS(handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"get", "post", "put", "patch", "delete"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}), // Add other headers as needed
	)

	log.Fatal(http.ListenAndServe(":5050", cors(mux)))
	var err error
	if err = db.MongoClient.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}
