package postgres

import (
	"GoVoteApi/models"
	"context"
)

func (p *psql) CreateVote(ctx context.Context, v *models.Vote) (uint, error) {
	var id uint
	err := p.db.QueryRowContext(ctx, "INSERT INTO vote (title, slug, vote_type, vote_status) VALUES ($1,$2,$3,$4) RETURNING vote_id", v.Title, v.Slug, v.Type, v.Status).Scan(&id)
	if err != nil {
		return 0, nil
	}

    return id, nil
}

func (p *psql) GetVoteBySlug(ctx context.Context, slug string) (*models.Vote, error) {
    panic("not implemented")
}
