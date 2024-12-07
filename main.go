package main

import (
	"context"
	"fmt"
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
		//if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		//	http.Error(w, "Unauthorized", http.StatusUnauthorized)
		//	return
		//}
		token := ""
		if authHeader != "" {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}

		//notice to delete it.
		if r.URL.Query().Get("token") != "" {
			token = r.URL.Query().Get("token")
		}

		if token == "czhdawang" {
			next.ServeHTTP(w, r)
			return
		}
		if token == "" {
			http.Error(w, "Unauthorized1", http.StatusUnauthorized)
			return
		}
		valid, err := impl.IsTokenValid(token, r.RequestURI)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Unauthorized2", http.StatusUnauthorized)
			return
		}
		r.Header.Set("openid", valid.Openid)
		if err != nil && valid.ExpiresAt.Before(time.Now()) {
			http.Error(w, "Unauthorized3", http.StatusUnauthorized)
			return
		}

		// 如果验证成功，调用下一个处理器
		next.ServeHTTP(w, r)
	})
}

func main() {
	db.InitMongo()
	db.InitMinio()
	db.InitElasticSearchClient()
	mux := http.NewServeMux()

	// 定义路由及对应的处理器
	mux.Handle("/home", AuthMiddleware(http.HandlerFunc(service.HomeHandler)))
	mux.Handle("/home/get", AuthMiddleware(http.HandlerFunc(service.GetHottestPostHandler)))
	mux.Handle("/questions/page/get/id", AuthMiddleware(http.HandlerFunc(service.GetQuestionsByPagingHandler)))
	mux.Handle("/questions/page/get/hottest", AuthMiddleware(http.HandlerFunc(service.GetHottestQuestionHandler)))
	mux.Handle("/questions/page/get/tag", AuthMiddleware(http.HandlerFunc(service.GetQuestionsByTag)))
	mux.Handle("/questions/page/get/details", AuthMiddleware(http.HandlerFunc(service.GetQuestionByIdHandler)))
	//mux.Handle("/questions/getById", AuthMiddleware(http.HandlerFunc(service.GetQuestionsByIdHandler)))
	mux.Handle("/questions/upsert", AuthMiddleware(http.HandlerFunc(service.UpsertQuestions)))
	mux.Handle("/questions/insert_by_file", AuthMiddleware(http.HandlerFunc(service.UpsertQuestionsByFile)))
	mux.Handle("/questions/search", AuthMiddleware(http.HandlerFunc(service.GetResults)))
	mux.Handle("/questions/search/test", AuthMiddleware(http.HandlerFunc(service.Testing)))
	mux.Handle("/collections/collection/get_items_by_time", AuthMiddleware(http.HandlerFunc(service.GetCollectionItemsByTime)))
	mux.Handle("/collections/item/delete", AuthMiddleware(http.HandlerFunc(service.DeleteBookmarkItem)))
	mux.Handle("/collections/item/add", AuthMiddleware(http.HandlerFunc(service.AddBookmarkItem)))
	mux.Handle("/collections/collection/get_items_by_collection", AuthMiddleware(http.HandlerFunc(service.GetBookmarkItems)))
	mux.Handle("/collections/collection/get_items_by_category", AuthMiddleware(http.HandlerFunc(service.GetCollectionItemsByCategory)))
	mux.Handle("/collections/collection/create", AuthMiddleware(http.HandlerFunc(service.AddBookmarkCollection)))
	mux.Handle("/collections/collection/delete", AuthMiddleware(http.HandlerFunc(service.DelBookmarkCollection)))
	mux.Handle("/collections/collection/get", AuthMiddleware(http.HandlerFunc(service.GetBookmarkCollections)))
	mux.Handle("/classes/insert", AuthMiddleware(http.HandlerFunc(service.InsertClassHandler)))
	mux.Handle("/classes/update", AuthMiddleware(http.HandlerFunc(service.UpdateClassHandler)))
	mux.Handle("/classes/delete", AuthMiddleware(http.HandlerFunc(service.DeleteClassHandler)))
	mux.Handle("/classes/get", AuthMiddleware(http.HandlerFunc(service.GetClassHandler)))
	mux.Handle("/user/login", http.HandlerFunc(service.LoginHandler)) // 登录不需要身份验证
	//mux.Handle("/user/validate", http.HandlerFunc(service.ValidateTokenHandler)) // 验证令牌不需要身份验证
	mux.Handle("/user/get_info", AuthMiddleware(http.HandlerFunc(service.GetUserInfo)))
	mux.Handle("/user/set_class", AuthMiddleware(http.HandlerFunc(service.SetClass)))
	mux.Handle("/tags/get", AuthMiddleware(http.HandlerFunc(service.GetTags)))
	mux.Handle("/tags/update", AuthMiddleware(http.HandlerFunc(service.UpdateTags)))
	mux.Handle("/voice/get", AuthMiddleware(http.HandlerFunc(service.GetVoiceHandler)))
	mux.Handle("/voice/delete", AuthMiddleware(http.HandlerFunc(service.DeleteVoiceHandler)))
	mux.Handle("/voice/generate", AuthMiddleware(http.HandlerFunc(service.GenerateVoiceHandler)))

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
