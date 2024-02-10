package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Connection *mongo.Database

const (
	ConnString = "mongodb+srv://challengeUser:WUMglwNBaydH8Yvu@challenge-xzwqd.mongodb.net/getir-case-study?retryWrites=true"
	DB         = "getircase-study"
)

func InitializeMongoDb() {
	Ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(Ctx, options.Client().ApplyURI(ConnString))
	if err != nil {
		log.Println(err)
	}
	session, err := client.StartSession()
	if err != nil {
		log.Println(err)
	}

	Connection = session.Client().Database(DB)

}
