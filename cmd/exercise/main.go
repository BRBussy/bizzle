package main

import (
	"flag"
	exerciseConfig "github.com/BRBussy/bizzle/configs/exercise"
	"github.com/BRBussy/bizzle/internal/app/exercise"
	jsonRpcHttpServer "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/server/http"
	jsonRpcServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	bizzleJSONRPCAuthenticator "github.com/BRBussy/bizzle/internal/pkg/authenticator/jsonRPC"
	exerciseAdminJsonRPCAdaptor "github.com/BRBussy/bizzle/internal/pkg/exercise/admin/adaptor/jsonRPC"
	basicExerciseAdmin "github.com/BRBussy/bizzle/internal/pkg/exercise/admin/basic"
	exerciseStoreJsonRPCAdaptor "github.com/BRBussy/bizzle/internal/pkg/exercise/store/adaptor/jsonRpc"
	mongoExerciseStore "github.com/BRBussy/bizzle/internal/pkg/exercise/store/mongo"
	"github.com/BRBussy/bizzle/internal/pkg/logs"
	"github.com/BRBussy/bizzle/internal/pkg/middleware"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	jsonRPCTokenValidator "github.com/BRBussy/bizzle/internal/pkg/security/token/validator/jsonRPC"
	basicValidator "github.com/BRBussy/bizzle/pkg/validate/validator/basic"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
)

var configFileName = flag.String("config-file-name", "config", "specify config file")

func main() {
	flag.Parse()
	logs.Setup()

	// get config
	config, err := exerciseConfig.GetConfig(configFileName)
	if err != nil {
		log.Fatal().Err(err).Msg("getting config from file")
	}

	// create validator
	BasicValidator := basicValidator.New()

	// create new mongo db connection
	mongoDb, err := mongo.New(config.MongoDbHosts, config.MongoDBConnectionString, config.MongoDbName)
	if err != nil {
		log.Fatal().Err(err).Msg("creating new mongo db client")
	}
	defer func() {
		if err := mongoDb.CloseConnection(); err != nil {
			log.Error().Err(err).Msg("closing mongo db client connection")
		}
	}()

	// create service providers
	MongoExerciseStore, err := mongoExerciseStore.New(
		BasicValidator,
		mongoDb,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("creating mongo exercise store")
	}
	BasicExerciseAdmin := basicExerciseAdmin.New(
		BasicValidator,
		MongoExerciseStore,
	)

	JSONRPCTokenValidator := jsonRPCTokenValidator.New(
		config.AuthURL,
		config.PreSharedSecret,
	)

	JSONRPCBizzleAuthenticator := bizzleJSONRPCAuthenticator.New(
		BasicValidator,
		config.AuthURL,
		config.PreSharedSecret,
	)

	authenticationMiddleware := middleware.NewAuthentication(
		config.PreSharedSecret,
		JSONRPCTokenValidator,
		JSONRPCBizzleAuthenticator,
	)

	// perform setup
	if err := exercise.Setup(BasicExerciseAdmin); err != nil {
		log.Fatal().Err(err).Msg("performing exercise creation")
	}

	// create rpc http server
	server := jsonRpcHttpServer.New(
		"/",
		"0.0.0.0",
		config.ServerPort,
		[]func(http.Handler) http.Handler{
			authenticationMiddleware.Apply,
		},
	)

	// register service providers
	if err := server.RegisterBatchServiceProviders([]jsonRpcServiceProvider.Provider{
		exerciseStoreJsonRPCAdaptor.New(MongoExerciseStore),
		exerciseAdminJsonRPCAdaptor.New(BasicExerciseAdmin),
	}); err != nil {
		log.Fatal().Err(err).Msg("registering batch service providers")
	}

	// start server
	go func() {
		if err := server.Start(); err != nil {
			log.Error().Err(err).Msg("json rpc http api server has stopped")
		}
	}()

	// wait for interrupt signal to stop
	systemSignalsChannel := make(chan os.Signal, 1)
	signal.Notify(systemSignalsChannel, os.Interrupt)
	for s := range systemSignalsChannel {
		log.Info().Msgf("Application is shutting down.. ( %s )", s)
		return
	}
}
