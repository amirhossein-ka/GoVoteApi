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
		CreateVote(ctx context.Context, req *dto.VoteRequest) (*dto.VoteResponse, error)
		GetAllVotes(ctx context.Context, limit, offset int, status dto.VoteStatus) ([]dto.VoteResponse, error)
		GetVoteByID(ctx context.Context, id uint) (*dto.VoteResponse, error)
		GetVoteBySlug(ctx context.Context, slug string) (*dto.VoteResponse, error)
		AddVote(ctx context.Context, vote *dto.Voters) (*dto.VoteResponse, error)
		AddVoteSlug(ctx context.Context, slug string, vote *dto.Voters) (*dto.VoteResponse, error)
		UpdateUserVote(ctx context.Context, vote *dto.Voters) (*dto.VoteResponse, error)
	}

	AuthService interface {
		GenerateToken(id uint, username string, role models.UserRole) (string, error)
		ClaimsFromToken(token string) (any, error)
	}
)
