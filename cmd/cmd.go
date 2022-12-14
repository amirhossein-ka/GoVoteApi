package cmd

import (
	"GoVoteApi/config"
	"GoVoteApi/controller/mux"
	"GoVoteApi/pkg/logger"
	"GoVoteApi/pkg/logger/zap"
	"GoVoteApi/repository/memcache"
	"GoVoteApi/repository/postgres"
	"GoVoteApi/service"
	"GoVoteApi/service/auth"
	"GoVoteApi/service/user"
	"GoVoteApi/service/vote"
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
	//log.Info(logger.LogData{Section: "Run", Message: "Running"})

	srv, err := getService(cfg)
	if err != nil {
		return err
	}

	rest := mux.New(srv, cfg, log)
	go func() {
		log.Info(logger.LogData{Section: "cmd/Run", Message: "Starting Server"})
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
	repo, err := postgres.New(cfg.DBPostgres)
	if err != nil {
		return nil, err
	}
	cache := memcache.New()

	a := auth.New(cfg.Secrets)
	validate := validator.New()
	userService := user.New(repo, validate, cache, a, cfg)
	voteService := vote.New(repo, validate, cfg)

	type srv struct {
		service.AuthService
		service.VoteService
		service.UserService
	}

	return srv{a, voteService, userService}, nil

}
