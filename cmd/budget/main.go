package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"

	budgetConfig "github.com/BRBussy/bizzle/configs/budget"
	jsonRpcHttpServer "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/server/http"
	jsonRPCServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	bizzleJSONRPCAuthenticator "github.com/BRBussy/bizzle/internal/pkg/authenticator/jsonRPC"
	budgetAdminJSONRPCAdaptor "github.com/BRBussy/bizzle/internal/pkg/budget/admin/adaptor/jsonRPC"
	basicBudgetAdmin "github.com/BRBussy/bizzle/internal/pkg/budget/admin/basic"
	budgetEntryAdminJSONRPCAdaptor "github.com/BRBussy/bizzle/internal/pkg/budget/entry/admin/adaptor/jsonRPC"
	basicBudgetEntryAdmin "github.com/BRBussy/bizzle/internal/pkg/budget/entry/admin/basic"
	basicBudgetCategoryRuleAdmin "github.com/BRBussy/bizzle/internal/pkg/budget/entry/categoryRule/admin/basic"
	budgetCategoryRuleStoreJSONRPCAdaptor "github.com/BRBussy/bizzle/internal/pkg/budget/entry/categoryRule/store/adaptor/jsonRPC"
	mongoBudgetCategoryRuleStore "github.com/BRBussy/bizzle/internal/pkg/budget/entry/categoryRule/store/mongo"
	mongoBudgetEntryStore "github.com/BRBussy/bizzle/internal/pkg/budget/entry/store/mongo"
	basicBudgetEntryValidator "github.com/BRBussy/bizzle/internal/pkg/budget/entry/validator/basic"
	xlsxStandardBankStatementParser "github.com/BRBussy/bizzle/internal/pkg/budget/statement/parser/XLSXStandardBank"
	"github.com/BRBussy/bizzle/internal/pkg/logs"
	"github.com/BRBussy/bizzle/internal/pkg/middleware"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	basicScopeAdmin "github.com/BRBussy/bizzle/internal/pkg/security/scope/basic"
	jsonRPCTokenValidator "github.com/BRBussy/bizzle/internal/pkg/security/token/validator/jsonRPC"
	requestValidator "github.com/BRBussy/bizzle/pkg/validate/validator/request"
	"github.com/rs/zerolog/log"
)

var configFileName = flag.String("config-file-name", "config", "specify config file")

func main() {
	flag.Parse()
	logs.Setup()

	// get config
	config, err := budgetConfig.GetConfig(configFileName)
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

	//
	// Budget
	//
	MongoBudgetCategoryRuleStore, err := mongoBudgetCategoryRuleStore.New(
		RequestValidator,
		BasicScopeAdmin,
		mongoDb,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("creating mongo budget category rule store")
	}
	BasicBudgetCategoryRuleAdmin := basicBudgetCategoryRuleAdmin.New(
		RequestValidator,
		MongoBudgetCategoryRuleStore,
	)
	XLSXStandardBankStatementParser := xlsxStandardBankStatementParser.New(
		RequestValidator,
		BasicBudgetCategoryRuleAdmin,
	)
	MongoBudgetEntryStore, err := mongoBudgetEntryStore.New(
		RequestValidator,
		BasicScopeAdmin,
		mongoDb,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("creating mongo budget entry store")
	}
	BasicBudgetEntryAdmin := basicBudgetEntryAdmin.New(
		RequestValidator,
		MongoBudgetEntryStore,
		XLSXStandardBankStatementParser,
	)
	BasicBudgetEntryValidator := basicBudgetEntryValidator.New(
		RequestValidator,
		MongoBudgetCategoryRuleStore,
	)
	BasicBudgetAdmin := basicBudgetAdmin.New(
		RequestValidator,
		MongoBudgetEntryStore,
		BasicBudgetEntryValidator,
	)

	//
	// Authentication
	//
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
			budgetEntryAdminJSONRPCAdaptor.New(BasicBudgetEntryAdmin),
			budgetAdminJSONRPCAdaptor.New(BasicBudgetAdmin),
			budgetCategoryRuleStoreJSONRPCAdaptor.New(MongoBudgetCategoryRuleStore),
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
