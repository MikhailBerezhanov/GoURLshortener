// MongoDB client implementing url.RecordStore interface

package datastore

import (
	"context"
	"errors"
	"log"
	"time"
	"url_shortener/url"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client         *mongo.Client
	urlsCollection *mongo.Collection
}

func NewMongoDB() *MongoDB {
	return &MongoDB{}
}

func CreateContextWithTimeoutSec(seconds time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), seconds*time.Second)
}

func (m *MongoDB) Connect(uri string) error {
	log.Println("Connecting to MongoDB ...")

	_ = uri // TODO: add to config

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	ctx, cancel := CreateContextWithTimeoutSec(10)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		return err
	}

	m.client = client

	// Check the connection
	if err = m.client.Ping(ctx, nil); err != nil {
		return err
	}

	log.Println("Connected to MongoDB")

	m.urlsCollection = client.Database("urlRecords").Collection("urls")

	return nil
}

func (m *MongoDB) Disconnect() error {
	log.Println("Disconnecting from MongoDB ...")
	ctx, cancel := CreateContextWithTimeoutSec(5)
	defer cancel()

	err := m.client.Disconnect(ctx)
	if err != nil {
		return err
	}
	log.Println("Disconnected from MongoDB")
	return nil
}

// Implements url.RecordStore
func (m *MongoDB) InsertRecord(r *url.Record) error {
	if _, err := m.urlsCollection.InsertOne(context.TODO(), r); err != nil {
		return err
	}

	log.Println("Inserted record:", r)
	return nil
}

func (m *MongoDB) SelectRecord(shortURL string) (rec url.Record, err error) {
	ctx, cancel := CreateContextWithTimeoutSec(5)
	defer cancel()

	err = m.urlsCollection.FindOne(ctx, bson.M{"shortCode": shortURL}).Decode(&rec)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return url.Record{}, url.ErrRecordNotExist
		}
	}

	return
}
