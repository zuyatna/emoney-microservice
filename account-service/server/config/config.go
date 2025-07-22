package config

import (
	"errors"

	"github.com/spf13/viper"
)

type Config struct {
	GRPCPORT    string `mapstructure:"GRPC_PORT"`
	HTTPPort    string `mapstructure:"HTTP_PORT"`
	PostgresURL string `mapstructure:"POSTGRES_URL"`
	RedisURL    string `mapstructure:"REDIS_URL"`
	JWTSecret   string `mapstructure:"JWT_SECRET"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			return config, nil
		}
		return config, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}
	return config, nil
}
