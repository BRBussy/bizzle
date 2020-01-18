package client

type Client interface {
	JsonRpcRequest(method string, request, response interface{}) error
}
