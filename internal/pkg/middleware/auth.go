package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	bizzleAuthenticator "github.com/BRBussy/bizzle/internal/pkg/authenticator"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
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
		// get json rpc method
		method, err := getMethod(r)
		if err != nil {
			log.Error().Err(err).Msg("cannot get jsonrpc method")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// if method is login no authentication is required
		// allow request to pass to service provider
		if method == bizzleAuthenticator.LoginService {
			next.ServeHTTP(w, r)
			return
		}

		// if a valid pre-shared secret is present then this is an inter-microservice call
		// if there is a claims header the value needs to be placed into the request context
		if pss := r.Header.Get("Pre-Shared-Secret"); pss == a.preSharedSecret {
			// look for claims header
			serializedClaims := r.Header.Get("Claims")
			if serializedClaims == "" {
				// if it is blank - do nothing and carry on
				next.ServeHTTP(w, r)
				return
			}
			// otherwise place claims into request context and carry on
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "Claims", serializedClaims)))
			return
		}

		// otherwise this request originates from outside of the bizzle cluster
		// an authorisation header with a valid jwt token should be present
		// derive auth claims from it and place it into the request context
		jwt := r.Header.Get("Authorization")
		if jwt == "" {
			log.Error().Err(err).Msg("no token in authentication header")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// validate token to get claims
		validateResponse, err := a.tokenValidator.Validate(
			tokenValidator.ValidateRequest{
				Token: jwt,
			},
		)
		if err != nil {
			log.Error().Err(err).Msg("token validation failure")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// validate service access
		if _, err := a.authenticator.AuthenticateService(
			bizzleAuthenticator.AuthenticateServiceRequest{
				Claims:  validateResponse.Claims,
				Service: method,
			},
		); err != nil {
			log.Error().Err(err).Msg("unauthorized to access service")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// marshall claims to put into context
		marshalledClaims, err := json.Marshal(claims.Serialized{Claims: validateResponse.Claims})
		if err != nil {
			log.Error().Err(err).Msg("could not marshall claims")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "Claims", marshalledClaims)))
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
