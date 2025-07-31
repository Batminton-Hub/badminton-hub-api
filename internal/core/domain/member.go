package domain

import (
	"time"
)

// Member Structure
type Member struct {
	Username     string    `json:"username"`      // Unique username
	DisplayName  string    `json:"display_name"`  // Display name
	Password     string    `json:"password"`      // Password hash
	Email        string    `json:"email"`         // Unique email
	Phone        string    `json:"phone"`         // Unique phone number
	Hash         string    `json:"hash"`          // Unique hash for password reset or verification
	Status       string    `json:"status"`        // ACTIVE, BANNED, DELETED
	CreatedAt    time.Time `json:"created_at"`    // Creation timestamp
	UpdatedAt    time.Time `json:"updated_at"`    // Last update timestamp
	Tag          []string  `json:"tag"`           // Tags for categorization
	MainTag      []string  `json:"main_tag"`      // Main tag for categorization
	Gender       string    `json:"gender"`        // Gender
	ProfileImage string    `json:"profile_image"` // URL to profile image
	DateOfBirth  string    `json:"date_of_birth"` // Date of birth in YYYY-MM-DD format
	Region       string    `json:"region"`        // Region or country
	// Address     Address   `json:"address"`      // Address details
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
	ErrorCode   int    `json:"error_code,omitempty"`
	Error       string `json:"error,omitempty"`
	Message     string `json:"message,omitempty"`
}

// Login Structure
type LoginForm struct {
	Email    string `json:"email" binding:"required,email"`    // email
	Password string `json:"password" binding:"required,min=6"` // password
	Clerk    string `json:"clerk"`                             // Clerk token ได้มาจากการ login ด้วย platform อื่น ๆ
}

type ResponseLogin struct {
	BearerToken string `json:"bearer_token,omitempty"`
	ErrorCode   int    `json:"error_code,omitempty"`
	Error       string `json:"error,omitempty"`
	Message     string `json:"message,omitempty"`
}
