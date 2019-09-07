package main

import (
	"cloud.google.com/go/compute/metadata"
	"flag"
	"fmt"
	gatewayConfig "github.com/BRBussy/bizzle/configs/gateway"
	basicJsonRpcClient "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/client/basic"
	jsonRpcHttpServer "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/server/http"
	jsonRpcServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	authenticatorJsonRpcAdaptor "github.com/BRBussy/bizzle/internal/pkg/authenticator/adaptor/jsonRpc"
	jsonRpcAuthenticator "github.com/BRBussy/bizzle/internal/pkg/authenticator/jsonRpc"
	"github.com/BRBussy/bizzle/internal/pkg/environment"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
)

var configFileName = flag.String("config-file-name", "config", "specify config file")

func main() {
	flag.Parse()

	// get gateway config
	config, err := gatewayConfig.GetConfig(configFileName)
	if err != nil {
		log.Fatal().Err(err).Msg("getting config from file")
	}

	// create json rpc client to access the authenticator
	BasicJsonRpcClient := basicJsonRpcClient.New(config.AuthenticatorURL)
	if config.Environment == environment.Production {
		// query the id_token with ?audience as the authentication serviceURL
		authenticationServiceToken, err := metadata.Get(fmt.Sprintf("/instance/service-accounts/default/identity?audience=%s", config.AuthenticatorURL))
		if err != nil {
			log.Fatal().Err(err).Msg("failed to get authentication service access token")
		}
		log.Info().Msg("obtained authentication service access token")

		// set this as an additional header
		BasicJsonRpcClient.AddAdditionalHeaderEntry("Authorization", fmt.Sprintf("Bearer %s", authenticationServiceToken))
	}

	// create service providers
	JsonRpcAuthenticator := jsonRpcAuthenticator.New(BasicJsonRpcClient)

	// create rpc http server
	server := jsonRpcHttpServer.New(
		"/",
		"0.0.0.0",
		config.ServerPort,
	)

	// register service providers
	if err := server.RegisterBatchServiceProviders([]jsonRpcServiceProvider.Provider{
		authenticatorJsonRpcAdaptor.New(JsonRpcAuthenticator),
	}); err != nil {
		log.Fatal().Err(err).Msg("registering batch service providers")
	}

	// start server
	log.Info().Msgf("starting gateway json rpc http api server started on port %s", config.ServerPort)
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
