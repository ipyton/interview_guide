package model

import "time"

type User struct {
	// 基本信息
	OpenId       string    `json:"open_id" bson:"open_id"`             // 用户ID
	Username     string    `json:"username" bson:"username"`           // 用户名
	AvatarURL    string    `json:"avatar_url" bson:"avatar_url"`       // 头像URL
	Email        string    `json:"email" bson:"email"`                 // 邮箱
	PhoneNumber  string    `json:"phone_number" bson:"phone_number"`   // 手机号
	RegisterDate time.Time `json:"register_date" bson:"register_date"` // 注册日期
	// 账户信息
	MembershipStatus string `json:"membership_status" bson:"membership_status"` // 会员状态
	Points           int    `json:"points" bson:"points"`                       // 积分
	ContinuousLogin  int    `json:"continuous_login" bson:"continuous_login"`   // 连续登录
}

type UserStatus struct {
	OpenId       string    `json:"open_id" bson:"open_id"`
	ThirdSession string    `json:"third_session" bson:"third_session"`
	SessionKey   string    `json:"session_key" bson:"session_key"`
	LastLogin    time.Time `json:"last_login" bson:"last_login"`       // 上次登录时间
	LastLoginIP  string    `json:"last_login_ip" bson:"last_login_ip"` // 上次登录IP
}
