package domain

import (
	"context"
	"time"
)

type AuthMember struct {
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	HashAuth  string    `json:"hash_auth"`
}

type HashAuth struct {
	CreateAt time.Time
	UserID   string
}

type AuthBody struct {
	CreateAt   time.Time  `json:"create_at"`
	Exp        int64      `json:"exp"`
	Data       AuthMember `json:"data"`
	Permission []string   `json:"permission"`
}

type BearerToken struct {
	Token string `json:"token"`
}

type LoginInfo struct {
	Context      context.Context
	TraceID      string
	SpanID       string
	ScopeName    string
	Path         string
	Job          string
	Platform     string
	PlatformData any
	LoginForm    LoginForm
	TypeSystem   string
}

type RegisterInfo struct {
	Platform     string
	PlatformData any
	RegisterForm RegisterForm
	TypeSystem   string
}

type AuthInfo struct {
	BearerToken BearerToken
	State       string
	Code        string
	Action      string
	Platform    string
	TypeSystem  string
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
