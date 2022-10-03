package postgres

import (
	"GoVoteApi/config"
	"GoVoteApi/repository"
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type psql struct {
	db *sql.DB
}


func New(cfg config.DBPostgres) (repository.Repository, error) {
	db, err := sql.Open("pgx", cfg.URI)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return &psql{db: db}, nil
}
