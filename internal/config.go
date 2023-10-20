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
	FirestoreConfig
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
	AppPort          string `env:"APP_PORT"`
	AppEnv           string `env:"APP_ENV"`
	TokenExpireHours string `env:"TOKEN_EXPIRE_HOURS"`
	TokenSecret      string `env:"TOKEN_SECRET"`
	GoogleConfig     string `env:"GOOGLE_CONFIG"`
}

// FirestoreConfig is the struct that holds the configuration for the firestore
type FirestoreConfig struct {
	AuthB64   string `env:"FIRESTORE_AUTH_BASE64"`
	ProjectID string `env:"FIRESTORE_PROJECT_ID"`
}
