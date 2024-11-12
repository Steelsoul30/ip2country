// Package config Description: This package contains the configuration for the ip2country service.
package config

import (
	"github.com/spf13/viper"
)

const (
	defaultLogLevel        = "info"
	defaultServiceName     = "ip2country"
	defaultServiceVersion  = "0.0.1"
	defaultActiveDataStore = "local"
	logLevel               = "LOG_LEVEL"
	serviceName            = "SERVICE_NAME"
	serviceVersion         = "SERVICE_VERSION"
	activeDataStore        = "ACTIVE_DATA_STORE"
)

type Config struct {
	LogLevel        string
	DB              []dbConfig
	ActiveDataStore string
	ServiceVersion  string
	ServiceName     string
}

type dbConfig struct {
	Host string
	Name string
}

func LoadConfig() (*Config, error) {
	viper.SetDefault(logLevel, defaultLogLevel)
	viper.SetDefault(serviceName, defaultServiceName)
	viper.SetDefault(serviceVersion, defaultServiceVersion)
	viper.SetDefault(activeDataStore, defaultActiveDataStore)
	viper.AutomaticEnv()

	return &Config{
		LogLevel:        viper.GetString(logLevel),
		ServiceName:     viper.GetString(serviceName),
		ServiceVersion:  viper.GetString(serviceVersion),
		ActiveDataStore: viper.GetString(activeDataStore),
		DB: []dbConfig{
			{Host: "",
				Name: defaultActiveDataStore,
			},
		},
	}, nil
}
