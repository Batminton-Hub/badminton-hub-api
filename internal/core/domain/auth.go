package domain

import "time"

type AuthMember struct {
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	HashAuth  string    `json:"hash_auth"`
}

type HashAuth struct {
	CreateAt time.Time
	UserID   string
}

// type AuthResponse struct {
// 	AuthBody AuthBody `json:"auth_body,omitempty"`
// 	Code     int      `json:"code,omitempty"`
// 	Message  string   `json:"message,omitempty"`
// }

type AuthBody struct {
	CreateAt   time.Time  `json:"create_at"`
	Exp        int64      `json:"exp"`
	Data       AuthMember `json:"data"`
	Permission []string   `json:"permission"`
}

type BearerToken struct {
	Token string `json:"token"`
}

//---New Flow---//

type LoginInfo struct {
	Platform     string
	PlatformData any
	LoginForm    LoginForm
}

type RegisterInfo struct {
	Platform     string
	PlatformData any
	RegisterForm RegisterForm
}

type AuthInfo struct {
	BearerToken BearerToken
	State       string
	Code        string
	Action      string
	Platform    string
}

type RespAuth struct {
	AuthBody     AuthBody `json:"auth_body,omitempty"`
	PlatformData any      `json:"platform_data,omitempty"`
	Resp         Resp     `json:"resp,omitempty"`
}

type RespLogin struct {
	BearerToken string `json:"bearer_token,omitempty"`
	Resp        Resp   `json:"resp,omitempty"`
}

type RespRegister struct {
	BearerToken string `json:"bearer_token,omitempty"`
	Resp        Resp   `json:"resp,omitempty"`
}
