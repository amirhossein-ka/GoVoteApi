package postgres

import (
	"GoVoteApi/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrNoUserFound = errors.New("no user found with given username")

// CreateUser implements repository.Repository
func (p *psql) CreateUser(ctx context.Context, u *models.User) (uint, error) {
	fmt.Printf("%#v\n\n", u)
	var id int
	err := p.db.QueryRowContext(ctx, "INSERT INTO users (username, fullname, password, email, role) VALUES ($1,$2,$3,$4,$5) RETURNING user_id",
		u.UserName, u.FullName, u.Password, u.Email, u.UserRole).Scan(&id)
	if err != nil {
		return 0, err
	}
	println(id)
	return uint(id), nil
}

// DeleteUser implements repository.Repository
func (p *psql) DeleteUser(ctx context.Context, id uint) error {
	_, err := p.db.ExecContext(ctx, "DELETE FROM users WHERE user_id=$1", id)
	return err
}

// GetUser implements repository.Repository
func (p *psql) GetUser(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := p.db.QueryRowContext(ctx, "SELECT user_id, username, fullname, email, role FROM users WHERE user_id=$1", id).
		Scan(&user.ID, &user.UserName, &user.FullName, &user.Email, &user.UserRole)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoUserFound
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername implements repository.Repository
func (p *psql) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := p.db.QueryRowContext(ctx, "SELECT user_id, username, fullname, email, password, role FROM users WHERE username=$1", username).
		Scan(&user.ID, &user.UserName, &user.FullName, &user.Email, &user.Password, &user.UserRole)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoUserFound
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUser implements repository.Repository
func (p *psql) UpdateUser(ctx context.Context, u *models.User) error {
	res, err := p.db.ExecContext(ctx, "UPDATE users SET username=$1,email=$2,fullname=$3,password=$4 WHERE user_id=$5",
		u.UserName, u.Email, u.FullName, u.Password, u.ID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return fmt.Errorf("no rows affected in database")
	}
	return nil
}
