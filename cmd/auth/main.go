package main

import (
	"flag"
	authConfig "github.com/BRBussy/bizzle/configs/auth"
	jsonRpcHttpServer "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/server/http"
	jsonRPCServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	authenticatorJSONRPCAdaptor "github.com/BRBussy/bizzle/internal/pkg/authenticator/adaptor/jsonRpc"
	basicAuthenticator "github.com/BRBussy/bizzle/internal/pkg/authenticator/basic"
	"github.com/BRBussy/bizzle/internal/pkg/logs"
	"github.com/BRBussy/bizzle/internal/pkg/middleware"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	jsonRPCRoleStore "github.com/BRBussy/bizzle/internal/pkg/security/role/store/jsonRpc"
	basicTokenGenerator "github.com/BRBussy/bizzle/internal/pkg/security/token/generator/basic"
	tokenValidatorJSONRPCAdaptor "github.com/BRBussy/bizzle/internal/pkg/security/token/validator/adaptor/jsonRPC"
	basicTokenValidator "github.com/BRBussy/bizzle/internal/pkg/security/token/validator/basic"
	jsonRPCUserStore "github.com/BRBussy/bizzle/internal/pkg/user/store/jsonRPC"
	"github.com/BRBussy/bizzle/pkg/key"
	requestValidator "github.com/BRBussy/bizzle/pkg/validate/validator/request"
	"github.com/rs/zerolog/log"
	"gopkg.in/square/go-jose.v2"
	"net/http"
	"os"
	"os/signal"
)

var configFileName = flag.String("config-file-name", "config", "specify config file")

func main() {
	flag.Parse()
	logs.Setup()

	// get config
	config, err := authConfig.GetConfig(configFileName)
	if err != nil {
		log.Fatal().Err(err).Msg("getting config from file")
	}

	// create validator
	RequestValidator := requestValidator.New()

	// create new mongo db connection
	mongoDb, err := mongo.New(
		config.MongoDBHosts,
		config.MongoDBUsername,
		config.MongoDBPassword,
		config.MongoDBConnectionString,
		config.MongoDBName,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("creating new mongo db client")
	}
	defer func() {
		if err := mongoDb.CloseConnection(); err != nil {
			log.Error().Err(err).Msg("closing mongo db client connection")
		}
	}()

	// fetch or generate RSA key pair
	rsaKeyPair, err := key.ParseRSAPrivateKeyFromString(config.PrivateKeyString)
	if err != nil {
		log.Fatal().Err(err).Msg("parsing private key")
	}

	// create a new signer using RSASSA-PSS (SHA512) with the given private key.
	joseSigner, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.PS512, Key: rsaKeyPair}, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("generating new jose signer")
	}

	JSORPCUserStore := jsonRPCUserStore.New(
		RequestValidator,
		config.UserURL,
		config.PreSharedSecret,
	)
	JSONRPCRoleStore := jsonRPCRoleStore.New(
		RequestValidator,
		config.RoleURL,
		config.PreSharedSecret,
	)

	BasicTokenGenerator := basicTokenGenerator.New(
		joseSigner,
		RequestValidator,
	)
	BasicTokenValidator := basicTokenValidator.New(
		rsaKeyPair,
		RequestValidator,
	)

	// create authenticator
	BasicAuthenticator := basicAuthenticator.New(
		JSORPCUserStore,
		JSONRPCRoleStore,
		BasicTokenGenerator,
		RequestValidator,
		mongoDb,
	)

	authenticationMiddleware := middleware.NewAuthentication(
		config.PreSharedSecret,
		BasicTokenValidator,
		BasicAuthenticator,
	)

	// create rpc http server
	server := jsonRpcHttpServer.New(
		"/",
		"0.0.0.0",
		config.ServerPort,
		[]func(http.Handler) http.Handler{
			authenticationMiddleware.Apply,
		},
		[]jsonRPCServiceProvider.Provider{
			authenticatorJSONRPCAdaptor.New(BasicAuthenticator),
			tokenValidatorJSONRPCAdaptor.New(BasicTokenValidator),
		},
	)

	// start server
	go func() {
		if err := server.Start(); err != nil {
			log.Error().Err(err).Msg("json rpc http api server has stopped")
		}
	}()

	// wait for interrupt signal to stop
	systemSignalsChannel := make(chan os.Signal, 1)
	signal.Notify(systemSignalsChannel, os.Interrupt)
	for s := range systemSignalsChannel {
		log.Info().Msgf("Application is shutting down.. ( %s )", s)
		return
	}
}
