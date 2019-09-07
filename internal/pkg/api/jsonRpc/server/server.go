package server

import (
	jsonRpcServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
)

type Server interface {
	Start() error
	SecureStart() error
	RegisterServiceProvider(jsonRpcServiceProvider.Provider) error
	RegisterBatchServiceProviders([]jsonRpcServiceProvider.Provider) error
}
