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

type database struct {
	mongoClient *mongoDriver.Client
	database    *mongoDriver.Database
}

func New(
	hosts []string,
	connectionString string,
	databaseName string,
) (*database, error) {

	var db *database
	var err error

	if connectionString != "" {
		db, err = NewFromConnectionString(connectionString)
	}
	if len(hosts) != 0 {
		db, err = NewFromHosts(hosts)
	}
	if err != nil {
		return nil, err
	}

	db.database = db.mongoClient.Database(databaseName)
	return db, nil
}

func NewFromHosts(hosts []string) (*database, error) {
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

	return &database{
		mongoClient: mongoClient,
	}, nil
}

func NewFromConnectionString(connectionString string) (*database, error) {
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

	return &database{
		mongoClient: mongoClient,
	}, nil
}

func (d *database) CloseConnection() error {
	if err := d.mongoClient.Disconnect(context.Background()); err != nil {
		log.Error().Err(err).Msg("disconnecting from mongo database")
		return err
	}
	return nil
}

func (d *database) CreateOne(document interface{}, collection string) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := d.database.Collection(collection).InsertOne(ctx, document)
	return err
}
