package user

import (
	dto "GoVoteApi/DTO"
	"GoVoteApi/models"
	"GoVoteApi/repository/mongo"
	"context"
	"errors"
	// "github.com/go-playground/validator/v10"
)

// const token_hash string = "bananana"

var ErrUsernameExists error = errors.New("given username already exists")

// Delete implements service.UserService
func (u *user_impl) Delete(ctx context.Context, ur *dto.UserRequest) error {
	u.validator.StructExcept(ur, "FullName", "Email", "Password")
	var id string
	if ur.ID != "" {
		id = ur.ID
	} else if len(ur.Username) >= 6 {
		user, err := u.repo.GetUserByUsername(ctx, ur.Username)
		if err != nil {
			return err
		}
		id = user.ID.Hex()
	}
	return u.repo.DeleteUser(ctx, id)
}

// Login implements service.UserService
func (u *user_impl) Login(ctx context.Context, username, pass string) (*dto.UserResponse, error) {
	user, err := u.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	var token string
	if comparePass(pass, user.Password) {
		token, err = u.auth.GenerateToken(user.ID.Hex(), user.UserName, user.UserRole)
		if err != nil {
			return nil, err
		}
	}

	return &dto.UserResponse{
		Status: dto.StatusFound,
		ID:     user.ID.Hex(),
		Token:  token,
	}, nil
}

// Register implements service.UserService
func (u *user_impl) Register(ctx context.Context, ur *dto.UserRequest) (*dto.UserResponse, error) {
	user, err := u.repo.GetUserByUsername(ctx, ur.Username)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	if user != nil {
		return nil, ErrUsernameExists
	}

	// validate data
	if err = u.validator.Struct(ur); err != nil {
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

func (u *user_impl) Info(ctx context.Context, username string) (*dto.UserResponse, error) {
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
