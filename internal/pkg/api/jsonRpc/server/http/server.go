package http

import (
	jsonRpcServerAuthoriser "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/server/authoriser"
	jsonRpcServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
)

type server struct {
	path             string
	host             string
	port             string
	rpcServer        *rpc.Server
	authoriser       jsonRpcServerAuthoriser.Authoriser
	serverMux        *mux.Router
	serviceProviders map[jsonRpcServiceProvider.Name]jsonRpcServiceProvider.Provider
}
