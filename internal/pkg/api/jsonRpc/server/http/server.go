package http

import (
	jsonRPCServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	"github.com/go-chi/chi"
	"github.com/gorilla/rpc/v2"
	gorillaJson "github.com/gorilla/rpc/v2/json2"
	"github.com/rs/zerolog/log"
	netHttp "net/http"
)

type server struct {
	path             string
	host             string
	port             string
	rpcServer        *rpc.Server
	rootRouter       *chi.Mux
	apiRouter        *chi.Mux
	serviceProviders []jsonRPCServiceProvider.Provider
}

func New(
	path string,
	host string,
	port string,
	middleware []func(netHttp.Handler) netHttp.Handler,
	serviceProviders []jsonRPCServiceProvider.Provider,
) *server {
	// create a new server
	newServer := new(server)
	newServer.path = path
	newServer.host = host
	newServer.port = port

	// create new gorilla rpc server
	newServer.rpcServer = rpc.NewServer()
	newServer.rpcServer.RegisterCodec(gorillaJson.NewCodec(), "application/json")

	for _, serviceProvider := range serviceProviders {
		log.Info().Msg("registering: " + serviceProvider.Name().String())
		if err := newServer.rpcServer.RegisterService(serviceProvider, serviceProvider.Name().String()); err != nil {
			log.Fatal().Err(err).Msg("could not register: " + serviceProvider.Name().String())
		}
	}

	// initialise middleware
	if middleware == nil {
		middleware = make([]func(netHttp.Handler) netHttp.Handler, 0)
	}

	// create chi root router and apply middleware
	newServer.rootRouter = chi.NewRouter()
	newServer.rootRouter.Use(preFlightAndCORSHandler)
	for _, m := range middleware {
		newServer.rootRouter.Use(m)
	}

	// create chi api router
	newServer.apiRouter = chi.NewRouter()

	newServer.apiRouter.Post("/", func() netHttp.HandlerFunc { return newServer.rpcServer.ServeHTTP }())
	newServer.rootRouter.Mount("/api", newServer.apiRouter)

	return newServer
}

func (s *server) Start() error {
	log.Info().Msg("starting http json rpc api server on: " + s.host + ":" + s.port + "/api")
	return netHttp.ListenAndServe(s.host+":"+s.port, s.rootRouter)
}

func preFlightAndCORSHandler(next netHttp.Handler) netHttp.Handler {
	return netHttp.HandlerFunc(func(w netHttp.ResponseWriter, r *netHttp.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Origin, X-Requested-With, Content-Type, Accept, Access-Control-Allow-Origin")
		w.WriteHeader(netHttp.StatusOK)
		if r.Method == netHttp.MethodPost {
			next.ServeHTTP(w, r)
		}
	})
}
