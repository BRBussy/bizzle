package main

import (
	"flag"
	jsonRPCRoleStore "github.com/BRBussy/bizzle/internal/pkg/security/role/store/jsonRpc"
	basicScopeAdmin "github.com/BRBussy/bizzle/internal/pkg/security/scope/basic"
	"net/http"
	"os"
	"os/signal"

	userConfig "github.com/BRBussy/bizzle/configs/user"
	jsonRpcHttpServer "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/server/http"
	jsonRPCServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	bizzleJSONRPCAuthenticator "github.com/BRBussy/bizzle/internal/pkg/authenticator/jsonRPC"
	"github.com/BRBussy/bizzle/internal/pkg/logs"
	"github.com/BRBussy/bizzle/internal/pkg/middleware"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	jsonRPCTokenValidator "github.com/BRBussy/bizzle/internal/pkg/security/token/validator/jsonRPC"
	userAdminJSONRPCAdaptor "github.com/BRBussy/bizzle/internal/pkg/user/admin/adaptor/jsonRpc"
	basicUserAdmin "github.com/BRBussy/bizzle/internal/pkg/user/admin/basic"
	userStoreJsonRpcAdaptor "github.com/BRBussy/bizzle/internal/pkg/user/store/adaptor/jsonRpc"
	mongoUserStore "github.com/BRBussy/bizzle/internal/pkg/user/store/mongo"
	basicUserValidator "github.com/BRBussy/bizzle/internal/pkg/user/validator/basic"
	requestValidator "github.com/BRBussy/bizzle/pkg/validate/validator/request"
	"github.com/rs/zerolog/log"
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

	//
	// Scope Admin
	//
	BasicScopeAdmin := basicScopeAdmin.New(
		RequestValidator,
	)

	// create service providers
	JSONRPCRoleStore := jsonRPCRoleStore.New(
		RequestValidator,
		config.RoleURL,
		config.PreSharedSecret,
	)

	MongoUserStore, err := mongoUserStore.New(
		RequestValidator,
		BasicScopeAdmin,
		mongoDb,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("creating mongo user role store")
	}
	BasicUserValidator := basicUserValidator.New(
		RequestValidator,
		MongoUserStore,
		JSONRPCRoleStore,
	)
	BasicUserAdmin := basicUserAdmin.New(
		RequestValidator,
		BasicUserValidator,
		MongoUserStore,
		JSONRPCRoleStore,
	)

	JSONRPCTokenValidator := jsonRPCTokenValidator.New(
		config.AuthURL,
		config.PreSharedSecret,
	)
	JSONRPCBizzleAuthenticator := bizzleJSONRPCAuthenticator.New(
		RequestValidator,
		config.AuthURL,
		config.PreSharedSecret,
	)
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
		[]jsonRPCServiceProvider.Provider{
			userStoreJsonRpcAdaptor.New(MongoUserStore),
			userAdminJSONRPCAdaptor.New(BasicUserAdmin),
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
