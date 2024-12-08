package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"mime/multipart"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

type UserInformationDaoImpl struct {
	UserInformationInterface
}

var storage FileManagerImpl

func (UserInformationDaoImpl) ChangeMembershipStatus(openid string, status bool) error {
	return nil
}

func (UserInformationDaoImpl) AddPoints(openid string, points int) error {
	return nil
}

func (UserInformationDaoImpl) UploadAvatar(openid string, file multipart.File) {
	// Parse the form data, which allows us to access the uploaded file

}

func (UserInformationDaoImpl) UpdateUserInfo(user model.User) error {
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"openid": user.OpenId}
	update := bson.D{
		{"$set", bson.D{
			{"open_id", user.OpenId},
			{"username", user.Username},
			{"avatar_url", user.AvatarURL},
			{"email", user.Email},
			{"phone_number", user.PhoneNumber},
		}},
	}
	_, err := db.MongoClient.Database("interview_guide").Collection("user_status").UpdateOne(context.TODO(), filter, update, opts)
	return err
}

func (UserInformationDaoImpl) GetUserInfo(openid string) (model.User, error) {
	filter := bson.M{"openid": openid}
	user := model.User{}
	err := db.MongoClient.Database("interview_guide").Collection("user_info").FindOne(context.TODO(), filter).Decode(&user)
	storage.GetFile(openid, "avatar")
	return user, err
}

func (UserInformationDaoImpl) UpdateUserName(openid string, username string) error {
	filter := bson.M{"openid": openid}
	update := bson.M{"$set": bson.M{"username": username}}
	db.MongoClient.Database("interview_guide").Collection("user_info").FindOneAndUpdate(context.TODO(), filter, update)
	return nil
}

func (UserInformationDaoImpl) UpdateUserClass(openid string, classId int, className string) error {
	opts := options.Update()
	filter := bson.M{"openid": openid}
	update := bson.D{
		{"$set", bson.D{
			{"class_id", classId},
			{"class_name", className},
		}},
	}
	_, err := db.MongoClient.Database("interview_guide").Collection("user_info").UpdateOne(context.TODO(), filter, update, opts)
	return err
}
