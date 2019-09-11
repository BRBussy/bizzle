package authenticated

import (
	"cloud.google.com/go/compute/metadata"
	"fmt"
	jsonRpcClient "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/client"
	basicClient "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/client/basic"
	"github.com/rs/zerolog/log"
)

type Client struct {
	jsonRpcClient.Client
}

func New(url string) (jsonRpcClient.Client, error) {
	// query the id_token with ?audience as the authentication serviceURL
	authenticationServiceToken, err := metadata.Get(fmt.Sprintf("/instance/service-accounts/default/identity?audience=%s", url))
	if err != nil {
		log.Error().Err(err).Msg("failed to get authentication service access token")
		return nil, err
	}

	// create a basic client
	bc := basicClient.New(url)

	// set additional authentication header
	bc.AddAdditionalHeaderEntry("Authorization", fmt.Sprintf("Bearer %s", authenticationServiceToken))

	return &Client{
		Client: bc,
	}, nil
}
