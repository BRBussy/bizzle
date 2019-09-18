package mongo

import (
	"context"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
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

func (c *Collection) FindOne(document interface{}, identifier identifier.Identifier) error {
	if identifier == nil {
		return ErrInvalidIdentifier{Reasons: []string{"nil identifier"}}
	} else if err := identifier.IsValid(); err != nil {
		return ErrInvalidIdentifier{Reasons: []string{err.Error()}}
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if err := c.driverCollection.FindOne(ctx, identifier.ToFilter()).Decode(document); err != nil {
		switch err {
		case mongoDriver.ErrNoDocuments:
			return ErrNotFound{}
		default:
			log.Error().Err(err).Msg("find one")
			return err
		}
	}
	return nil
}
