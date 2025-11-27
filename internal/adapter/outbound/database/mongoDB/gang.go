package mongoDB

import (
	"Badminton-Hub/internal/core/domain"
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	GangCollection = "gang"
)

func (db *MongoDB) SaveGang(ctx context.Context, gang domain.Gang) domain.ErrInfo {
	result := &mongo.InsertOneResult{}
	errResult := domain.ErrInfo{}
	collection := db.Database.Collection(GangCollection)
	result, errResult.Err = collection.InsertOne(ctx, gang)
	if errResult.Err != nil {
		errResult.Resp = domain.ErrCreateGangFail
	}

	if result != nil && result.InsertedID == nil {
		errResult.Resp = domain.ErrCreateGangFail
		errResult.Err = domain.ErrCreateGangFail.Err
	}

	return errResult
}
