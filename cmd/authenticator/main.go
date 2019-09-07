package main

import (
	"flag"
	authenticatorConfig "github.com/BRBussy/bizzle/configs/authenticator"
	jsonRpcHttpServer "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/server/http"
	jsonRpcServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	authenticatorJsonRpcAdaptor "github.com/BRBussy/bizzle/internal/pkg/authenticator/adaptor/jsonRpc"
	basicAuthenticator "github.com/BRBussy/bizzle/internal/pkg/authenticator/basic"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
)

var configFileName = flag.String("config-file-name", "config", "specify config file")

func main() {
	flag.Parse()

	// get authenticator config
	config, err := authenticatorConfig.GetConfig(configFileName)
	if err != nil {
		log.Fatal().Err(err).Msg("getting config from file")
	}

	// create service providers
	BasicAuthenticator := basicAuthenticator.New()

	// create rpc http server
	server := jsonRpcHttpServer.New(
		"/",
		"0.0.0.0",
		config.ServerPort,
	)

	// register service providers
	if err := server.RegisterBatchServiceProviders([]jsonRpcServiceProvider.Provider{
		authenticatorJsonRpcAdaptor.New(BasicAuthenticator),
	}); err != nil {
		log.Fatal().Err(err).Msg("registering batch service providers")
	}

	// start server
	log.Info().Msgf("starting authenticator json rpc http api server started on port %s", config.ServerPort)
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
