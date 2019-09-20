package mongo

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

type Database struct {
	mongoClient *mongoDriver.Client
	database    *mongoDriver.Database
}

func New(
	hosts []string,
	connectionString string,
	databaseName string,
) (*Database, error) {

	var db *Database
	var err error

	// try connect with a connection string if one is provided
	if connectionString != "" {
		db, err = NewFromConnectionString(connectionString)
		if err != nil {
			log.Error().Err(err).Msg("connecting to mongo")
			return nil, ErrUnexpected{}
		}
	} else if len(hosts) != 0 {
		db, err = NewFromHosts(hosts)
		if err != nil {
			log.Error().Err(err).Msg("connecting to mongo")
			return nil, ErrUnexpected{}
		}
	} else {
		err = ErrInvalidConfig{Reasons: []string{"no hosts or connection string"}}
		log.Error().Err(err).Msg("connecting to mongo")
		return nil, err
	}

	// connection successful populate and return database
	db.database = db.mongoClient.Database(databaseName)
	return db, nil
}

func NewFromHosts(hosts []string) (*Database, error) {
	log.Info().Msg(fmt.Sprintf(
		"Connecting to mongo cluster on nodes: [%s]",
		strings.Join(hosts, ","),
	))

	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()
	mongoClient, err := mongoDriver.Connect(
		ctx,
		&options.ClientOptions{
			Hosts: hosts,
		})
	if err != nil {
		log.Error().Err(err).Msg("error connecting to mongo")
		return nil, err
	}

	// confirm that the client is connected
	ctx, cancelFn = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()
	if err := mongoClient.Ping(ctx, nil); err != nil {
		log.Error().Err(err).Msg("could not ping mongo")
		return nil, err
	} else {
		log.Info().Msg("connected to mongo")
	}

	return &Database{
		mongoClient: mongoClient,
	}, nil
}

func NewFromConnectionString(connectionString string) (*Database, error) {
	log.Info().Msg("Connecting to mongo with connection string")

	// create a new mongo client
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()
	mongoClient, err := mongoDriver.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Error().Err(err).Msg("connecting to mongo")
		return nil, err
	}

	// confirm that the client is connected
	ctx, cancelFn = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()
	if err := mongoClient.Ping(ctx, nil); err != nil {
		log.Error().Err(err).Msg("could not ping mongo")
		return nil, err
	} else {
		log.Info().Msg("connected to mongo")
	}

	return &Database{
		mongoClient: mongoClient,
	}, nil
}

func (d *Database) CloseConnection() error {
	if err := d.mongoClient.Disconnect(context.Background()); err != nil {
		log.Error().Err(err).Msg("disconnecting from mongo Database")
		return err
	}
	return nil
}

func (d *Database) Collection(collectionName string) *Collection {
	return &Collection{
		driverCollection: d.database.Collection(collectionName),
	}
}
