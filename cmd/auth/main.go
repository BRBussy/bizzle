package main

import (
	"flag"
	authConfig "github.com/BRBussy/bizzle/configs/auth"
	jsonRpcHttpServer "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/server/http"
	jsonRpcServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	authenticatorJsonRpcAdaptor "github.com/BRBussy/bizzle/internal/pkg/authenticator/adaptor/jsonRpc"
	basicAuthenticator "github.com/BRBussy/bizzle/internal/pkg/authenticator/basic"
	"github.com/BRBussy/bizzle/internal/pkg/middleware"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
)

var configFileName = flag.String("config-file-name", "config", "specify config file")

func main() {
	flag.Parse()

	// get config
	config, err := authConfig.GetConfig(configFileName)
	if err != nil {
		log.Fatal().Err(err).Msg("getting config from file")
	}

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

	// create authenticator
	BasicAuthenticator := new(basicAuthenticator.Authenticator).Setup()

	authenticationMiddleware := new(middleware.Authentication).Setup(
		config.PreSharedSecret,
	)

	// create rpc http server
	server := jsonRpcHttpServer.New(
		"/",
		"0.0.0.0",
		config.ServerPort,
		[]mux.MiddlewareFunc{
			authenticationMiddleware.Apply,
		},
	)

	// register service providers
	if err := server.RegisterBatchServiceProviders([]jsonRpcServiceProvider.Provider{
		authenticatorJsonRpcAdaptor.New(BasicAuthenticator),
	}); err != nil {
		log.Fatal().Err(err).Msg("registering batch service providers")
	}

	// start server
	log.Info().Msgf("starting auth json rpc http api server started on port %s", config.ServerPort)
	go func() {
		if err := server.Start(); err != nil {
			log.Error().Err(err).Msg("json rpc http api server has stopped")
		}
	}()

	// wait for interrupt signal to stop
	systemSignalsChannel := make(chan os.Signal, 1)
	signal.Notify(systemSignalsChannel, os.Interrupt)
	for {
		select {
		case s := <-systemSignalsChannel:
			log.Info().Msgf("Application is shutting down.. ( %s )", s)
			return
		}
	}
}
