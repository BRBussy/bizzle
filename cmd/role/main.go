package main

import (
	"flag"
	roleConfig "github.com/BRBussy/bizzle/configs/role"
	"github.com/BRBussy/bizzle/internal/app/role"
	jsonRpcHttpServer "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/server/http"
	jsonRpcServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	basicRoleAdmin "github.com/BRBussy/bizzle/internal/pkg/security/role/admin/basic"
	roleStoreJsonRpcAdaptor "github.com/BRBussy/bizzle/internal/pkg/security/role/store/adaptor/jsonRpc"
	mongoRoleStore "github.com/BRBussy/bizzle/internal/pkg/security/role/store/mongo"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
)

var configFileName = flag.String("config-file-name", "config", "specify config file")

func main() {
	flag.Parse()

	// get exercise/store config
	config, err := roleConfig.GetConfig(configFileName)
	if err != nil {
		log.Fatal().Err(err).Msg("getting config from file")
	}

	// create new mongo db connection
	mongoDb, err := mongo.New(config.MongoDbHosts, config.MongoDBConnectionString, "bizzle")
	if err != nil {
		log.Fatal().Err(err).Msg("creating new mongo db client")
	}
	defer func() {
		if err := mongoDb.CloseConnection(); err != nil {
			log.Error().Err(err).Msg("closing mongo db client connection")
		}
	}()

	// create service providers
	MongoRoleStore, err := mongoRoleStore.New(mongoDb)
	if err != nil {
		log.Fatal().Err(err).Msg("creating mongo role store")
	}
	BasicRoleAdmin := basicRoleAdmin.New(MongoRoleStore)

	// run setup
	if err := role.Setup(
		BasicRoleAdmin,
		MongoRoleStore,
	); err != nil {
		log.Fatal().Err(err).Msg("role setup")
	}

	// create rpc http server
	server := jsonRpcHttpServer.New(
		"/",
		"0.0.0.0",
		config.ServerPort,
	)

	// register service providers
	if err := server.RegisterBatchServiceProviders([]jsonRpcServiceProvider.Provider{
		roleStoreJsonRpcAdaptor.New(MongoRoleStore),
	}); err != nil {
		log.Fatal().Err(err).Msg("registering batch service providers")
	}

	// start server
	log.Info().Msgf("starting exercise json rpc http api server started on port %s", config.ServerPort)
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
