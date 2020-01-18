package exercise

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func init() {
	err := viper.BindEnv("ServerPort", "PORT")
	if err != nil {
		log.Fatal().Err(err).Msgf("binding environment variables to configuration keys")
	}

	// set default configuration
	viper.SetDefault("ServerPort", "8083")
	viper.SetDefault("MongoDbHosts", []string{"localhost:27017"})
	viper.SetDefault("MongoDbName", "bizzle")
	viper.SetDefault("AuthURL", "http://localhost:8080/api")
	viper.SetDefault("RoleURL", "http://localhost:8081/api")
	viper.SetDefault("UserURL", "http://localhost:8082/api")
	viper.SetDefault("PreSharedSecret", "1234")
}

type Config struct {
	ServerPort              string
	MongoDBConnectionString string
	MongoDbHosts            []string
	MongoDbName             string
	RoleURL                 string
	AuthURL                 string
	PreSharedSecret         string
}

func GetConfig(configFileName *string) (*Config, error) {
	// set places to look for config file
	viper.AddConfigPath("configs/exercise")
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
