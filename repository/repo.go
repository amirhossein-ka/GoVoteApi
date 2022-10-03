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
	}

	Repository interface {
		UserRepo
		VoteRepo
	}
	// key/value cache
	Cache interface {
		Get(key string) any
		Set(key string, val any)
	}
)
