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

type AuthResponse struct {
	AuthBody AuthBody `json:"auth_body,omitempty"`
	Code     int      `json:"code,omitempty"`
	Message  string   `json:"message,omitempty"`
}

type AuthBody struct {
	CreateAt   time.Time  `json:"create_at"`
	Exp        int64      `json:"exp"`
	Data       AuthMember `json:"data"`
	Permission []string   `json:"permission"`
}
