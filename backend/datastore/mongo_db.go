// MongoDB client implementing url.RecordStore interface

package datastore

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongoDb(url string) (*mongo.Client, error) {

	_ = url // TODO: add to config

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return nil, err
	}

	// Check the connection
	if err = client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	}

	log.Println("MongoClient connected")

	return client, nil
}

// func Disconnect() {
// 	err := client.Disconnect(context.TODO())
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Println("Disconnected from MongoDB!")
// }
