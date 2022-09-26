package repository

import (
	"GoVoteApi/models"
	"context"
)

type (
	Repository interface {
		CreateUser(ctx context.Context, u *models.User) (string, error)
		GetUser(ctx context.Context, id string) (*models.User, error)
		UpdateUser(ctx context.Context, u *models.User) error
		DeleteUser(ctx context.Context, id string) error
	}
    // key/value cache
	Cache interface{
        Get(key string) any
        Set(key string, val any)
    }
)
