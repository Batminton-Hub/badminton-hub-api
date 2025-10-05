package mongoDB

import (
	"Badminton-Hub/internal/core/domain"
	"context"
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	MemberCollection = "member"
)

func (db *MongoDB) SaveMember(ctx context.Context, member domain.Member) domain.ErrInfo {
	result := &mongo.InsertOneResult{}
	errResult := domain.ErrInfo{}
	collection := db.Database.Collection(MemberCollection)
	result, errResult.Err = collection.InsertOne(ctx, member)
	if errResult.Err != nil {
		if strings.Contains(errResult.Err.Error(), "index: email_1 dup key") {
			errResult.Resp = domain.ErrMemberRegisterFailDuplicateEmail
		} else if strings.Contains(errResult.Err.Error(), "index: hash_1 dup key") {
			errResult.Resp = domain.ErrMemberRegisterFailDuplicateHash
		} else {
			errResult.Resp = domain.ErrCreateMemberFail
		}
	}

	if result != nil && result.InsertedID == nil {
		errResult.Resp = domain.ErrCreateMemberFail
	}

	return errResult
}

func (db *MongoDB) FindEmailMember(ctx context.Context, email string) (domain.Member, domain.ErrInfo) {
	errInfo := domain.ErrInfo{}
	collection := db.Database.Collection(MemberCollection)
	member := domain.Member{}
	filter := bson.M{
		"email": email,
	}
	errInfo.Err = collection.FindOne(ctx, filter).Decode(&member)
	if errInfo.Err != nil {
		if errInfo.Err == mongo.ErrNoDocuments {
			errInfo.Resp = domain.ErrMemberEmailNotFound
		}
		errInfo.Resp = domain.ErrMemberEmailNotFound
	}

	return member, errInfo
}

func (db *MongoDB) GetMemberByUserID(ctx context.Context, userID string) (domain.Member, domain.ErrInfo) {
	errInfo := domain.ErrInfo{}
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
	errInfo.Err = collection.FindOne(ctx, filter, option).Decode(&member)
	if errInfo.Err != nil {
		errInfo.Resp = domain.ErrGetMember
	}
	return member, errInfo
}

func (db *MongoDB) UpdateMember(ctx context.Context, userID string, request domain.ReqUpdateProfile) domain.ErrInfo {
	result := &mongo.UpdateResult{}
	errInfo := domain.ErrInfo{}
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
	result, errInfo.Err = collection.UpdateOne(ctx, filter, update)
	if errInfo.Err != nil {
		errInfo.Resp = domain.ErrUpdateMemberFail
		return errInfo
	}
	if result.MatchedCount == 0 {
		errInfo.Resp = domain.ErrUpdateMemberFail
		errInfo.Err = errors.New("matched count is 0")
		return errInfo
	}
	return errInfo
}
