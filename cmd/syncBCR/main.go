package main

import (
	"flag"
	basicBudgetConfigAdmin "github.com/BRBussy/bizzle/internal/pkg/budget/config/admin/basic"
	mongoBudgetConfigStore "github.com/BRBussy/bizzle/internal/pkg/budget/config/store/mongo"
	basicBudgetConfigValidator "github.com/BRBussy/bizzle/internal/pkg/budget/config/validator/basic"
	basicScopeAdmin "github.com/BRBussy/bizzle/internal/pkg/security/scope/basic"

	syncBudgetCategoryRulesForUserConfig "github.com/BRBussy/bizzle/configs/syncBCR"
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

	BasicScopeAdmin := basicScopeAdmin.New(
		RequestValidator,
	)

	//
	// User
	//
	UserMongoStore, err := userMongoStore.New(
		RequestValidator,
		BasicScopeAdmin,
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
		BasicScopeAdmin,
		mongoDb,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("creating budget entry category rule mongo store")
	}

	MongoBudgetConfigStore, err := mongoBudgetConfigStore.New(
		RequestValidator,
		BasicScopeAdmin,
		mongoDb,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("creating mongo budget config store")
	}
	BasicBudgetConfigValidator := basicBudgetConfigValidator.New(
		RequestValidator,
		BudgetEntryCategoryRuleMongoStore,
	)
	BasicBudgetConfigAdmin := basicBudgetConfigAdmin.New(
		RequestValidator,
		MongoBudgetConfigStore,
		BasicBudgetConfigValidator,
	)
	BudgetEntryCategoryRuleBasicAdmin := budgetEntryCategoryRuleBasicAdmin.New(
		RequestValidator,
		BudgetEntryCategoryRuleMongoStore,
		BasicBudgetConfigAdmin,
	)

	if err := budget.SyncBudgetCategoryRulesForUser(
		config.UserID,
		BudgetEntryCategoryRuleBasicAdmin,
		BudgetEntryCategoryRuleMongoStore,
		UserMongoStore,
		config.Rules,
	); err != nil {
		log.Fatal().Err(err).Msg("running SyncBudgetCategoryRulesForUser")
	}
}
