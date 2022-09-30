package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB interface {
	Close() error
	getDatabase() *mongo.Database
}

type mongoDBImpl struct {
	database *mongo.Database
	client   *mongo.Client
}

func NewDB(database string) DB {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://0.0.0.0:27017").SetAuth(options.Credential{Username: "guest", Password: "guest"}))
	if err != nil {
		log.Fatal(err)
	}
	return &mongoDBImpl{
		database: client.Database(database),
		client:   client,
	}
}

func (m *mongoDBImpl) Close() error {
	return m.client.Disconnect(context.TODO())
}

func (m *mongoDBImpl) getDatabase() *mongo.Database {
	return m.database
}
