package mongoDb

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

type MongoDB struct {
	mongoClient *mongo.Client
}

func New(
	hosts []string,
	connectionString string,
) (*MongoDB, error) {

	if connectionString != "" {
		return NewFromConnectionString(connectionString)
	}
	if len(hosts) != 0 {
		return NewFromHosts(hosts)
	}
	return nil, errors.New("invalid configuration")
}

func NewFromHosts(hosts []string) (*MongoDB, error) {
	log.Info().Msg(fmt.Sprintf(
		"Connecting to mongo cluster on nodes: [%s]",
		strings.Join(hosts, ","),
	))

	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()
	mongoClient, err := mongo.Connect(
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

	return &MongoDB{
		mongoClient: mongoClient,
	}, nil
}

func NewFromConnectionString(connectionString string) (*MongoDB, error) {
	log.Info().Msg("Connecting to mongo with connection string")

	// create a new mongo client
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
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

	return &MongoDB{
		mongoClient: mongoClient,
	}, nil
}

func (m *MongoDB) CloseConnection() error {
	if err := m.mongoClient.Disconnect(context.Background()); err != nil {
		log.Error().Err(err).Msg("disconnecting from mongo database")
		return err
	}
	return nil
}
