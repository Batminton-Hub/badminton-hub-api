package domain

import (
	"time"
)

// Member Structure
type Member struct {
	UserID       string    `json:"user_id" bson:"user_id"`             // Unique user ID
	Username     string    `json:"username" bson:"username"`           // Unique username
	DisplayName  string    `json:"display_name" bson:"display_name"`   // Display name
	Password     string    `json:"password" bson:"password"`           // Password hash
	Email        string    `json:"email" bson:"email"`                 // Unique email
	Phone        string    `json:"phone" bson:"phone"`                 // Unique phone number
	Hash         string    `json:"hash" bson:"hash"`                   // Unique hash for password reset or verification
	Status       string    `json:"status" bson:"status"`               // ACTIVE, BANNED, DELETED
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`       // Creation timestamp
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`       // Last update timestamp
	Tag          []string  `json:"tag" bson:"tag"`                     // Tags for categorization
	MainTag      []string  `json:"main_tag" bson:"main_tag"`           // Main tag for categorization
	Gender       string    `json:"gender" bson:"gender"`               // Gender
	ProfileImage string    `json:"profile_image" bson:"profile_image"` // URL to profile image
	DateOfBirth  string    `json:"date_of_birth" bson:"date_of_birth"` // Date of birth in YYYY-MM-DD format
	Region       string    `json:"region" bson:"region"`               // Region or country
	Permission   []string  `json:"permission" bson:"permission"`       // Permission
	GoogleID     string    `json:"google_id" bson:"google_id"`         // Google ID
	// Address     Address   `json:"address"`      // Address details
}

type HashMember struct {
	Email    string
	Username string
}

type Address struct {
	Home    string `json:"home"`     // บ้าน
	Street  string `json:"street"`   // ถนน
	City    string `json:"city"`     // เมือง
	Country string `json:"country"`  // ประเทศ
	State   string `json:"state"`    // รัฐ
	ZipCode string `json:"zip_code"` // รหัสไปรษณีย์
}

// RegisterForm Structure
type RegisterForm struct {
	Email    string `json:"email" binding:"required,email"`    // email
	Password string `json:"password" binding:"required,min=6"` // password
	Gender   string `json:"gender" binding:"required"`         // เพศ
}

type ResponseRegisterMember struct {
	BearerToken string `json:"bearer_token,omitempty"`
	Code        int    `json:"code,omitempty"`
	Message     string `json:"message,omitempty"`
}

// Login Structure
type LoginForm struct {
	Email      string `json:"email" binding:"required,email"`    // email
	Password   string `json:"password" binding:"required,min=6"` // password
	Platform   string `json:"platform"`                          // platform
	PlatformID string `json:"platform_id"`                       // platform id
}

type ResponseLogin struct {
	BearerToken string `json:"bearer_token,omitempty"`
	Code        int    `json:"code,omitempty"`
	Message     string `json:"message,omitempty"`
}

type ResponseRedirectGoogleLogin struct {
	URL     string `json:"url,omitempty"`
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type ResponseRedirectGoogleRegister struct {
	URL     string `json:"url,omitempty"`
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type ResponseGoogleLoginCallback struct {
	UserInfo     GoogleUserInfo `json:"user_info,omitempty"`
	AccessToken  string         `json:"access_token,omitempty"`
	RefreshToken string         `json:"refresh_token,omitempty"`
	Code         int            `json:"code,omitempty"`
	Message      string         `json:"message,omitempty"`
}

type ResponseGoogleRegisterCallback struct {
	UserInfo     GoogleUserInfo `json:"user_info,omitempty"`
	AccessToken  string         `json:"access_token,omitempty"`
	RefreshToken string         `json:"refresh_token,omitempty"`
	Code         int            `json:"code,omitempty"`
	Message      string         `json:"message,omitempty"`
}

type GoogleUserInfo struct {
	Email         string `json:"email"`
	GivenName     string `json:"given_name"`
	ID            string `json:"id"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}
