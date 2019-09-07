package main

import (
	"cloud.google.com/go/compute/metadata"
	"flag"
	"fmt"
	gatewayConfig "github.com/BRBussy/bizzle/configs/gateway"
	"github.com/rs/zerolog/log"
	"net/http"
)

var configFileName = flag.String("config-file-name", "config", "specify config file")

func main() {
	flag.Parse()
	var err error

	config, err := gatewayConfig.GetConfig(configFileName)
	if err != nil {
		log.Fatal().Err(err).Msg("getting config from file")
	}

	fmt.Println("conf", *config)

	http.HandleFunc("/", handler)

	log.Fatal().Err(http.ListenAndServe(fmt.Sprintf(":%s", config.ServerPort), nil))
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
