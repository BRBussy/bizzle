package main

import (
	"cloud.google.com/go/compute/metadata"
	"flag"
	"fmt"
	"github.com/BRBussy/bizzle/internal/pkg/secrets"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var env = flag.String("env", "prod", "specify environment")
var secretsFile = flag.String("secrets-file", "secrets.json", "specify environment")
var bizzleSecrets *secrets.Secrets

func handler(w http.ResponseWriter, r *http.Request) {
	log.Info("Gateway received a request.")

	response, err := makeGetRequest(bizzleSecrets.AuthenticatorURL)
	if err != nil {
		fmt.Fprintf(w, "could not get to authenticator :( %s!\n", err.Error())
		log.Error("could not get to authenticator!! ", err)
		return
	}

	var bodyBytes []byte
	if _, err := response.Body.Read(bodyBytes); err != nil {
		log.Error("error reading response body! ", err)
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

	if *env == "prod" {
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

func main() {
	flag.Parse()
	var err error

	// get the secrets
	bizzleSecrets, err = secrets.GetSecrets(*secretsFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("%v\n", *bizzleSecrets)

	log.Info("The bizzle gateway has started!")

	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
