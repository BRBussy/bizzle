package main

import (
	"cloud.google.com/go/compute/metadata"
	"flag"
	"fmt"
	gatewayConfig "github.com/BRBussy/bizzle/configs/gateway"
	jsonRpcHttpServer "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/server/http"
	jsonRpcServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	authenticatorJsonRpcAdaptor "github.com/BRBussy/bizzle/internal/pkg/authenticator/adaptor/jsonRpc"
	basicAuthenticator "github.com/BRBussy/bizzle/internal/pkg/authenticator/basic"
	"github.com/rs/zerolog/log"
	"net/http"
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

func handler(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Gateway received a request.")

	response, err := makeGetRequest("/")
	if err != nil {
		fmt.Fprintf(w, "could not get to authenticator :( %s!\n", err.Error())
		log.Error().Err(err).Msg("could not get to authenticator!! ")
		return
	}

	var bodyBytes []byte
	if _, err := response.Body.Read(bodyBytes); err != nil {
		log.Error().Err(err).Msg("error reading response body! ")
		fmt.Fprintf(w, "could not read authenticators body bytes :( %s!\n", err.Error())
		return
	}

	fmt.Fprintf(w, "Response from authenticator!!!!! %s!\n", string(bodyBytes))
}

// makeGetRequest makes a GET request to the specified Cloud Run endpoint in
// serviceURL (must be a complete URL) by authenticating with the ID token
// obtained from the Metadata API.
func makeGetRequest(serviceURL string) (*http.Response, error) {
	// create a request
	req, err := http.NewRequest("GET", serviceURL, nil)
	if err != nil {
		return nil, err
	}

	if "prod" == "prod" {
		// query the id_token with ?audience as the serviceURL
		tokenURL := fmt.Sprintf("/instance/service-accounts/default/identity?audience=%s", serviceURL)
		idToken, err := metadata.Get(tokenURL)
		if err != nil {
			return nil, fmt.Errorf("metadata.Get: failed to query id_token: %+v", err)
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", idToken))
	}

	return http.DefaultClient.Do(req)
}
