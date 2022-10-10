package user

import (
	"GoVoteApi/config"
	"GoVoteApi/repository"
	"GoVoteApi/service"

	"github.com/go-playground/validator/v10"
)

type userImpl struct {
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
	return &userImpl{
		cfg:       cfg,
		repo:      r,
		validator: v,
		cache:     c,
		auth:      a,
	}
}
