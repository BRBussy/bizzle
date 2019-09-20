package main

import (
	"flag"
	userConfig "github.com/BRBussy/bizzle/configs/user"
	"github.com/BRBussy/bizzle/internal/app/user"
	jsonRpcHttpServer "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/server/http"
	jsonRpcServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	"github.com/BRBussy/bizzle/internal/pkg/firebase"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	jsonRpcRoleStore "github.com/BRBussy/bizzle/internal/pkg/security/role/store/jsonRpc"
	basicUserAdmin "github.com/BRBussy/bizzle/internal/pkg/user/admin/basic"
	userStoreJsonRpcAdaptor "github.com/BRBussy/bizzle/internal/pkg/user/store/adaptor/jsonRpc"
	mongoUserStore "github.com/BRBussy/bizzle/internal/pkg/user/store/mongo"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
)

var configFileName = flag.String("config-file-name", "config", "specify config file")

func main() {
	flag.Parse()

	// get config
	config, err := userConfig.GetConfig(configFileName)
	if err != nil {
		log.Fatal().Err(err).Msg("getting config from file")
	}

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

	// create firebase
	Firebase, err := firebase.New(config.FirebaseCredentialsPath)
	if err != nil {
		log.Fatal().Err(err).Msg("creating firebase")
	}

	// create service providers
	MongoUserStore, err := mongoUserStore.New(mongoDb)
	if err != nil {
		log.Fatal().Err(err).Msg("creating mongo user role store")
	}
	JSONRPCRoleStore := jsonRpcRoleStore.New(
		config.RoleURL,
	)
	BasicUserAdmin := basicUserAdmin.New(
		JSONRPCRoleStore,
		Firebase,
	)

	// perform setup
	if err := user.Setup(
		BasicUserAdmin,
		MongoUserStore,
		JSONRPCRoleStore,
		Firebase,
		config.RootPassword,
	); err != nil {
		log.Fatal().Err(err).Msg("user setup")
	}

	// create rpc http server
	server := jsonRpcHttpServer.New(
		"/",
		"0.0.0.0",
		config.ServerPort,
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
