package domain

import "time"

type AuthMember struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	// UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	HashAuth  string    `json:"hash"`
}
