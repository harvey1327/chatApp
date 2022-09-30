package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CollectionCommands interface {
	FindByID(id string) (interface{}, error)
	FindSingleByQuery(query findBy) (interface{}, error)
	FindMultipleByQuery(query findBy) ([]interface{}, error)
	InsertOne(object interface{}) error
}

type mongoDBCollectionImpl struct {
	database   *mongo.Database
	collection *mongo.Collection
}

func NewCollection(database DB, collection string) CollectionCommands {
	return &mongoDBCollectionImpl{
		database:   database.getDatabase(),
		collection: database.getDatabase().Collection(collection),
	}
}

func (m *mongoDBCollectionImpl) InsertOne(object interface{}) error {
	_, err := m.collection.InsertOne(context.TODO(), object)
	if err != nil {
		return err
	}
	return nil
}

func (m *mongoDBCollectionImpl) FindByID(id string) (interface{}, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var result interface{}
	err = m.collection.FindOne(context.TODO(), bson.M{"_id": objectId}).Decode(&result)
	if err != nil {
		return nil, err
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

func (m *mongoDBCollectionImpl) FindSingleByQuery(query findBy) (interface{}, error) {
	var result interface{}
	err := m.collection.FindOne(context.TODO(), query.convert()).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *mongoDBCollectionImpl) FindMultipleByQuery(query findBy) ([]interface{}, error) {
	results := make([]interface{}, 0)
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
