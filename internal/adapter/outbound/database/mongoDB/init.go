package mongoDB

import (
	"Badminton-Hub/util"
	"context"
	"log"
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

func NewMongoDB(dbName string) *MongoDB {
	config := util.LoadConfig()
	client, err := mongo.Connect(options.Client().ApplyURI(config.MongoDBURL))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
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
	if err := db.Client.Disconnect(db.Ctx); err != nil {
		panic(err)
	}
	db.Cancel()
}
