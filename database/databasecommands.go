package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CollectionCommands[T any] interface {
	FindByID(id string) (T, error)
	FindSingleByQuery(query findBy) (T, error)
	FindMultipleByQuery(query findBy) ([]T, error)
	InsertOne(object T) error
}

var EMPTY = mongo.ErrNoDocuments

type mongoDBCollectionImpl[T any] struct {
	database   *mongo.Database
	collection *mongo.Collection
}

func NewCollection[T any](database DB, collection string) CollectionCommands[T] {
	return &mongoDBCollectionImpl[T]{
		database:   database.getDatabase(),
		collection: database.getDatabase().Collection(collection),
	}
}

func (m *mongoDBCollectionImpl[T]) InsertOne(object T) error {
	_, err := m.collection.InsertOne(context.TODO(), object)
	if err != nil {
		return err
	}
	return nil
}

func (m *mongoDBCollectionImpl[T]) FindByID(id string) (T, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return *new(T), err
	}
	var result T
	err = m.collection.FindOne(context.TODO(), bson.M{"_id": objectId}).Decode(&result)
	if err != nil {
		return *new(T), err
	}
	return result, nil
}

type findBy map[string]interface{}

func Query(field string, value interface{}) findBy {
	m := make(findBy)
	m.And(field, value)
	return m
}

func (fb findBy) And(field string, value interface{}) {
	fb[field] = value
}

func (fb findBy) convert() bson.M {
	res := bson.M{}
	for k, v := range fb {
		res[k] = v
	}
	return res
}

func (m *mongoDBCollectionImpl[T]) FindSingleByQuery(query findBy) (T, error) {
	var result T
	err := m.collection.FindOne(context.TODO(), query.convert()).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return *new(T), err
		}
		return *new(T), err
	}
	return result, nil
}

func (m *mongoDBCollectionImpl[T]) FindMultipleByQuery(query findBy) ([]T, error) {
	results := make([]T, 0)
	curr, err := m.collection.Find(context.TODO(), query.convert())
	if err != nil {
		return nil, err
	}
	err = curr.All(context.TODO(), &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}
