package mongo

import (
	"context"
	"github.com/rs/zerolog/log"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Collection struct {
	driverCollection *mongoDriver.Collection
}

func (c *Collection) SetupIndex(model mongoDriver.IndexModel) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if _, err := c.driverCollection.Indexes().CreateOne(ctx, model); err != nil {
		log.Error().Err(err).Msg("setting up mongo collection index")
		return err
	}
	return nil
}

func (c *Collection) SetupIndices(models []mongoDriver.IndexModel) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if _, err := c.driverCollection.Indexes().CreateMany(ctx, models); err != nil {
		log.Error().Err(err).Msg("setting up mongo collection indices")
		return err
	}
	return nil
}

func (c *Collection) CreateOne(document interface{}) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := c.driverCollection.InsertOne(ctx, document)
	return err
}

func (c *Collection) FindOne(document interface{}, filter map[string]interface{}) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if err := c.driverCollection.FindOne(ctx, filter).Decode(document); err != nil {
		switch err {
		case mongoDriver.ErrNoDocuments:
		default:
			log.Error().Err(err).Msg("find one")
			return err
		}
	}
	return nil
}
