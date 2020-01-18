package main

import (
	"flag"
	exerciseConfig "github.com/BRBussy/bizzle/configs/exercise"
	jsonRpcHttpServer "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/server/http"
	jsonRpcServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	bizzleJSONRPCAuthenticator "github.com/BRBussy/bizzle/internal/pkg/authenticator/jsonRPC"
	exerciseAdminJsonRPCAdaptor "github.com/BRBussy/bizzle/internal/pkg/exercise/admin/adaptor/jsonRPC"
	basicExerciseAdmin "github.com/BRBussy/bizzle/internal/pkg/exercise/admin/basic"
	sessionAdminJSONRPCAdaptor "github.com/BRBussy/bizzle/internal/pkg/exercise/session/admin/adaptor/jsonRpc"
	basicSessionAdmin "github.com/BRBussy/bizzle/internal/pkg/exercise/session/admin/basic"
	sessionStoreJSONRPCAdaptor "github.com/BRBussy/bizzle/internal/pkg/exercise/session/store/adaptor/jsonRpc"
	mongoSessionStore "github.com/BRBussy/bizzle/internal/pkg/exercise/session/store/mongo"
	exerciseStoreJsonRPCAdaptor "github.com/BRBussy/bizzle/internal/pkg/exercise/store/adaptor/jsonRpc"
	mongoExerciseStore "github.com/BRBussy/bizzle/internal/pkg/exercise/store/mongo"
	"github.com/BRBussy/bizzle/internal/pkg/logs"
	"github.com/BRBussy/bizzle/internal/pkg/middleware"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	jsonRPCTokenValidator "github.com/BRBussy/bizzle/internal/pkg/security/token/validator/jsonRPC"
	requestValidator "github.com/BRBussy/bizzle/pkg/validate/validator/request"
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
	RequestValidator := requestValidator.New()

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

	//
	// Exercise
	//
	MongoExerciseStore, err := mongoExerciseStore.New(
		RequestValidator,
		mongoDb,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("creating mongo exercise store")
	}
	BasicExerciseAdmin := basicExerciseAdmin.New(
		RequestValidator,
		MongoExerciseStore,
	)

	//
	// Session
	//
	MongoSessionStore, err := mongoSessionStore.New(
		RequestValidator,
		mongoDb,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("creating session store")
	}
	BasicSessionAdmin := basicSessionAdmin.New(
		RequestValidator,
		MongoSessionStore,
	)

	//
	// Authentication
	//
	JSONRPCTokenValidator := jsonRPCTokenValidator.New(
		config.AuthURL,
		config.PreSharedSecret,
	)
	JSONRPCBizzleAuthenticator := bizzleJSONRPCAuthenticator.New(
		RequestValidator,
		config.AuthURL,
		config.PreSharedSecret,
	)
	authenticationMiddleware := middleware.NewAuthentication(
		config.PreSharedSecret,
		JSONRPCTokenValidator,
		JSONRPCBizzleAuthenticator,
	)

	// create rpc http server
	server := jsonRpcHttpServer.New(
		"/",
		"0.0.0.0",
		config.ServerPort,
		[]func(http.Handler) http.Handler{
			authenticationMiddleware.Apply,
		},
		[]jsonRpcServiceProvider.Provider{
			exerciseStoreJsonRPCAdaptor.New(MongoExerciseStore),
			exerciseAdminJsonRPCAdaptor.New(BasicExerciseAdmin),
			sessionStoreJSONRPCAdaptor.New(MongoSessionStore),
			sessionAdminJSONRPCAdaptor.New(BasicSessionAdmin),
		},
	)

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
