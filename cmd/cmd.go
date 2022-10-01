package cmd

import (
	"GoVoteApi/config"
	"GoVoteApi/controller/mux"
	"GoVoteApi/pkg/logger"
	"GoVoteApi/pkg/logger/zap"
	"GoVoteApi/repository/memcache"
	"GoVoteApi/repository/mongo"
	"GoVoteApi/service"
	"GoVoteApi/service/auth"
	"GoVoteApi/service/user"
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator/v10"
)

func Run(cfg *config.Config) error {
	log, err := zap.New(&cfg.Log)
	if err != nil {
		return err
	}
	log.Info(logger.LogData{Section: "Run", Message: "Bananana"})

	srv, err := getService(cfg)
	if err != nil {
		return err
	}

	rest := mux.New(srv, cfg)
	go func() {
		if err = rest.Start(":8000"); err != nil {
			panic(err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	if err = rest.Stop(); err != nil {
		panic(err)
	}
	return nil
}

func getService(cfg *config.Config) (service.Service, error) {
	mongo, err := mongo.New(context.Background(), &cfg.DBMongo)
	if err != nil {
		return nil, err
	}
	memcache := memcache.New()

	a := auth.New(cfg.Secrets)
	validate := validator.New()
	userService := user.New(mongo, validate, memcache, a, cfg)

	type srv struct {
		service.AuthService
		service.VoteService
		service.UserService
	}

	return srv{a, nil, userService}, nil

}
