package main

import (
	"flag"
	setupConfig "github.com/BRBussy/bizzle/configs/setup"
	"github.com/BRBussy/bizzle/internal/app/exercise"
	"github.com/BRBussy/bizzle/internal/app/role"
	"github.com/BRBussy/bizzle/internal/app/user"
	basicExerciseAdmin "github.com/BRBussy/bizzle/internal/pkg/exercise/admin/basic"
	mongoExerciseStore "github.com/BRBussy/bizzle/internal/pkg/exercise/store/mongo"
	"github.com/BRBussy/bizzle/internal/pkg/logs"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	basicRoleAdmin "github.com/BRBussy/bizzle/internal/pkg/security/role/admin/basic"
	mongoRoleStore "github.com/BRBussy/bizzle/internal/pkg/security/role/store/mongo"
	basicUserAdmin "github.com/BRBussy/bizzle/internal/pkg/user/admin/basic"
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

	//
	// User
	//
	MongoUserStore, err := mongoUserStore.New(mongoDb)
	if err != nil {
		log.Fatal().Err(err).Msg("creating mongo user role store")
	}
	BasicUserValidator := basicUserValidator.New(MongoRoleStore)
	BasicUserAdmin := basicUserAdmin.New(
		BasicUserValidator,
		MongoUserStore,
		MongoRoleStore,
	)

	//
	// Exercise
	//
	MongoExerciseStore, err := mongoExerciseStore.New(
		RequestValidator,
		mongoDb,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("creating mongo exercise store")
	}
	BasicExerciseAdmin := basicExerciseAdmin.New(
		RequestValidator,
		MongoExerciseStore,
	)

	log.Info().Msg("Running role setup")
	if err := role.Setup(
		BasicRoleAdmin,
		MongoRoleStore,
	); err != nil {
		log.Fatal().Err(err).Msg("role setup")
	}

	log.Info().Msg("Running user setup")
	if err := user.Setup(
		BasicUserAdmin,
		MongoUserStore,
		MongoRoleStore,
		config.RootPassword,
	); err != nil {
		log.Fatal().Err(err).Msg("user setup")
	}

	log.Info().Msg("Running exercise setup")
	if err := exercise.Setup(
		BasicExerciseAdmin,
		MongoExerciseStore,
	); err != nil {
		log.Fatal().Err(err).Msg("performing exercise setup")
	}
}
