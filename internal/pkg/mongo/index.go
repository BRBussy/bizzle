package mongo

import (
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

func NewUniqueIndex(key string) mongoDriver.IndexModel {
	unique := true

	return mongoDriver.IndexModel{
		Keys: bsonx.Doc{{Key: key, Value: bsonx.Int32(1)}},
		Options: &mongoOptions.IndexOptions{
			Unique: &unique,
		},
	}
}
