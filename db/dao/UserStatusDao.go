package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

type UserStatusDaoImpl struct {
	UserStatusInterface
}

var jwtKey = []byte("your_secret_key") // 建议在配置文件或环境变量中管理密钥

type Claims struct {
	Openid string `json:"openid"`
	jwt.RegisteredClaims
}

func (UserStatusDaoImpl) UpsertLoginStatus(userStatus model.UserStatus) error {
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"openid": userStatus.OpenId}
	update := bson.D{
		{"$set", bson.D{
			{"openid", userStatus.OpenId},
			{"third_session", userStatus.ThirdSession},
			{"session_id", userStatus.SessionKey},
			{"last_login", time.Now()},
			{"last_login_ip", userStatus.LastLoginIP},
		}},
	}
	_, err := db.MongoClient.Database("interview_guide").Collection("user_status").UpdateOne(context.TODO(), filter, update, opts)
	return err
}
func (UserStatusDaoImpl) DeleteLoginStatus(userCode string) error {
	_, err := db.MongoClient.Database("interview_guide").Collection("user_status").DeleteOne(context.TODO(), bson.M{"openid": userCode})
	return err
}

func (UserStatusDaoImpl) LoginCount() {

}

func (UserStatusDaoImpl) IsUserExists(openid string) bool {
	one := db.MongoClient.Database("interview_guide").Collection("user_status").FindOne(context.TODO(), bson.M{"openid": openid})
	if errors.Is(one.Err(), mongo.ErrNoDocuments) {
		return false
	}
	return true
}
func verifyJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}

func (UserStatusDaoImpl) IsTokenValid(token string, requestPath string) (*Claims, error) {
	var canVisit = false
	if len(requestPath) != 0 {
		canVisit = true
	}
	if !canVisit {
		return nil, fmt.Errorf("Do not have access")
	}
	return verifyJWT(token)
}

func (UserStatusDaoImpl) Registration(userStatus model.UserStatus, user model.User) error {
	session, err := db.MongoClient.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(context.Background())

	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {
		collection1 := db.MongoClient.Database("interview_guide").Collection("user_status")
		collection2 := db.MongoClient.Database("interview_guide").Collection("user_info")

		// 插入订单
		_, err := collection1.InsertOne(sessCtx, userStatus)
		if err != nil {
			return nil, err
		}

		// 更新库存
		_, err = collection2.InsertOne(sessCtx, user)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
	_, err = session.WithTransaction(context.Background(), transactionFunc)
	if err != nil {
		return err
	}
	return err
}

func (UserStatusDaoImpl) CancelRegistration(user model.User) error {
	return nil
}
