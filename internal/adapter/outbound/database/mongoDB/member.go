package mongoDB

import (
	"Badminton-Hub/internal/core/domain"
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	MemberCollection = "member"
)

func (db *MongoDB) SaveMember(ctx context.Context, member domain.Member) error {
	collection := db.Database.Collection(MemberCollection)
	result, err := collection.InsertOne(ctx, member)
	if err != nil {
		if strings.Contains(err.Error(), "index: email_1 dup key") {
			return domain.ErrMemberRegisterFailDuplicateEmail.Err
		} else if strings.Contains(err.Error(), "index: hash_1 dup key") {
			return domain.ErrMemberRegisterFailDuplicateHash.Err
		} else {
			return err
		}
	}

	if result.InsertedID == nil {
		return domain.ErrCreateMemberFail.Err
	}

	return nil
}

func (db *MongoDB) FindEmailMember(ctx context.Context, email string) (domain.Member, error) {
	collection := db.Database.Collection(MemberCollection)
	member := domain.Member{}
	filter := bson.M{
		"email": email,
	}
	err := collection.FindOne(ctx, filter).Decode(&member)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return member, domain.ErrMemberEmailNotFound.Err
		}
		return member, err
	}

	return member, nil
}

func (db *MongoDB) GetMemberByUserID(ctx context.Context, userID string) (domain.Member, error) {
	collection := db.Database.Collection(MemberCollection)
	option := options.FindOne()
	project := bson.M{
		"google_id":  0,
		"hash":       0,
		"password":   0,
		"created_at": 0,
		"updated_at": 0,
	}
	option.SetProjection(project)
	member := domain.Member{}
	filter := bson.M{
		"user_id": userID,
	}
	err := collection.FindOne(ctx, filter, option).Decode(&member)
	if err != nil {
		return member, err
	}
	return member, nil
}

func (db *MongoDB) UpdateMember(ctx context.Context, userID string, request domain.ReqUpdateProfile) error {
	collection := db.Database.Collection(MemberCollection)

	filter := bson.M{
		"user_id": userID,
	}

	if request.Gender != "" {
		request.Status = domain.ACTIVE
	}

	update := bson.M{
		"$set": request,
	}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return domain.ErrUpdateMemberFail.Err
	}
	return nil
}
