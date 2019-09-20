package mongo

import (
	"context"
	"github.com/BRBussy/bizzle/pkg/search/criteria"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	"github.com/BRBussy/bizzle/pkg/search/query"
	"github.com/rs/zerolog/log"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	mongoDriverOptions "go.mongodb.org/mongo-driver/mongo/options"
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

func (c *Collection) FindMany(documents interface{}, criteria criteria.Criteria, query query.Query) error {
	// build filter
	filter := make(map[string]interface{})
	if criteria != nil {
		reasonsInvalid := make([]string, 0)
		for i := range criteria {
			if err := criteria[i].IsValid(); err != nil {
				reasonsInvalid = append(reasonsInvalid, err.Error())
			}
		}
		if len(reasonsInvalid) > 0 {
			return ErrInvalidCriteria{Reasons: reasonsInvalid}
		}
		filter = criteria.ToFilter()
	}

	// build options
	findOptions := mongoDriverOptions.FindOptions{}
	findOptions.SetSort()

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cur, err := c.driverCollection.Find(ctx, filter, &mongoDriverOptions.FindOptions{
		AllowPartialResults: nil,
		BatchSize:           nil,
		Collation:           nil,
		Comment:             nil,
		CursorType:          nil,
		Hint:                nil,
		Limit:               nil,
		Max:                 nil,
		MaxAwaitTime:        nil,
		MaxTime:             nil,
		Min:                 nil,
		NoCursorTimeout:     nil,
		OplogReplay:         nil,
		Projection:          nil,
		ReturnKey:           nil,
		ShowRecordID:        nil,
		Skip:                nil,
		Snapshot:            nil,
		Sort:                nil,
	})
}

func (c *Collection) UpdateOne(document interface{}, identifier identifier.Identifier) error {
	if identifier == nil {
		return ErrInvalidIdentifier{Reasons: []string{"nil identifier"}}
	} else if err := identifier.IsValid(); err != nil {
		return ErrInvalidIdentifier{Reasons: []string{err.Error()}}
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if _, err := c.driverCollection.ReplaceOne(ctx, identifier.ToFilter(), document); err != nil {
		log.Error().Err(err).Msg("update one")
		return err
	}

	return nil
}
