package main

import (
	"cloud.google.com/go/compute/metadata"
	"fmt"
	"github.com/BRBussy/bizzle/package/authenticator"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Print("Gateway received a request.")
	target := os.Getenv("TARGET")
	if target == "" {
		target = "World"
	}

	fmt.Fprintf(w, "Hello from gateway yet again!!! %s!\n", target)
}

// makeGetRequest makes a GET request to the specified Cloud Run endpoint in
// serviceURL (must be a complete URL) by authenticating with the ID token
// obtained from the Metadata API.
func makeGetRequest(serviceURL string) (*http.Response, error) {
	// query the id_token with ?audience as the serviceURL
	tokenURL := fmt.Sprintf("/instance/service-accounts/default/identity?audience=%s", serviceURL)
	idToken, err := metadata.Get(tokenURL)
	if err != nil {
		return nil, fmt.Errorf("metadata.Get: failed to query id_token: %+v", err)
	}
	req, err := http.NewRequest("GET", serviceURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", idToken))
	return http.DefaultClient.Do(req)
}

func main() {
	logrus.Info("The bizzle gateway has started!")
	authenticator.Auth()

	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
