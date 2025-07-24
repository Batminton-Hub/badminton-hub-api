package mongodb

import (
	"Badminton-Hub/internal/core/domain"
	"context"
	"time"
)

func (db *MongoDB) RegisterMember(member domain.Member) error {
	// Implementation for registering a member in MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	collection := db.Database.Collection("members")
	_, err := collection.InsertOne(ctx, member)
	if err != nil {
		return err
	}
	return nil
}
