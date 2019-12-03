package mongo

import (
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

func NewUniqueIndex(keys ...string) mongoDriver.IndexModel {
	unique := true
	bsonIndexDoc := bsonx.Doc{}

	for _, key := range keys {
		bsonIndexDoc = bsonIndexDoc.Append(key, bsonx.Int32(1))
	}

	return mongoDriver.IndexModel{
		Keys: bsonIndexDoc,
		Options: &mongoOptions.IndexOptions{
			Unique: &unique,
		},
	}
}
