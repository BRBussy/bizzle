package basic

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	jsonRpcClient "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/client"
	"github.com/rs/zerolog/log"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"strings"
)

type client struct {
	url                  string
	additionalHeaderKeys map[string]string
}

func New(
	url string,
) jsonRpcClient.Client {
	return &client{
		url:                  url,
		additionalHeaderKeys: make(map[string]string),
	}
}

func (c *client) AddAdditionalHeaderEntry(key, value string) {
	c.additionalHeaderKeys[key] = value
}

func (c *client) Post(request *jsonRpcClient.Request) (*jsonRpcClient.Response, error) {
	// marshal the request message
	marshalledRequest, err := json.Marshal(*request)
	if err != nil {
		return nil, errors.New("error marshalling request " + err.Error())
	}

	// put the bytes of the marshalled request into a buffer
	body := bytes.NewBuffer(marshalledRequest)

	// build the post request
	postRequest, err := http.NewRequest("POST", fmt.Sprintf("%s", c.url), body)
	if err != nil {
		return nil, errors.New("error creating post request " + err.Error())
	}

	// set the required headers on the request
	postRequest.Header.Set("Content-Type", "application/json")
	postRequest.Header.Set("Access-Control-Allow-Origin", "*")

	// set additional headers on request
	for key, value := range c.additionalHeaderKeys {
		postRequest.Header.Set(key, value)
	}

	// create the http client
	httpClient := &http.Client{
		//Timeout: time.Second * 5,
	}

	// perform the request
	postResponse, err := httpClient.Do(postRequest)
	if err != nil {
		return nil, errors.New("error performing post request " + err.Error())
	}

	// read the body bytes of the response
	postResponseBytes, err := ioutil.ReadAll(postResponse.Body)
	defer func() {
		if err := postResponse.Body.Close(); err != nil {
			log.Error().Err(err).Msg("closing response body error")
		}
	}()
	if err != nil {
		return nil, errors.New("error reading post response body bytes " + err.Error())
	}

	// check for an rpc error
	if strings.Contains(string(postResponseBytes), "rpc: can't find service") {
		return nil, errors.New("rpc error: method not found")
	}

	// unmarshal the body into the response
	response := jsonRpcClient.Response{}
	err = json.Unmarshal(postResponseBytes, &response)
	if err != nil {
		return nil, errors.New("error unmarshalling response bytes into json rpc response: " + err.Error())
	}

	if response.Error != "" {
		return &response, errors.New("json rpc service error " + response.Error)
	}

	return &response, nil
}

func (c *client) JsonRpcRequest(method string, request, response interface{}) error {
	jsonRpcRequest := jsonRpcClient.NewRequest(uuid.NewV4().String(), method, [1]interface{}{request})

	jsonRpcResponse, err := c.Post(&jsonRpcRequest)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(jsonRpcResponse.Result, response); err != nil {
		return err
	}

	return nil
}
