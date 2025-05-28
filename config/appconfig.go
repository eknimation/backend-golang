package config

import (
	"context"

	"backend-service/internal/infrastructure/database"

	"github.com/sethvargo/go-envconfig"
)

var conf *AppConfig

type AppConfig struct {
	MongoDB   *database.MongoConfig
	Port      string `env:"API_PORT"`
	Env       string `env:"APP_ENV"`
	BasePath  string `env:"BASE_PATH"`
	JWTSecret string `env:"JWT_SECRET"`
}

func GetAppConfig() *AppConfig {
	var config AppConfig
	if err := envconfig.Process(context.Background(), &config); err != nil {
		panic(err)
	}

	conf = &config

	return conf
}
