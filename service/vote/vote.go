package vote

import (
	"GoVoteApi/config"
	"GoVoteApi/repository"
	"GoVoteApi/service"

	"github.com/go-playground/validator/v10"
)

type voteImpl struct {
	repo     repository.Repository
	validate *validator.Validate
	cfg      *config.Config
}

func New(repo repository.Repository, v *validator.Validate, cfg *config.Config) service.VoteService {
	return &voteImpl{
		repo:     repo,
		validate: v,
		cfg:      cfg,
	}
}
