package secrets

import (
	"encoding/json"
	"github.com/BRBussy/bizzle/internal/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Secrets struct {
	AuthenticatorURL string `json:"authenticatorURL"`
}

func GetSecrets(pathToSecretsFile string) (*Secrets, error) {
	// validate given path to secrets file
	pathToSecretsFile, err := filepath.Abs(pathToSecretsFile)
	if err != nil {
		return nil, ErrCannotFindSecretFile{
			Path: pathToSecretsFile,
			Reasons: []string{
				"could not convert to absolute path",
				err.Error(),
			},
		}
	}
	secretsFileInfo, err := os.Stat(pathToSecretsFile)
	if err != nil {
		return nil, ErrCannotFindSecretFile{
			Path: pathToSecretsFile,
			Reasons: []string{
				"could not get file info at path",
				err.Error(),
			},
		}
	}
	if secretsFileInfo.IsDir() {
		return nil, ErrCannotFindSecretFile{
			Path: pathToSecretsFile,
			Reasons: []string{
				"path is a directory",
			},
		}
	}
	if !strings.HasSuffix(pathToSecretsFile, ".json") {
		return nil, ErrCannotFindSecretFile{
			Path: pathToSecretsFile,
			Reasons: []string{
				"does not end in '.json'",
			},
		}
	}

	// open a file reader for the secrets file
	secretsFile, err := ioutil.ReadFile(pathToSecretsFile)
	if err != nil {
		return nil, errors.ErrUnexpected{
			Reasons: []string{
				"reading secrets file",
				err.Error(),
			},
		}
	}

	// parse the json file
	parsedSecrets := new(Secrets)
	if err := json.Unmarshal(secretsFile, parsedSecrets); err != nil {
		return nil, errors.ErrUnexpected{
			Reasons: []string{
				"unmarshalling secrets file",
				err.Error(),
			},
		}
	}

	return parsedSecrets, nil
}
