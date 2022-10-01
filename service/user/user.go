package user

import (
	"GoVoteApi/config"
	"GoVoteApi/repository"
	"GoVoteApi/service"

	"github.com/go-playground/validator/v10"
)

type user_impl struct {
	repo      repository.Repository
	cache     repository.Cache
    validator *validator.Validate
	auth      service.AuthService
	cfg       *config.Config
}

func New(
	r repository.Repository,
	v *validator.Validate,
	c repository.Cache,
	a service.AuthService,
	cfg *config.Config,
) service.UserService {
	return &user_impl{
		cfg:       cfg,
		repo:      r,
		validator: v,
		cache:     c,
		auth:      a,
	}
}
