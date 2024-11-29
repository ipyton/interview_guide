package dao

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

type VoiceDaoImpl struct {
	VoiceInterface
}

var questionImpl QuestionInterfaceImpl
var fileImpl FileManagerImpl

type TTSServResponse struct {
	ReqID     string `json:"reqid"`
	Code      int    `json:"code"`
	Message   string `json:"Message"`
	Operation string `json:"operation"`
	Sequence  int    `json:"sequence"`
	Data      string `json:"data"`
}

func httpPost(url string, headers map[string]string, body []byte,
	timeout time.Duration) ([]byte, error) {
	client := &http.Client{
		Timeout: timeout,
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	retBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return retBody, err
}

func synthesis(text string) ([]byte, error) {
	BearerToken := os.Getenv("tts_token")
	fmt.Println(BearerToken)
	reqID := uuid.New().String()
	fmt.Println(reqID)
	params := make(map[string]map[string]interface{})
	params["app"] = make(map[string]interface{})
	//填写平台申请的appid
	params["app"]["appid"] = "9275620043"
	//这部分的token不生效，填写下方的默认值就好
	params["app"]["token"] = "ERsSZqW2bGVh1ai2thWIVla8p0HzRenM"
	//填写平台上显示的集群名称
	params["app"]["cluster"] = "volcano_tts"
	params["user"] = make(map[string]interface{})
	//这部分如有需要，可以传递用户真实的ID，方便问题定位
	params["user"]["uid"] = "uid"
	params["audio"] = make(map[string]interface{})
	//填写选中的音色代号
	params["audio"]["voice_type"] = "BV001_streaming"
	params["audio"]["encoding"] = "mp3"
	params["audio"]["speed_ratio"] = 1.0
	params["audio"]["volume_ratio"] = 1.0
	params["audio"]["pitch_ratio"] = 1.0
	params["request"] = make(map[string]interface{})
	params["request"]["reqid"] = reqID
	params["request"]["text"] = text
	params["request"]["text_type"] = "plain"
	params["request"]["operation"] = "query"

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	//bearerToken为saas平台对应的接入认证中的Token
	headers["Authorization"] = fmt.Sprintf("Bearer;%s", BearerToken)

	// URL查看上方第四点: 4.并发合成接口(POST)
	url := "https://openspeech.bytedance.com/api/v1/tts"
	timeo := 30 * time.Second
	bodyStr, _ := json.Marshal(params)
	synResp, err := httpPost(url, headers,
		[]byte(bodyStr), timeo)
	if err != nil {
		fmt.Printf("http post fail [err:%s]\n", err.Error())
		return nil, err
	}
	//fmt.Printf("resp body:%s\n", synResp)
	var respJSON TTSServResponse
	err = json.Unmarshal(synResp, &respJSON)
	if err != nil {
		fmt.Printf("unmarshal response fail [err:%s]\n", err.Error())
		return nil, err
	}
	code := respJSON.Code
	if code != 3000 {
		fmt.Printf("code fail [code:%d]\n", code)
		return nil, errors.New("resp code fail")
	}

	audio, _ := base64.StdEncoding.DecodeString(respJSON.Data)
	return audio, nil
}

func (*VoiceDaoImpl) GenerateVoice(questionId string) ([]byte, error) {
	id, err := strconv.ParseInt(questionId, 10, 64)
	question, err := questionImpl.GetQuestionById(id)
	if err != nil {
		return []byte{}, err
	}
	details := question.Details
	title := question.Title
	voice, err := synthesis(title + details)
	collection := db.MongoClient.Database("interview_guide").Collection("question_voice")
	id, err = strconv.ParseInt(questionId, 10, 64)
	if err != nil {
		return []byte{}, err
	}

	voiceModel := model.VoiceModel{QuestionId: id, Path: questionId[0:2] + "/" + questionId}
	_, err = collection.InsertOne(context.TODO(), voiceModel)
	if err != nil {
		return nil, err
	}
	err = fileImpl.UploadFile(questionId, "question-voice", voice)
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err

	}
	return voice, err
}

func (*VoiceDaoImpl) GetVoice(questionId int64) (*minio.Object, error) {
	var collection = db.MongoClient.Database("interview_guide").Collection("question_voice")
	filter := bson.M{"QuestionId": questionId}

	one := collection.FindOne(context.TODO(), filter)
	var voice model.VoiceModel
	err := one.Decode(&voice)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	fmt.Println(voice.Path)

	file, err := fileImpl.GetFile(voice.Path, "question-voice")
	return file, err
}

func (*VoiceDaoImpl) DeleteVoice(questionId string) error {
	var collection = db.MongoClient.Database("interview_guide").Collection("question_voice")
	filter := bson.M{"question_id": questionId}
	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	err = fileImpl.DeleteFile(questionId, "question-voice")
	return err
}
