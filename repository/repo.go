package repository

import (
	"GoVoteApi/models"
	"context"
)

type (
	UserRepo interface {
		CreateUser(ctx context.Context, u *models.User) (uint, error)
		GetUser(ctx context.Context, id uint) (*models.User, error)
		GetUserByUsername(ctx context.Context, username string) (*models.User, error)
		UpdateUser(ctx context.Context, u *models.User) error
		DeleteUser(ctx context.Context, id uint) error
	}

	VoteRepo interface {
		CreateVote(ctx context.Context, v *models.Vote) (uint, error)
		GetAllVotesInfo(ctx context.Context, status models.VoteStatus, limit, offset int) ([]models.Vote, error)
		GetVoteInfo(ctx context.Context, id uint) (*models.Vote, error)
		GetVoteInfoBySlug(ctx context.Context, slug string) (*models.Vote, error)
		GetVoteOptions(ctx context.Context, id uint) (v []models.VoteOptions, err error)

		AddUserVote(ctx context.Context, uv *models.UserVotes) error
		GetVoters(ctx context.Context, voteID uint) (uv []models.UserVotes, err error)
		GetUserVote(ctx context.Context, voteID, userID uint) (*models.UserVotes, error)
		UpdateUserVote(ctx context.Context, uv *models.UserVotes) error
		DeleteUserVote(ctx context.Context, uv *models.UserVotes) error
	}

	Repository interface {
		UserRepo
		VoteRepo
	}
	// Cache is key/value cache
	Cache interface {
		Get(key string) any
		Set(key string, val any)
	}
)
