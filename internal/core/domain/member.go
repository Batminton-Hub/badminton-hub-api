package domain

import (
	"time"
)

// Member Structure
type Member struct {
	UserID       string    `json:"user_id" bson:"user_id"`                             // Unique user ID
	Username     string    `json:"username" bson:"username"`                           // Unique username
	DisplayName  string    `json:"display_name" bson:"display_name"`                   // Display name
	Password     string    `json:"password,omitempty" bson:"password,omitempty"`       // Password hash
	Email        string    `json:"email" bson:"email"`                                 // Unique email
	Phone        string    `json:"phone" bson:"phone"`                                 // Unique phone number
	Hash         string    `json:"hash,omitempty" bson:"hash,omitempty"`               // Unique hash for password reset or verification
	Status       string    `json:"status" bson:"status"`                               // ACTIVE, BANNED, DELETED
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`                       // Creation timestamp
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`                       // Last update timestamp
	Tag          []string  `json:"tag" bson:"tag"`                                     // Tags for categorization
	MainTag      []string  `json:"main_tag" bson:"main_tag"`                           // Main tag for categorization
	Gender       string    `json:"gender" bson:"gender"`                               // Gender
	ProfileImage string    `json:"profile_image" bson:"profile_image"`                 // URL to profile image
	DateOfBirth  string    `json:"date_of_birth" bson:"date_of_birth"`                 // Date of birth in YYYY-MM-DD format
	TypeMember   string    `json:"type_member,omitempty" bson:"type_member,omitempty"` // Type member
	Region       string    `json:"region" bson:"region"`                               // Region or country
	Permission   []string  `json:"permission,omitempty" bson:"permission,omitempty"`   // Permission
	GoogleID     string    `json:"google_id" bson:"google_id"`                         // Google ID
	Address      Address   `json:"address,omitempty" bson:"address,omitempty"`         // Address details
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
	Password string `json:"password" binding:"required,min=8"` // password
	Gender   string `json:"gender" binding:"required"`         // เพศ
}

// Login Structure
type LoginForm struct {
	Email      string `json:"email" binding:"required,email"`    // email
	Password   string `json:"password" binding:"required,min=6"` // password
	Platform   string `json:"platform"`                          // platform
	PlatformID string `json:"platform_id"`                       // platform id
}

// Profile Request and Response
type ReqGetProfile struct {
	UserID string
}
type RespGetProfile struct {
	Member Member `json:"member,omitempty"`
	Resp   Resp   `json:"resp,omitempty"`
}

type ReqUpdateProfile struct {
	DisplayName  string   `json:"display_name,omitempty" bson:"display_name,omitempty"`
	ProfileImage string   `json:"profile_image,omitempty" bson:"profile_image,omitempty"`
	DateOfBirth  string   `json:"date_of_birth,omitempty" bson:"date_of_birth,omitempty" binding:"omitempty,datetime=2006-01-02"`
	Region       string   `json:"region,omitempty" bson:"region,omitempty"`
	Gender       string   `json:"gender,omitempty" bson:"gender,omitempty"`
	Phone        string   `json:"phone,omitempty" bson:"phone,omitempty" binding:"omitempty,numeric,min=10,max=10"`
	Tag          []string `json:"tag,omitempty" bson:"tag,omitempty"`
	Status       string   `json:"status,omitempty" bson:"status,omitempty"`
}

type RespUpdateProfile struct {
	Member Member `json:"member,omitempty"`
	Resp   Resp   `json:"resp,omitempty"`
}
