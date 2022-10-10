package postgres

import (
	"GoVoteApi/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrNoUserVoteFound = errors.New("no vote found with on vote from user")

func (p *psql) CreateVote(ctx context.Context, v *models.Vote) (uint, error) {
	var id uint
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, "INSERT INTO votes (title, slug, user_id, vote_type, vote_status) VALUES ($1,$2,$3,$4,$5) RETURNING vote_id",
		v.Title, v.Slug, v.UserID, v.Type, v.Status).Scan(&id)
	if err != nil {
		return 0, err
	}

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO vote_options  (vote_id, option_name, all_votes, quiz_answer) VALUES ($1,$2,$3,$4)")
	if err != nil {
		return 0, err
	}

	for i := 0; i < len(v.VoteOptions); i++ {
		_, err = stmt.ExecContext(ctx, id, v.VoteOptions[i].Option, v.VoteOptions[i].Count, v.VoteOptions[i].IsQuizAnswer)
		if err != nil {
			//if err = tx.Rollback(); err != nil {
			//	return 0, err
			//}
			return 0, err
		}
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	println(id)
	return id, nil
}

func (p *psql) GetAllVotesInfo(ctx context.Context, status models.VoteStatus, limit, offset int) (v []models.Vote, err error) {
	rows, err := p.db.QueryContext(
		ctx,
		"SELECT vote_id,user_id, title, slug, vote_type,vote_status FROM votes WHERE vote_status=$1 LIMIT $2 OFFSET $3",
		status, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		e := rows.Close()
		if e != nil {
			err = e
		}
	}()

	for rows.Next() {
		var vote models.Vote
		err = rows.Scan(&vote.ID, &vote.UserID, &vote.Title, &vote.Slug, &vote.Type, &vote.Status)
		if err != nil {
			return nil, err
		}
		v = append(v, vote)
	}
	if rows.Err() != nil {
		return nil, err
	}
	return
}

func (p *psql) GetVoteInfo(ctx context.Context, id uint) (*models.Vote, error) {
	var v models.Vote
	err := p.db.QueryRowContext(
		ctx,
		`SELECT vote_id,user_id, title, slug, vote_type,vote_status FROM votes WHERE vote_id=$1`, id,
	).Scan(&v.ID, &v.UserID, &v.Title, &v.Slug, &v.Type, &v.Status)
	if err != nil {
		return nil, err
	}

	return &v, nil
}

func (p *psql) GetVoteInfoBySlug(ctx context.Context, slug string) (*models.Vote, error) {
	var v models.Vote
	err := p.db.QueryRowContext(
		ctx,
		`SELECT vote_id,user_id, title, slug, vote_type,vote_status FROM votes WHERE slug=$1`, slug,
	).Scan(&v.ID, &v.UserID, &v.Title, &v.Slug, &v.Type, &v.Status)
	if err != nil {
		return nil, err
	}

	return &v, nil
}

func (p *psql) GetVoteOptions(ctx context.Context, id uint) (v []models.VoteOptions, err error) {
	rows, err := p.db.QueryContext(
		ctx,
		"SELECT opts_id, option_name, all_votes, quiz_answer FROM vote_options WHERE vote_id=$1", id,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		e := rows.Close()
		if e != nil {
			err = e
		}
	}()
	for i := 0; rows.Next(); i++ {
		var voteOpt models.VoteOptions
		if err = rows.Scan(&voteOpt.ID, &voteOpt.Option, &voteOpt.Count, &voteOpt.IsQuizAnswer); err != nil {
			return nil, err
		}
		v = append(v, voteOpt)
	}
	if rows.Err() != nil {
		return nil, err
	}

	return
}

func (p *psql) GetVoters(ctx context.Context, voteID uint) (uv []models.UserVotes, err error) {
	rows, err := p.db.QueryContext(
		ctx,
		"SELECT id,vote_id,user_username, opts_id FROM user_votes WHERE vote_id=$1", voteID,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		e := rows.Close()
		if e != nil {
			err = e
		}
	}()
	for i := 0; rows.Next(); i++ {
		var v models.UserVotes
		if err = rows.Scan(&v.ID, &v.VoteID, &v.Username, &v.OptionID); err != nil {
			return nil, err
		}
		uv = append(uv, v)
	}

	return
}

func (p *psql) AddUserVote(ctx context.Context, uv *models.UserVotes) error {
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
	if err != nil {
		return err
	}
	//defer tx.Rollback()

	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO user_votes (vote_id, user_id, opts_id, user_username) VALUES ($1,$2,$3, $4)",
		uv.VoteID, uv.UserID, uv.OptionID, uv.Username,
	)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	_, err = tx.ExecContext(
		ctx,
		"UPDATE vote_options SET all_votes=all_votes+1 WHERE vote_id=$1 AND opts_id=$2", uv.VoteID, uv.OptionID,
	)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (p *psql) GetUserVote(ctx context.Context, voteID, userID uint) (*models.UserVotes, error) {
	var uv models.UserVotes
	if err := p.db.QueryRowContext(
		ctx,
		"SELECT id, vote_id, user_id, user_username, opts_id FROM user_votes WHERE vote_id=$1 AND user_id=$2", voteID, userID,
	).Scan(&uv.ID, &uv.VoteID, &uv.UserID, &uv.Username, &uv.OptionID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoUserVoteFound
		}
		return nil, err
	}

	return &uv, nil
}

func (p *psql) UpdateUserVote(ctx context.Context, uv *models.UserVotes) error {
	// I guess there is no need for transaction in this method
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(
		ctx,
		"UPDATE user_votes SET opts_id=$1 WHERE id=$2 AND vote_id=$3 AND user_id=$4", uv.OptionID, uv.ID, uv.VoteID, uv.UserID,
	)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (p *psql) DeleteUserVote(ctx context.Context, uv *models.UserVotes) error {
	res, err := p.db.ExecContext(ctx, "DELETE FROM user_votes WHERE id=$1", uv.UserID)
	if err != nil {
		return err
	}
	if aff, err := res.RowsAffected(); err == nil {
		if aff < 0 {
			return fmt.Errorf("no rows deleted")
		}
	} else {
		return err
	}

	return nil
}
