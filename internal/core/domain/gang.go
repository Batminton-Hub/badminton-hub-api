package domain

import "time"

type Gang struct {
	GangID    string       `json:"gang_id" bson:"gang_id"`
	OwnerID   string       `json:"owner_id" bson:"owner_id"` // User ID ของผู้สร้าง
	Title     string       `json:"title" bson:"title"`
	CreatedAt time.Time    `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time    `json:"gang_updated_at" bson:"gang_updated_at"`
	Members   []GangMember `json:"members" bson:"members"`
	Tag       []string     `json:"tag" bson:"tag"`
	MainTag   []string     `json:"main_tag" bson:"main_tag"`
	Address   Address      `json:"address" bson:"address"`
}

type GangMember struct {
	UserID     string    `json:"user_id" bson:"user_id"`
	Username   string    `json:"username" bson:"username"`
	JoinAt     time.Time `json:"join_at" bson:"join_at"`
	Permission string    `json:"permission" bson:"permission"` // ADMIN, MEMBER
}

type RespCreateGang struct {
	GangID string `json:"gang_id"`
	Resp   Resp   `json:"resp"`
}
