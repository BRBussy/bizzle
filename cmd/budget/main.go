package main

import (
	"flag"
	budgetConfig "github.com/BRBussy/bizzle/configs/budget"
	jsonRpcHttpServer "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/server/http"
	jsonRPCServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	bizzleJSONRPCAuthenticator "github.com/BRBussy/bizzle/internal/pkg/authenticator/jsonRPC"
	budgetAdminJSONRPCAdaptor "github.com/BRBussy/bizzle/internal/pkg/budget/admin/adaptor/jsonRPC"
	basicBudgetAdmin "github.com/BRBussy/bizzle/internal/pkg/budget/admin/basic"
	budgetEntryAdminJSONRPCAdaptor "github.com/BRBussy/bizzle/internal/pkg/budget/entry/admin/adaptor/jsonRPC"
	basicBudgetEntryAdmin "github.com/BRBussy/bizzle/internal/pkg/budget/entry/admin/basic"
	mongoBudgetEntryStore "github.com/BRBussy/bizzle/internal/pkg/budget/entry/store/mongo"
	xlsxStandardBankStatementParser "github.com/BRBussy/bizzle/internal/pkg/budget/statement/parser/XLSXStandardBank"
	"github.com/BRBussy/bizzle/internal/pkg/logs"
	"github.com/BRBussy/bizzle/internal/pkg/middleware"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	jsonRPCTokenValidator "github.com/BRBussy/bizzle/internal/pkg/security/token/validator/jsonRPC"
	requestValidator "github.com/BRBussy/bizzle/pkg/validate/validator/request"
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
	config, err := budgetConfig.GetConfig(configFileName)
	if err != nil {
		log.Fatal().Err(err).Msg("getting config from file")
	}

	// create validator
	RequestValidator := requestValidator.New()

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

	//
	// Budget
	//
	XLSXStandardBankStatementParser := xlsxStandardBankStatementParser.New(
		RequestValidator,
	)
	MongoBudgetEntryStore, err := mongoBudgetEntryStore.New(
		RequestValidator,
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
	BasicBudgetAdmin := basicBudgetAdmin.New(
		RequestValidator,
		MongoBudgetEntryStore,
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
