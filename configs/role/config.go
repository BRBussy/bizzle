package role

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
	viper.SetDefault("ServerPort", "8181")
	viper.SetDefault("MongoDBHosts", []string{"localhost:27017"})
	viper.SetDefault("MongoDBName", "bizzle")
	viper.SetDefault("PreSharedSecret", "1234")
	viper.SetDefault("AuthURL", "http://localhost:8180/api")
}

type Config struct {
	ServerPort              string
	MongoDBConnectionString string
	MongoDBHosts            []string
	MongoDBName             string
	MongoDBUsername         string
	MongoDBPassword         string
	PreSharedSecret         string
	AuthURL                 string
}

func GetConfig(configFileName *string) (*Config, error) {
	// set places to look for config file
	viper.AddConfigPath("configs/role")
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
