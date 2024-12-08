package application

import "github.com/kelseyhightower/envconfig"

type Config struct {
	DB struct {
		SQLite struct {
			Path string `envconfig:"DB_SQLITE_PATH" default:"./data/db/sqlite.db"`
		}
	}

	HTTP struct {
		Port uint `envconfig:"HTTP_PORT" default:"8080"`
	}
}

func NewConfig() (Config, error) {
	cfg := Config{}

	err := envconfig.Process("", &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
