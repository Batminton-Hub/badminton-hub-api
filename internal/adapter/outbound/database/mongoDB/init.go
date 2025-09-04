package mongodb

import (
	"Badminton-Hub/util"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDB struct {
	Database *mongo.Database
	Client   *mongo.Client
	Ctx      context.Context
	Cancel   context.CancelFunc
}

// func NewMongoDB() port.DatabaseService {
func NewMongoDB(dbName string) *MongoDB {
	config, err := util.LoadConfig()
	if err != nil {
		panic(err)
	}

	client, err := mongo.Connect(options.Client().ApplyURI(config.MongoDBURL))
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	return &MongoDB{
		Database: client.Database(dbName),
		Client:   client,
		Ctx:      ctx,
		Cancel:   cancel,
	}
}

func (db *MongoDB) CloseDB() {
	defer db.Cancel()
	if err := db.Client.Disconnect(db.Ctx); err != nil {
		panic(err)
	}
}
