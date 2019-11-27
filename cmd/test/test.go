package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/gorilla/rpc/v2"
	gorillaJson "github.com/gorilla/rpc/v2/json2"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	testAdaptor := &TestAdaptor{}

	// create new gorilla rpc server
	rpcServer := rpc.NewServer()
	rpcServer.RegisterCodec(gorillaJson.NewCodec(), "application/json")

	if err := rpcServer.RegisterService(testAdaptor, "Test"); err != nil {
		log.Error().Err(err).Msg("could not register Test service provider")
	}

	// create chi root router and apply middleware
	router := chi.NewRouter()
	router.Use(
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Println("pre flight!")
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "POST")
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Access-Control-Allow-Origin")
				w.WriteHeader(http.StatusOK)
				next.ServeHTTP(w, r)
			})
		},
	)

	// create chi api router
	apiRouter := chi.NewRouter()

	router.Mount("/api", apiRouter)

	apiRouter.Post("/", func() http.HandlerFunc { return rpcServer.ServeHTTP }())
	apiRouter.Options("/", preFlightHandler)

	go func() {
		if err := http.ListenAndServe("0.0.0.0:8080", router); err != nil {
			log.Error().Err(err).Msg("server stopped")
		}
	}()
	log.Info().Msg("started server")

	// wait for interrupt signal to stop
	systemSignalsChannel := make(chan os.Signal, 1)
	signal.Notify(systemSignalsChannel, os.Interrupt)
	for s := range systemSignalsChannel {
		log.Info().Msgf("Application is shutting down.. ( %s )", s)
		return
	}
}

type TestAdaptor struct {
}

type TestRequest struct {
}

type TestResponse struct {
}

func (t *TestAdaptor) Test(r *http.Request, request *TestRequest, response *TestResponse) error {
	fmt.Println("test running!!!")
	return nil
}

func preFlightHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("pre flight!")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Access-Control-Allow-Origin")
	w.WriteHeader(http.StatusOK)
}
