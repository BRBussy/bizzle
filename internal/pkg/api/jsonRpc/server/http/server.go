package http

import (
	"fmt"
	jsonRpcServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	"github.com/BRBussy/bizzle/internal/pkg/cors"
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/gorilla/rpc"
	gorillaJson "github.com/gorilla/rpc/json"
	"github.com/rs/zerolog/log"
	netHttp "net/http"
	"time"
)

type server struct {
	path             string
	host             string
	port             string
	rpcServer        *rpc.Server
	rootRouter       *chi.Mux
	apiRouter        *chi.Mux
	serviceProviders map[jsonRpcServiceProvider.Name]jsonRpcServiceProvider.Provider
	middleware       []func(netHttp.Handler) netHttp.Handler
}

func New(
	path string,
	host string,
	port string,
	middleware []func(netHttp.Handler) netHttp.Handler,
) *server {
	// create a new server
	newServer := new(server)
	newServer.serviceProviders = make(map[jsonRpcServiceProvider.Name]jsonRpcServiceProvider.Provider)
	newServer.path = path
	newServer.host = host
	newServer.port = port

	// create new gorilla rpc server
	newServer.rpcServer = rpc.NewServer()
	newServer.rpcServer.RegisterCodec(cors.CodecWithCors([]string{"*"}, gorillaJson.NewCodec()), "application/json")

	// initialise middleware
	if middleware == nil {
		middleware = make([]func(netHttp.Handler) netHttp.Handler, 0)
	}

	// create chi root router and apply middleware
	newServer.rootRouter = chi.NewRouter()

	// create chi api router
	newServer.apiRouter = chi.NewRouter()

	// apply middleware to api router
	newServer.apiRouter.Use(
		chiMiddleware.Timeout(time.Second * 60),
	)
	for _, middleware := range middleware {
		newServer.apiRouter.Use(middleware)
	}

	// mount api router on root router
	newServer.rootRouter.Mount("/api", newServer.apiRouter)

	// handle post requests to api router
	newServer.apiRouter.Post("/", newServer.rpcServer.ServeHTTP)

	return newServer
}

func (s *server) Start() error {
	log.Info().Msg("staring http json rpc api server on: " + s.host + ":" + s.port)
	return netHttp.ListenAndServe(s.host+":"+s.port, s.rootRouter)
}

func (s *server) RegisterServiceProvider(serviceProvider jsonRpcServiceProvider.Provider) error {
	log.Info().Msg(fmt.Sprintf("register %s jsonrpc service", serviceProvider.Name()))
	s.serviceProviders[serviceProvider.Name()] = serviceProvider
	if err := s.rpcServer.RegisterService(serviceProvider, string(serviceProvider.Name())); err != nil {
		log.Error().Err(err).Msgf("registering service %s with json rpc http server", serviceProvider.Name())
		return err
	}
	return nil
}

func (s *server) RegisterBatchServiceProviders(serviceProviders []jsonRpcServiceProvider.Provider) error {
	for _, serviceProvider := range serviceProviders {
		if err := s.RegisterServiceProvider(serviceProvider); err != nil {
			return err
		}
	}
	return nil
}

func preFlightHandler(w netHttp.ResponseWriter, r *netHttp.Request) {
	w.Header().Set(
		"Access-Control-Allow-Origin",
		"*",
	)
	w.Header().Set(
		"Content-Type",
		"application/json",
	)
	w.Header().Set(
		"Access-Control-Allow-Headers",
		"Origin, X-Requested-With, Content-Type, Accept, Access-Control-Allow-Origin, Authorization",
	)
	w.WriteHeader(netHttp.StatusOK)
}
