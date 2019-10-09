package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/BRBussy/bizzle/internal/pkg/authenticator"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
)

type Authentication struct {
	preSharedSecret string
}

func (a *Authentication) Setup(
	preSharedSecret string,
) *Authentication {
	a.preSharedSecret = preSharedSecret
	return a
}

func (a *Authentication) Apply(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// look for pre-shared secret header
		pss := r.Header.Get("Pre-Shared-Secret")
		if pss == a.preSharedSecret {
			// if pre-shared secret present no authentication required
			next.ServeHTTP(w, r)
			return
		}

		// get json rpc method
		method, err := getMethod(r)
		if err != nil {
			log.Error().Err(err).Msg("cannot get jsonrpc method")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// if method is login not authentication is required
		// allow request to pass to service provider
		if method == authenticator.LoginService {
			next.ServeHTTP(w, r)
			return
		}

		// all other requests need to be authenticated

		next.ServeHTTP(w, r)
	})
}

func getMethod(r *http.Request) (string, error) {
	// Confirm that body of request has data
	if r.Body == nil {
		return "", errors.New("body is nil")
	}

	// Extract body of http Request
	var bodyBytes []byte
	bodyBytes, _ = ioutil.ReadAll(r.Body)

	// Reset body of request
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	// Retrieve id and method of json rpc request
	var req struct {
		// To unmarshal the received json
		Id     string `json:"id"`
		Method string `json:"method"`
	}
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		return "", err
	}

	return req.Method, nil
}
