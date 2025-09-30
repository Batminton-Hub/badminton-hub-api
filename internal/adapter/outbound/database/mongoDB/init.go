package mongoDB

import (
	"Badminton-Hub/util"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type CloseMongoDB func()

type MongoDB struct {
	Database *mongo.Database
	Client   *mongo.Client
	Ctx      context.Context
}

func NewMongoDB() (*MongoDB, CloseMongoDB) {
	ctx, cancel := util.InitConText(2 * time.Second)
	defer cancel()
	config := util.LoadConfig()
	client, err := mongo.Connect(options.Client().ApplyURI(config.MongoDBURL))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	mongoDB := &MongoDB{
		Database: client.Database(config.DBName),
		Client:   client,
		Ctx:      ctx,
	}
	closeDB := closeMongoDB(ctx, client)
	return mongoDB, closeDB
}

func closeMongoDB(ctx context.Context, client *mongo.Client) CloseMongoDB {
	return func() {
		fmt.Println("closeMongoDB client")
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}
}
