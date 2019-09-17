package mongo

import (
	"context"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Collection struct {
	*mongoDriver.Collection
}

func (c *Collection) CreateOne(document interface{}, collection string) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := c.InsertOne(ctx, document)
	return err
}
