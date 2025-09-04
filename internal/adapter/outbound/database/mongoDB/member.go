package mongodb

import (
	"Badminton-Hub/internal/core/domain"
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (db *MongoDB) RegisterMember(ctx context.Context, member domain.Member) error {
	// Implementation for registering a member in MongoDB
	collection := db.Database.Collection("members")
	_, err := collection.InsertOne(ctx, member)
	if err != nil {
		if strings.Contains(err.Error(), "index: email_1 dup key") {
			return domain.ErrMemberRegisterFailDuplicateEmail.Err
		} else if strings.Contains(err.Error(), "index: hash_1 dup key") {
			return domain.ErrMemberRegisterFailDuplicateHash.Err
		} else {
			return err
		}
	}

	return nil
}

func (db *MongoDB) FindEmailMember(ctx context.Context, loginForm domain.LoginForm) (domain.Member, error) {
	collection := db.Database.Collection("members")
	member := domain.Member{}
	filter := bson.M{
		"email": loginForm.Email,
	}
	err := collection.FindOne(ctx, filter).Decode(&member)
	if err != nil {
		fmt.Println("Error finding member:", err)
		if err == mongo.ErrNoDocuments {
			return member, domain.ErrMemberEmailNotFound.Err
		}
		return member, err
	}

	return member, nil
}
