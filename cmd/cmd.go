package cmd

import (
	"GoVoteApi/config"
	"GoVoteApi/pkg/logger"
	"GoVoteApi/pkg/logger/zap"
	"GoVoteApi/repository/memcache"
	"GoVoteApi/repository/mongo"
	"context"
)

func Run(cfg *config.Config) error {
	log, err := zap.New(&cfg.Log)
	if err != nil {
		return err
	}
	log.Info(logger.LogData{Section: "Run", Message: "Bananana"})

	mongo, err := mongo.New(context.Background(), &cfg.DBMongo)
	if err != nil {
		return err
	}
	memcache := memcache.New()

	_ = mongo
	_ = memcache

	return nil
}
