package service

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"io"
	"math/big"
	r "math/rand"
	"net/http"
	"os"
	"strings"
	"time"
	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
)

var userStatusDao = dao.UserStatusDaoImpl{}

type Request struct {
	Code string `json:"code"`
}
type Response struct {
	Token string `json:"token"`
}

type UserSession struct {
	SessionKey string `json:"session_key"` // 会话密钥
	UnionID    string `json:"unionid"`     // 用户在开放平台的唯一标识符
	ErrMsg     string `json:"errmsg"`      // 错误信息
	OpenID     string `json:"openid"`      // 用户唯一标识
	ErrCode    int32  `json:"errcode"`     // 错误码
}

type UserRecord struct {
	OpenID       string `json:"openid" bson:"openid"`
	ThirdSession string `json:"third_session" bson:"third_session"`
	SessionKey   string `json:"session_key" bson:"session_key"`
}

func generateName() string {
	// 修饰词数组
	adjectives := []string{"可爱", "调皮", "聪明", "勇敢", "慵懒", "机灵", "快乐", "胆小", "温柔", "勤劳"}

	// 动物名称数组
	animals := []string{"小猫", "小狗", "小兔", "小熊", "小鸟", "小鱼", "小鹿", "小猴", "小猪", "小羊"}

	// 地点数组
	locations := []string{"森林", "湖边", "草地", "山顶", "沙滩", "河边", "田野", "树下", "山洞", "农场"}

	// 活动数组
	activities := []string{"玩耍", "睡觉", "觅食", "打滚", "跳舞", "晒太阳", "打盹", "跑步", "游泳", "唱歌"}

	// 随机生成器种子
	r.New(r.NewSource(time.Now().UnixNano()))
	adjective := adjectives[r.Intn(len(adjectives))]
	animal := animals[r.Intn(len(animals))]
	location := locations[r.Intn(len(locations))]
	activity := activities[r.Intn(len(activities))]
	username := fmt.Sprintf("%s的%s在%s%s", adjective, animal, location, activity)
	return username

}

func generateRandomString(length int) (string, error) {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var result strings.Builder
	// 创建一个字符集的大小
	charSetLen := big.NewInt(int64(len(chars)))

	for i := 0; i < length; i++ {
		// 随机选择一个字符
		randomIndex, err := rand.Int(rand.Reader, charSetLen)
		if err != nil {
			return "", err
		}
		// 将随机字符附加到结果中
		result.WriteByte(chars[randomIndex.Int64()])
	}

	return result.String(), nil
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	var req Request
	err = json.Unmarshal(body, &req)

}

var jwtKey = []byte("your_secret_key") // 建议在配置文件或环境变量中管理密钥

func generateJWT(openid string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &dao.Claims{
		Openid: openid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	fmt.Println(body)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var request Request
	err = json.Unmarshal(body, &request)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	req, err := http.NewRequest("GET", "https://api.weixin.qq.com/sns/jscode2session", nil)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to create request", http.StatusBadRequest)
		return
	}
	q := req.URL.Query()
	q.Add("secret", os.Getenv("app_secret"))
	q.Add("js_code", request.Code)
	q.Add("appid", os.Getenv("appId"))
	q.Add("grant_type", "client_credential")
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to send request", http.StatusBadRequest)
		return
	}
	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	var response UserSession
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to parse JSON response", http.StatusBadRequest)
		return
	}

	randomString, err := generateJWT(response.OpenID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to generate random string", http.StatusInternalServerError)
		return
	}

	status := model.UserStatus{
		OpenId:       response.OpenID,
		ThirdSession: randomString,
		SessionKey:   response.SessionKey,
		LastLogin:    time.Now(),
		LastLoginIP:  r.RemoteAddr,
	}

	if userStatusDao.IsUserExists(response.OpenID) {
		err := userStatusDao.UpsertLoginStatus(status)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to update user status", http.StatusBadRequest)
			return
		}
	} else {
		user := model.User{OpenId: response.OpenID, Username: generateName()}
		err = userStatusDao.Registration(status, user)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to register user", http.StatusBadRequest)
			return
		}
	}

	// Send the JSON response
	resJson := JsonResult{Code: 1, Data: randomString}
	responseBody, err := json.Marshal(resJson)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to encode response JSON", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseBody)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}
