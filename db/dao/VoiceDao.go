package dao

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"os"
	"time"
	"wxcloudrun-golang/db"
)

type VoiceDaoImpl struct {
	VoiceInterface
}

var questionImpl QuestionInterfaceImpl

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
	BearerToken := os.Getenv("tss_token")
	reqID := uuid.New().String()
	params := make(map[string]map[string]interface{})
	params["app"] = make(map[string]interface{})
	//填写平台申请的appid
	params["app"]["appid"] = "9275620043"
	//这部分的token不生效，填写下方的默认值就好
	params["app"]["token"] = "access_token"
	//填写平台上显示的集群名称
	params["app"]["cluster"] = "volcano_tts"
	params["user"] = make(map[string]interface{})
	//这部分如有需要，可以传递用户真实的ID，方便问题定位
	params["user"]["uid"] = "uid"
	params["audio"] = make(map[string]interface{})
	//填写选中的音色代号
	params["audio"]["voice_type"] = "xxxx"
	params["audio"]["encoding"] = "wav"
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
	fmt.Printf("resp body:%s\n", synResp)
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

func (*VoiceDaoImpl) GenerateVoice(questionId int64) ([]byte, error) {
	question, err := questionImpl.GetQuestionById(questionId)
	if err != nil {
		return []byte{}, err
	}
	details := question.Details
	title := question.Title
	voice, err := synthesis(details + title)
	fmt.Println(voice)
	return voice, err
}

func (*VoiceDaoImpl) GetVoice(questionId string) error {
	var collection = db.MongoClient.Database("interview_guide").Collection("question")
	collection.FindOne()
	return nil
}
