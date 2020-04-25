package main

import (
	"flag"
	setupConfig "github.com/BRBussy/bizzle/configs/setup"
	"github.com/BRBussy/bizzle/internal/app/role"
	"github.com/BRBussy/bizzle/internal/pkg/logs"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	basicRoleAdmin "github.com/BRBussy/bizzle/internal/pkg/security/role/admin/basic"
	mongoRoleStore "github.com/BRBussy/bizzle/internal/pkg/security/role/store/mongo"
	requestValidator "github.com/BRBussy/bizzle/pkg/validate/validator/request"
	"github.com/rs/zerolog/log"
)

var configFileName = flag.String("config-file-name", "config", "specify config file")

func main() {
	flag.Parse()
	logs.Setup()

	// get config
	config, err := setupConfig.GetConfig(configFileName)
	if err != nil {
		log.Fatal().Err(err).Msg("getting config from file")
	}

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
	// Request Validator
	//
	RequestValidator := requestValidator.New()

	//
	// Role
	//
	MongoRoleStore, err := mongoRoleStore.New(
		RequestValidator,
		mongoDb,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("creating mongo role store")
	}
	BasicRoleAdmin := basicRoleAdmin.New(MongoRoleStore)

	log.Info().Msg("Running role setup")
	if err := role.Setup(
		BasicRoleAdmin,
		MongoRoleStore,
	); err != nil {
		log.Fatal().Err(err).Msg("role setup failed")
	}
}
