package config

import (
	"github.com/kelseyhightower/envconfig"
)

func ParseEnv(cfg *Config) error {
	if err := envconfig.Process("mongo", &cfg.DBMongo); err != nil {
		return err
	}
	if err := envconfig.Process("redis", &cfg.DBRedis); err != nil {
		return err
	}
	if err := envconfig.Process("log", &cfg.Log); err != nil {
		return err
	}

	return nil
}
