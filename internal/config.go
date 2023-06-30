package internal

import (
	"github.com/caarlos0/env/v6"
	"github.com/sirupsen/logrus"
)

const (
	EnvDevelop    = "development"
	EnvProduction = "production"
)

// config is the struct that holds all the configuration
type config struct {
	AppConfig
	DBConfig
}

// Config is the global variable that holds the configuration for parse the environment variables
var Config config

// init is the function that parse the environment variables
func init() {
	if err := env.Parse(&Config); err != nil {
		logrus.Fatalf("Error initializing: %s", err.Error())
	}
}

// AppConfig is the struct that holds the configuration for the application
type AppConfig struct {
	AppPort string `env:"APP_PORT"`
	AppEnv  string `env:"APP_ENV"`
}
