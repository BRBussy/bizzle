package main

import (
	"flag"
	syncBudgetCategoryRulesForUserConfig "github.com/BRBussy/bizzle/configs/syncBudgetCategoryRulesForUser"
	"github.com/BRBussy/bizzle/internal/app/budget"
	budgetEntryCategoryRuleBasicAdmin "github.com/BRBussy/bizzle/internal/pkg/budget/entry/categoryRule/admin/basic"
	budgetEntryCategoryRuleMongoStore "github.com/BRBussy/bizzle/internal/pkg/budget/entry/categoryRule/store/mongo"
	"github.com/BRBussy/bizzle/internal/pkg/logs"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	userMongoStore "github.com/BRBussy/bizzle/internal/pkg/user/store/mongo"
	requestValidator "github.com/BRBussy/bizzle/pkg/validate/validator/request"
	"github.com/rs/zerolog/log"
)

var configFileName = flag.String("config-file-name", "config", "specify config file")

func main() {
	flag.Parse()
	logs.Setup()

	// get config
	config, err := syncBudgetCategoryRulesForUserConfig.GetConfig(configFileName)
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

	// create validator
	RequestValidator := requestValidator.New()

	//
	// User
	//
	UserMongoStore, err := userMongoStore.New(
		RequestValidator,
		mongoDb,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("creating user mongo store")
	}

	//
	// Budget
	//
	BudgetEntryCategoryRuleMongoStore, err := budgetEntryCategoryRuleMongoStore.New(
		RequestValidator,
		mongoDb,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("creating budget entry category rule mongo store")
	}

	BudgetEntryCategoryRuleBasicAdmin := budgetEntryCategoryRuleBasicAdmin.New(
		RequestValidator,
		BudgetEntryCategoryRuleMongoStore,
	)

	if err := budget.SyncBudgetCategoryRulesForUser(
		config.UserID,
		BudgetEntryCategoryRuleBasicAdmin,
		BudgetEntryCategoryRuleMongoStore,
		UserMongoStore,
	); err != nil {
		log.Fatal().Err(err).Msg("running SyncBudgetCategoryRulesForUser")
	}
}
