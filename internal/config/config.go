// Package config Description: This package contains the configuration for the ip2country service.
package config

import (
	"fmt"
	"log/slog"
	"reflect"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

type DatabaseType string

const (
	defaultLogLevel                     = "info"
	defaultServiceName                  = "ip2country"
	defaultServiceVersion               = "0.0.1"
	defaultActiveDataStore              = Local
	logLevel                            = "LOG_LEVEL"
	serviceName                         = "SERVICE_NAME"
	serviceVersion                      = "SERVICE_VERSION"
	activeDataStore                     = "ACTIVE_DATA_STORE"
	isDebug                             = "IP2COUNTRY_DEBUG"
	port                                = "PORT"
	rateLimit                           = "RATE_LIMIT"
	burstLimit                          = "BURST_LIMIT"
	configLogPrefix                     = "[Config]"
	Local                  DatabaseType = "local"
	API                    DatabaseType = "api"
	Relational             DatabaseType = "some_relational_db"
)

type Config struct {
	DB              []dbConfig
	Logger          loggerConfig
	ActiveDataStore DatabaseType `mapstructure:"ACTIVE_DATA_STORE"`
	RateLimit       int          `mapstructure:"RATE_LIMIT"`
	BurstLimit      int          `mapstructure:"BURST_LIMIT"`
	Port            int
	IsDebug         bool
}

type dbConfig struct {
	Host string
	Name DatabaseType
}

type loggerConfig struct {
	Level          string
	ServiceName    string
	ServiceVersion string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // or viper.SetConfigType("YAML")
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	viper.SetDefault(logLevel, defaultLogLevel)
	viper.SetDefault(serviceName, defaultServiceName)
	viper.SetDefault(serviceVersion, defaultServiceVersion)
	viper.SetDefault(activeDataStore, string(defaultActiveDataStore))
	viper.SetDefault(port, 8080)
	viper.SetDefault(isDebug, false)
	viper.SetDefault(rateLimit, 1)
	viper.SetDefault(burstLimit, 5)
	//viper.AutomaticEnv()

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	viper.AutomaticEnv() // read in environment variables that match
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &cfg, nil

	//return &Config{
	//	ActiveDataStore: DatabaseType(viper.GetString(activeDataStore)),
	//	IsDebug:         viper.GetBool(isDebug),
	//	DB: []dbConfig{
	//		{
	//			Host: "db/geolite2.zip",
	//			Name: defaultActiveDataStore,
	//		},
	//	},
	//	Logger: loggerConfig{
	//		Level:          viper.GetString(logLevel),
	//		ServiceName:    viper.GetString(serviceName),
	//		ServiceVersion: viper.GetString(serviceVersion),
	//	},
	//}, nil
}

func PrintConfigToLog(cfg interface{}, prefix string) {
	fields := reflect.TypeOf(cfg)
	values := reflect.ValueOf(cfg)
	filKind := fields.Kind()
	if filKind == reflect.Ptr {
		fields = fields.Elem()
		values = values.Elem()
		filKind = fields.Kind()
	}
	if filKind != reflect.Struct {
		printSingle(values, prefix, "")
		return
	}
	num := fields.NumField()
	for i := range num {
		value := values.Field(i)
		field := fields.Field(i)
		printSingle(value, prefix, field.Name)
	}
}

func printSingle(value reflect.Value, prefix, fieldName string) {
	valKind := value.Kind()
	switch valKind {
	case reflect.Ptr:
		if !value.IsNil() {
			PrintConfigToLog(value.Interface(), fieldName+".")
		} else {
			slog.Info(fmt.Sprintf("%s[%s] = [nil]", configLogPrefix, prefix+fieldName))
		}
	case reflect.Slice:
		for i := range value.Len() {
			PrintConfigToLog(value.Index(i).Interface(), sprintfConfigSliceElement(prefix, fieldName, i))
		}
	case reflect.Struct:
		PrintConfigToLog(value.Interface(), fieldName+".")
	case reflect.String:
		value := fmt.Sprint(value.Interface())
		// value = sanitizeString(fieldName, value) // Create a sanitizer if and when we get values that need hiding
		slog.Info(fmt.Sprintf("%s[%s] = [%s]", configLogPrefix, prefix+fieldName, value))
	case reflect.Bool:
		slog.Info(fmt.Sprintf("%s[%s] = [%s]", configLogPrefix, prefix+fieldName, strconv.FormatBool(value.Interface().(bool))))
	case reflect.Int:
		slog.Info(fmt.Sprintf("%s[%s] = [%s]", configLogPrefix, prefix+fieldName, strconv.Itoa(value.Interface().(int))))
	case reflect.Uint:
		uint64Val := uint64(value.Interface().(uint))
		slog.Info(fmt.Sprintf("%s[%s] = [%s]", configLogPrefix, prefix+fieldName, strconv.FormatUint(uint64Val, 10)))
	case reflect.Int64:
		if isDuration(value) {
			slog.Info(fmt.Sprintf("%s[%s] = [%s]", configLogPrefix, prefix+fieldName, value.Interface().(time.Duration).String()))
		} else {
			slog.Info(fmt.Sprintf("%s[%s] = [%s]", configLogPrefix, prefix+fieldName, strconv.FormatInt(value.Interface().(int64), 10)))
		}
	}
}

func sprintfConfigSliceElement(prefix, fieldName string, i int) string {
	return prefix + fieldName + "[" + strconv.Itoa(i) + "]."
}

func isDuration(obj reflect.Value) bool {
	_, ok := obj.Interface().(time.Duration)
	return ok
}
