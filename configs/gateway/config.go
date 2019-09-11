package gateway

import (
	"github.com/BRBussy/bizzle/internal/pkg/environment"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func init() {
	err := viper.BindEnv("ServerPort", "PORT")
	if err != nil {
		log.Fatal().Err(err).Msgf("binding environment variables to configuration keys")
	}

	// set default configuration
	viper.SetDefault("Environment", environment.Development)
	viper.SetDefault("ServerPort", "8080")
	viper.SetDefault("AuthenticatorURL", "http://localhost:8081")
	viper.SetDefault("ExerciseURL", "http://localhost:8082")
}

type Config struct {
	Environment      environment.Environment
	ServerPort       string
	AuthenticatorURL string
	ExerciseURL      string
}

func GetConfig(configFileName *string) (*Config, error) {
	// set places to look for config file
	viper.AddConfigPath("configs/gateway")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")

	// set the name of the config file
	viper.SetConfigName(*configFileName)
	if err := viper.ReadInConfig(); err != nil {
		log.Error().Err(err).Msgf("could not parse config file")
		return nil, err
	}

	// parse the config file
	cfg := new(Config)
	if err := viper.Unmarshal(cfg); err != nil {
		log.Error().Err(err).Msg("unmarshalling config file")
		return nil, err
	}

	return cfg, nil
}
