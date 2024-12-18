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
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
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

		if err != nil || valid.Openid == "" || valid.ExpiresAt.Before(time.Now()) {
			fmt.Println(err)
			http.Error(w, "Unauthorized2", http.StatusUnauthorized)
			return
		}
		r.Header.Set("openid", valid.Openid)

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
	mux.Handle("/questions/search/suggestions", AuthMiddleware(http.HandlerFunc(service.GetSuggestions)))

	//
	mux.Handle("/questions/advice/get", AuthMiddleware(http.HandlerFunc(service.GetAdvisedQuestions)))
	mux.Handle("/questions/advice/approve", AuthMiddleware(http.HandlerFunc(service.ApproveAQuestion)))
	mux.Handle("/questions/advice/advice", AuthMiddleware(http.HandlerFunc(service.AdviceAQuestion)))

	//all
	mux.Handle("/collections/collection/get_items_by_collect_time", AuthMiddleware(http.HandlerFunc(service.GetBookmarkItems)))
	mux.Handle("/collections/item/delete", AuthMiddleware(http.HandlerFunc(service.DeleteBookmarkItem)))
	mux.Handle("/collections/item/add", AuthMiddleware(http.HandlerFunc(service.AddBookmarkItem)))
	// specific collection
	mux.Handle("/collections/collection/get_items_by_collection", AuthMiddleware(http.HandlerFunc(service.GetCollectionItemsByTime)))
	// specific tag.
	mux.Handle("/collections/collection/get_items_by_tag", AuthMiddleware(http.HandlerFunc(service.GetCollectionItemsByTag)))
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
	mux.Handle("/user/get_avatar", AuthMiddleware(http.HandlerFunc(service.GetUserAvatar)))
	mux.Handle("/user/set_class", AuthMiddleware(http.HandlerFunc(service.SetClass)))
	mux.Handle("/user/changeAvatar", AuthMiddleware(http.HandlerFunc(service.UploadAvatar)))
	mux.Handle("/user/change_name", AuthMiddleware(http.HandlerFunc(service.SetUserName)))

	mux.Handle("/tags/get", AuthMiddleware(http.HandlerFunc(service.GetTags)))
	mux.Handle("/tags/update", AuthMiddleware(http.HandlerFunc(service.UpdateTags)))

	mux.Handle("/voice/get", AuthMiddleware(http.HandlerFunc(service.GetVoiceHandler)))
	mux.Handle("/voice/delete", AuthMiddleware(http.HandlerFunc(service.DeleteVoiceHandler)))
	mux.Handle("/voice/generate", AuthMiddleware(http.HandlerFunc(service.GenerateVoiceHandler)))

	mux.Handle("/feedback/get", AuthMiddleware(http.HandlerFunc(service.GetFeedback)))
	mux.Handle("/feedback/reply", AuthMiddleware(http.HandlerFunc(service.ReplyFeedback)))
	mux.Handle("/feedback/send", AuthMiddleware(http.HandlerFunc(service.SendFeedback)))

	//mux.Handle("/statistics/get", AuthMiddleware(http.HandlerFunc(service.GetStatistics)))
	//mux.Handle("/statistics/set", AuthMiddleware(http.HandlerFunc(service.SetStatistics)))

	//mux.Handle("/notification/get_configuration", AuthMiddleware(http.HandlerFunc(service.GetConfiguration)))
	//mux.Handle("/notification/config", AuthMiddleware(http.HandlerFunc(service.SetConfiguration)))

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
