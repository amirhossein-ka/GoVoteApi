package service

import (
	dto "GoVoteApi/DTO"
	"GoVoteApi/models"
	"context"
)

type (
	Service interface {
		UserService
		VoteService
		AuthService
	}
	UserService interface {
		Login(ctx context.Context, user, pass string) (*dto.UserResponse, error)
		Register(ctx context.Context, ur *dto.UserRequest) (*dto.UserResponse, error)
		Delete(ctx context.Context, ur *dto.UserRequest) error
		Info(ctx context.Context, username string) (*dto.UserResponse, error)
	}
	VoteService interface {
	}

	AuthService interface {
		GenerateToken(id uint, username string, role models.UserRole) (string, error)
		ClaimsFromToken(token string) (any, error)
	}
)
