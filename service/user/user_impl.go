package user

import (
	dto "GoVoteApi/DTO"
	"GoVoteApi/models"
	"GoVoteApi/repository/postgres"

	// "GoVoteApi/repository/postgres"
	"context"
	"errors"
	"fmt"
	// "github.com/go-playground/validator/v10"
)

var ErrUsernameExists = errors.New("given username already exists")

// Delete implements service.UserService
func (u *userImpl) Delete(ctx context.Context, ur *dto.UserRequest) error {
	if err := u.validator.StructPartial(ur, "ID"); err != nil {
		return err
	}
	return u.repo.DeleteUser(ctx, ur.ID)
}

// Login implements service.UserService
func (u *userImpl) Login(ctx context.Context, username, pass string) (*dto.UserResponse, error) {
	user, err := u.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	var token string
	if comparePass(pass, user.Password) {
		token, err = u.auth.GenerateToken(user.ID, user.UserName, user.UserRole)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("invalid password")
	}

	return &dto.UserResponse{
		Status: dto.StatusFound,
		ID:     user.ID,
		Token:  token,
	}, nil
}

// Register implements service.UserService
func (u *userImpl) Register(ctx context.Context, ur *dto.UserRequest) (*dto.UserResponse, error) {
	user, err := u.repo.GetUserByUsername(ctx, ur.Username)
	if err != postgres.ErrNoUserFound {
		println("other error occurred")
		return nil, err
	}
	if user != nil {
		println("username exists, bad")
		return nil, ErrUsernameExists
	}

	//validate data
	if err := u.validator.Struct(ur); err != nil {
		// err.(validator.ValidationErrors)
		return nil, err
	}

	// hash password to save in database
	pass, err := hashPass(ur.Password)
	if err != nil {
		return nil, err
	}
	umodel := &models.User{
		FullName: ur.Fullname,
		UserName: ur.Username,
		Email:    ur.Email,
		Password: pass,
		UserRole: models.NormalUser,
	}

	// add user to repository(mongodb)
	uid, err := u.repo.CreateUser(ctx, umodel)
	if err != nil {
		return nil, err
	}

	token, err := u.auth.GenerateToken(uid, ur.Username, models.NormalUser)
	if err != nil {
		return nil, err
	}
	// generate jwt token
	for {
		// check if token exists in cache
		if t := u.cache.Get(token); t != nil {
			token, err = u.auth.GenerateToken(uid, ur.Username, models.NormalUser)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}

	return &dto.UserResponse{
		Status: dto.StatusCreated,
		ID:     uid,
		Token:  token,
	}, nil
}

func (u *userImpl) Info(ctx context.Context, username string) (*dto.UserResponse, error) {
	// if len(username) >= 6 {
	// 	return nil, fmt.Errorf("invalid username")
	// }
	user, err := u.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	user.Password = ""

	return &dto.UserResponse{
		Status: dto.StatusFound,
		Data:   user,
	}, nil
}
