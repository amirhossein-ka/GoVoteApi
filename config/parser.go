package config

import (
	"github.com/kelseyhightower/envconfig"
)

func ParseEnv(cfg *Config) error {
	if err := envconfig.Process("postgres", &cfg.DBPostgres); err != nil {
		return err
	}
	if err := envconfig.Process("redis", &cfg.DBRedis); err != nil {
		return err
	}
	if err := envconfig.Process("log", &cfg.Log); err != nil {
		return err
	}
	if err := envconfig.Process("server", &cfg.Server); err != nil {
		return err
	}
	if err := envconfig.Process("secrets", &cfg.Secrets); err != nil {
		return err
	}
	return nil
}
