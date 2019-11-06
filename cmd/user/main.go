package main

import (
	"flag"
	userConfig "github.com/BRBussy/bizzle/configs/user"
	"github.com/BRBussy/bizzle/internal/app/user"
	jsonRpcHttpServer "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/server/http"
	jsonRpcServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	bizzleJSONRPCAuthenticator "github.com/BRBussy/bizzle/internal/pkg/authenticator/jsonRPC"
	"github.com/BRBussy/bizzle/internal/pkg/logs"
	"github.com/BRBussy/bizzle/internal/pkg/middleware"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	jsonRpcRoleStore "github.com/BRBussy/bizzle/internal/pkg/security/role/store/jsonRpc"
	jsonRPCTokenValidator "github.com/BRBussy/bizzle/internal/pkg/security/token/validator/jsonRPC"
	basicUserAdmin "github.com/BRBussy/bizzle/internal/pkg/user/admin/basic"
	userStoreJsonRpcAdaptor "github.com/BRBussy/bizzle/internal/pkg/user/store/adaptor/jsonRpc"
	mongoUserStore "github.com/BRBussy/bizzle/internal/pkg/user/store/mongo"
	basicUserValidator "github.com/BRBussy/bizzle/internal/pkg/user/validator/basic"
	basicValidator "github.com/BRBussy/bizzle/pkg/validate/validator/basic"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
)

var configFileName = flag.String("config-file-name", "config", "specify config file")

func main() {
	flag.Parse()
	logs.Setup()

	// get config
	config, err := userConfig.GetConfig(configFileName)
	if err != nil {
		log.Fatal().Err(err).Msg("getting config from file")
	}

	// create validator
	BasicValidator := basicValidator.New()

	// create new mongo db connection
	mongoDb, err := mongo.New(config.MongoDbHosts, config.MongoDBConnectionString, config.MongoDbName)
	if err != nil {
		log.Fatal().Err(err).Msg("creating new mongo db client")
	}
	defer func() {
		if err := mongoDb.CloseConnection(); err != nil {
			log.Error().Err(err).Msg("closing mongo db client connection")
		}
	}()

	// create service providers
	MongoUserStore, err := mongoUserStore.New(mongoDb)
	if err != nil {
		log.Fatal().Err(err).Msg("creating mongo user role store")
	}
	JSONRPCRoleStore := jsonRpcRoleStore.New(
		BasicValidator,
		config.RoleURL,
		config.PreSharedSecret,
	)
	BasicUserValidator := basicUserValidator.New(JSONRPCRoleStore)
	BasicUserAdmin := basicUserAdmin.New(
		BasicUserValidator,
		MongoUserStore,
		JSONRPCRoleStore,
	)
	JSONRPCTokenValidator := jsonRPCTokenValidator.New(
		config.AuthURL,
		config.PreSharedSecret,
	)
	JSONRPCBizzleAuthenticator := bizzleJSONRPCAuthenticator.New(
		BasicValidator,
		config.AuthURL,
		config.PreSharedSecret,
	)

	// perform setup
	if err := user.Setup(
		BasicUserAdmin,
		MongoUserStore,
		JSONRPCRoleStore,
		config.RootPassword,
	); err != nil {
		log.Fatal().Err(err).Msg("user setup")
	}

	authenticationMiddleware := middleware.NewAuthentication(
		config.PreSharedSecret,
		JSONRPCTokenValidator,
		JSONRPCBizzleAuthenticator,
	)

	// create rpc http server
	server := jsonRpcHttpServer.New(
		"/",
		"0.0.0.0",
		config.ServerPort,
		[]func(http.Handler) http.Handler{
			authenticationMiddleware.Apply,
		},
	)

	// register service providers
	if err := server.RegisterBatchServiceProviders([]jsonRpcServiceProvider.Provider{
		userStoreJsonRpcAdaptor.New(MongoUserStore),
	}); err != nil {
		log.Fatal().Err(err).Msg("registering batch service providers")
	}

	// start server
	log.Info().Msgf("starting user json rpc http api server started on port %s", config.ServerPort)
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
