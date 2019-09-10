package mongoDb

import (
	"context"
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

func NewFromNodes(hosts []string) (*MongoDB, error) {
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

	return &MongoDB{
		mongoClient: mongoClient,
	}, nil
}

func NewFromConnectionString(connectionString string) (*MongoDB, error) {
	log.Info().Msg("Connecting to mongo with connection string")

	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Error().Err(err).Msg("connecting to mongo")
		return nil, err
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
