package mongo

import (
	"context"
	"github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/pkg/search/criteria"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	"github.com/rs/zerolog/log"
	mongoBSON "go.mongodb.org/mongo-driver/bson"
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

func (c *Collection) CreateMany(documents []interface{}) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := c.driverCollection.InsertMany(ctx, documents)
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

func (c *Collection) FindMany(documents interface{}, criteria criteria.Criteria, query Query) (int64, error) {
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
			return 0, ErrInvalidCriteria{Reasons: reasonsInvalid}
		}
		filter = criteria.ToFilter()
	}

	// get options
	findOptions, err := query.ToMongoFindOptions()
	if err != nil {
		return 0, err
	}

	// perform find
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cur, err := c.driverCollection.Find(ctx, filter, findOptions)
	if err != nil {
		log.Error().Err(err).Msg("performing find")
		return 0, ErrUnexpected{}
	}

	// decode the results
	if err := cur.All(ctx, documents); err != nil {
		log.Error().Err(err).Msg("decoding documents")
		return 0, ErrUnexpected{}
	}

	// get document count
	count, err := c.driverCollection.CountDocuments(ctx, filter)
	if err != nil {
		log.Error().Err(err).Msg("counting documents")
		return 0, ErrUnexpected{}
	}

	return count, nil
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

type aggregationCountHolder struct {
	Count int64 `bson:"count"`
}

func (c *Collection) Aggregate(pipeline mongoDriver.Pipeline, query Query, entities interface{}) (int64, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	// perform aggregation and output count
	countCursor, err := c.driverCollection.Aggregate(
		ctx,
		append(
			pipeline,
			mongoBSON.D{{Key: "$count", Value: "count"}},
		),
	)
	if err != nil {
		log.Error().Err(err).Msg("could not perform count")
		return -1, err
	}
	var countResults []aggregationCountHolder
	if err := countCursor.All(ctx, &countResults); err != nil {
		log.Error().Err(err).Msg("could not decode count")
		return -1, err
	}
	var count int64
	if len(countResults) == 1 {
		count = countResults[0].Count
	} else if len(countResults) == 0 {
		count = 0
	} else {
		log.Error().Msg("invalid count result")
		return -1, exception.ErrUnexpected{}
	}

	// perform aggregation and output documents with query applied
	cursor, err := c.driverCollection.Aggregate(
		ctx,
		append(pipeline, query.ToPipelineStages()...),
	)
	if err != nil {
		log.Error().Err(err).Msg("Could not create cursor for '" + c.driverCollection.Name() + " collection")
		return -1, err
	}

	// decode the results
	if err := cursor.All(ctx, entities); err != nil {
		log.Error().Err(err).Msg("decoding documents")
		return -1, err
	}

	return count, nil
}
