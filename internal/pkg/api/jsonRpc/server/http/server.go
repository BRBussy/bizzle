package http

import (
	"fmt"
	jsonRpcServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	"github.com/BRBussy/bizzle/internal/pkg/cors"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	gorillaJson "github.com/gorilla/rpc/json"
	"github.com/rs/zerolog/log"
	netHttp "net/http"
)

type server struct {
	path             string
	host             string
	port             string
	rpcServer        *rpc.Server
	serverMux        *mux.Router
	serviceProviders map[jsonRpcServiceProvider.Name]jsonRpcServiceProvider.Provider
	middleware       []mux.MiddlewareFunc
}

func New(
	path string,
	host string,
	port string,
	middleware []mux.MiddlewareFunc,
) *server {
	// create new gorilla mux rpc server
	rpcServer := rpc.NewServer()
	rpcServer.RegisterCodec(cors.CodecWithCors([]string{"*"}, gorillaJson.NewCodec()), "application/json")

	return &server{
		path:             path,
		host:             host,
		port:             port,
		serverMux:        mux.NewRouter(),
		rpcServer:        rpcServer,
		serviceProviders: make(map[jsonRpcServiceProvider.Name]jsonRpcServiceProvider.Provider),
		middleware:       middleware,
	}
}

func (s *server) Start() error {
	// set preflight handler for options method
	s.serverMux.Methods(netHttp.MethodOptions).HandlerFunc(preFlightHandler)

	// wrap rpc server with middleware
	var middlewareWrappedServer netHttp.Handler = s.rpcServer
	for _, middleware := range s.middleware {
		middlewareWrappedServer = middleware.Middleware(middlewareWrappedServer)
	}

	// start server
	s.serverMux.Handle(
		s.path,
		middlewareWrappedServer,
	).Methods(netHttp.MethodPost)
	if err := netHttp.ListenAndServe(s.host+":"+s.port, s.serverMux); err != nil {
		log.Error().Err(err).Msg("json rpc api server stopped")
	}
	return nil
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
