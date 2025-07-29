package config

import (
	"errors"
	"github.com/spf13/viper"
)

type Config struct {
	GRPCPort             string `mapstructure:"GRPC_PORT"`
	HTTPPort             string `mapstructure:"HTTP_PORT"`
	PostgresURL          string `mapstructure:"POSTGRES_URL"`
	ElasticsearchURL     string `mapstructure:"ELASTICSEARCH_URL"`
	JWTSecret            string `mapstructure:"JWT_SECRET"`
	AccountServiceTarget string `mapstructure:"ACCOUNT_SERVICE_TARGET"`
	RABBITMQURL          string `mapstructure:"RABBITMQ_URL"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return
		}
	}
	err = viper.Unmarshal(&config)
	return
}
