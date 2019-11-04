package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	bizzleAuthenticator "github.com/BRBussy/bizzle/internal/pkg/authenticator"
	tokenValidator "github.com/BRBussy/bizzle/internal/pkg/security/token/validator"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
)

type Authentication struct {
	preSharedSecret string
	tokenValidator  tokenValidator.Validator
	authenticator   bizzleAuthenticator.Authenticator
}

func NewAuthentication(
	preSharedSecret string,
	tokenValidator tokenValidator.Validator,
	authenticator bizzleAuthenticator.Authenticator,
) *Authentication {

	return &Authentication{
		preSharedSecret: preSharedSecret,
		tokenValidator:  tokenValidator,
		authenticator:   authenticator,
	}
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
		if method == bizzleAuthenticator.LoginService {
			next.ServeHTTP(w, r)
			return
		}

		// all other requests need to be authenticated with a token header
		// look for token in header
		jwt := r.Header.Get("Authorization")
		if jwt == "" {
			log.Error().Err(err).Msg("not token in authentication header")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// validate token to get claims
		validateResponse, err := a.tokenValidator.Validate(&tokenValidator.ValidateRequest{
			Token: jwt,
		})
		if err != nil {
			log.Error().Err(err).Msg("token validation failure")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		fmt.Println(validateResponse)

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
		Id     int    `json:"id"`
		Method string `json:"method"`
	}
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		return "", err
	}

	return req.Method, nil
}
