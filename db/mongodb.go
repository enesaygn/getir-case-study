package db

import (
	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connection is the global MongoDB connection
var (
	Connection *mongo.Database
	once       sync.Once
)

// MongoDB constants
const (
	ConnString = "mongodb+srv://challengeUser:WUMglwNBaydH8Yvu@challenge-xzwqd.mongodb.net/getir-case-study?retryWrites=true"
	DB         = "getircase-study"
)

// InitializeMongoDB initializes the MongoDB connection
func InitializeMongoDB() error {
	once.Do(func() {

		// Set the context with a timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Connect to the MongoDB server
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(ConnString))
		if err != nil {
			return
		}
		// Check the connection
		err = client.Ping(ctx, nil)
		if err != nil {
			return
		}
		// Set the global connection
		Connection = client.Database(DB)

		// Log the successful connection
		log.Println("Connected to MongoDB")
	})

	return nil
}
