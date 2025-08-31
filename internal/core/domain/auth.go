package domain

import "time"

type AuthMember struct {
	Username  string    `json:"username"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	HashAuth  string    `json:"hash_auth"`
}

type AuthResponse struct {
	Code int   `json:"error_code,omitempty"`
	Err  error `json:"error,omitempty"`
}
