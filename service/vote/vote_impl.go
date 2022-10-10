package vote

import (
	dto "GoVoteApi/DTO"
	"GoVoteApi/models"
	"GoVoteApi/repository/postgres"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gosimple/slug"
)

var (
	ErrNoResult       = errors.New("no result found")
	ErrAlreadyVoted   = errors.New("vote type is quiz, you already voted")
	ErrCantUpdate     = errors.New("cant update vote")
	ErrCantDeleteVote = errors.New("you cant your delete vote in quiz mode")
)

// CreateVote implements service.VoteService
func (v *voteImpl) CreateVote(ctx context.Context, req *dto.VoteRequest) (*dto.VoteResponse, error) {
	var err error
	if err = v.validate.Struct(req); err != nil {
		return nil, err
	}

	vo := make([]models.VoteOptions, len(req.VoteOptions))
	for i := 0; i < len(req.VoteOptions); i++ {
		vo[i].Option = req.VoteOptions[i].Option
		vo[i].IsQuizAnswer = req.VoteOptions[i].IsQuizAnswer
	}

	if req.Slug == "" {
		req.Slug = slug.Make(req.Title)
	}

	vote := &models.Vote{
		UserID:      req.UserID,
		Title:       req.Title,
		Slug:        req.Slug,
		Type:        models.VoteType(req.Type),
		VoteOptions: vo,
	}

	id, err := v.repo.CreateVote(ctx, vote)
	if err != nil {
		return nil, err
	}

	return &dto.VoteResponse{
		ID:     id,
		Status: dto.StatusCreated,
		VoteID: id,
		Title:  req.Title,
		Slug:   req.Slug,
	}, nil
}

// GetAllVotes implements service.VoteService
func (v *voteImpl) GetAllVotes(ctx context.Context, limit int, offset int, status dto.VoteStatus) ([]dto.VoteResponse, error) {
	if limit > 100 || limit < 0 {
		limit = 10
	}

	votes, err := v.repo.GetAllVotesInfo(ctx, models.VoteStatus(status), limit, offset)
	if err != nil {
		return nil, err
	}
	if len(votes) == 0 {
		return nil, ErrNoResult
	}

	vr := make([]dto.VoteResponse, len(votes))
	for i := 0; i < len(votes); i++ {
		vr[i].ID = votes[i].ID
		vr[i].UserID = votes[i].UserID
		vr[i].VoteID = votes[i].ID
		vr[i].Slug = votes[i].Slug
		vr[i].Title = votes[i].Title
		vr[i].VoteStatus = dto.VoteStatus(votes[i].Status)
		vr[i].Type = dto.VoteType(votes[i].Type)
		vr[i].Status = dto.StatusFound
	}

	return vr, nil
}

// GetVoteByID implements service.VoteService
func (v *voteImpl) GetVoteByID(ctx context.Context, id uint) (*dto.VoteResponse, error) {
	vote, err := v.repo.GetVoteInfo(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoResult
		}
		return nil, err
	}
	vo, err := v.repo.GetVoteOptions(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoResult
		}
		return nil, err
	}
	uv, err := v.repo.GetVoters(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoResult
		}
		return nil, err
	}

	fmt.Printf("len of voteOptions: %d\n", len(vo))
	dvo := make([]dto.VoteOptions, len(vo))
	for i := 0; i < len(vo); i++ {
		dvo[i].ID = vo[i].ID
		dvo[i].Option = vo[i].Option
		dvo[i].Count = vo[i].Count
		if vote.Type == models.VoteQuiz || vote.Type == models.VoteQuizAnon {
			dvo[i].IsQuizAnswer = vo[i].IsQuizAnswer
		}
	}

	fmt.Printf("len of voters: %d\n", len(uv))
	var voters []dto.Voters
	if vote.Type == models.VoteMulti || vote.Type == models.VoteQuiz {
		voters = make([]dto.Voters, len(uv))
		for i := 0; i < len(uv); i++ {
			voters[i].ID = uv[i].ID
			voters[i].VoteID = uv[i].VoteID
			voters[i].OptionID = uv[i].OptionID
			voters[i].UserID = uv[i].UserID
			voters[i].Username = uv[i].Username
		}
	} else {
		voters = nil
	}

	resp := &dto.VoteResponse{
		ID:          vote.ID,
		Status:      dto.StatusFound,
		VoteID:      vote.ID,
		UserID:      vote.UserID,
		Title:       vote.Title,
		Slug:        vote.Slug,
		Type:        dto.VoteType(vote.Type),
		VoteStatus:  dto.VoteStatus(vote.Type),
		VoteOptions: dvo,
		Voters:      voters,
	}

	return resp, nil
}

// GetVoteBySlug implements service.VoteService
func (v *voteImpl) GetVoteBySlug(ctx context.Context, slug string) (*dto.VoteResponse, error) {
	vote, err := v.repo.GetVoteInfoBySlug(ctx, slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoResult
		}
		return nil, err
	}
	vo, err := v.repo.GetVoteOptions(ctx, vote.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoResult
		}
		return nil, err
	}
	uv, err := v.repo.GetVoters(ctx, vote.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoResult
		}
		return nil, err
	}

	dvo := make([]dto.VoteOptions, len(vo))
	for i := 0; i < len(vo); i++ {
		dvo[i].ID = vo[i].ID
		dvo[i].Option = vo[i].Option
		dvo[i].Count = vo[i].Count
		if vote.Type == models.VoteQuiz {
			dvo[i].IsQuizAnswer = vo[i].IsQuizAnswer
		}
	}

	var voters []dto.Voters
	if vote.Type == models.VoteMulti || vote.Type == models.VoteQuiz {
		voters = make([]dto.Voters, len(uv))
		for i := 0; i < len(uv); i++ {
			voters[i].ID = uv[i].ID
			voters[i].VoteID = uv[i].VoteID
			voters[i].OptionID = uv[i].OptionID
			voters[i].UserID = uv[i].UserID
			voters[i].Username = uv[i].Username
		}
	} else {
		voters = nil
	}

	resp := &dto.VoteResponse{
		ID:          vote.ID,
		Status:      dto.StatusFound,
		VoteID:      vote.ID,
		UserID:      vote.UserID,
		Title:       vote.Title,
		Slug:        vote.Slug,
		Type:        dto.VoteType(vote.Type),
		VoteStatus:  dto.VoteStatus(vote.Type),
		VoteOptions: dvo,
		Voters:      voters,
	}

	return resp, nil
}

func (v *voteImpl) AddVoteSlug(ctx context.Context, slug string, vote *dto.Voters) (*dto.VoteResponse, error) {
	if err := v.validate.Struct(vote); err != nil {
		return nil, err
	}

	voteInfo, err := v.repo.GetVoteInfoBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	var userVote *models.UserVotes
	uv, err := v.repo.GetUserVote(ctx, vote.VoteID, vote.UserID)
	if !errors.Is(err, postgres.ErrNoUserVoteFound) {
		return nil, err
	}

	if uv != nil {
		// if vote type is quiz or anonymous quiz
		if voteInfo.Type == models.VoteQuiz || voteInfo.Type == models.VoteQuizAnon {
			return nil, ErrAlreadyVoted
		}
		uv.OptionID = vote.OptionID
		userVote = uv
	} else {
		userVote = &models.UserVotes{
			ID:       vote.ID,
			UserID:   vote.UserID,
			Username: vote.Username,
			VoteID:   vote.VoteID,
			OptionID: vote.OptionID,
		}
	}

	if err = v.repo.AddUserVote(ctx, userVote); err != nil {
		return nil, err
	}

	return &dto.VoteResponse{
		Status: dto.StatusCreated,
	}, nil
}

func (v *voteImpl) AddVote(ctx context.Context, vote *dto.Voters) (*dto.VoteResponse, error) {
	if err := v.validate.Struct(vote); err != nil {
		return nil, err
	}

	voteInfo, err := v.repo.GetVoteInfo(ctx, vote.VoteID)
	if err != nil {
		return nil, err
	}

	var userVote *models.UserVotes
	uv, err := v.repo.GetUserVote(ctx, vote.VoteID, vote.UserID)
	if !errors.Is(err, postgres.ErrNoUserVoteFound) {
		return nil, err
	}

	// if user already voted...
	if uv != nil {
		// if vote type is quiz or anonymous quiz
		if voteInfo.Type == models.VoteQuiz || voteInfo.Type == models.VoteQuizAnon {
			return nil, ErrAlreadyVoted
		}
		uv.OptionID = vote.OptionID
		userVote = uv
	} else {
		userVote = &models.UserVotes{
			ID:       vote.ID,
			UserID:   vote.UserID,
			Username: vote.Username,
			VoteID:   vote.VoteID,
			OptionID: vote.OptionID,
		}
	}

	if err = v.repo.AddUserVote(ctx, userVote); err != nil {
		return nil, err
	}

	return &dto.VoteResponse{
		Status: dto.StatusCreated,
	}, nil
}

func (v *voteImpl) UpdateUserVote(ctx context.Context, vote *dto.Voters) (*dto.VoteResponse, error) {
	if err := v.validate.Struct(vote); err != nil {
		return nil, err
	}

	voteInfo, err := v.repo.GetVoteInfo(ctx, vote.ID)
	if err != nil {
		return nil, err
	}

	var userVote *models.UserVotes
	uv, err := v.repo.GetUserVote(ctx, vote.VoteID, vote.UserID)
	if !errors.Is(err, postgres.ErrNoUserVoteFound) {
		return nil, err
	}

	if uv != nil {
		if voteInfo.Type == models.VoteQuiz || voteInfo.Type == models.VoteQuizAnon {
			return nil, ErrCantUpdate
		}
		uv.OptionID = vote.OptionID
		userVote = uv
	} else {
		userVote = &models.UserVotes{
			ID:       uv.ID,
			UserID:   vote.UserID,
			Username: vote.Username,
			VoteID:   vote.VoteID,
			OptionID: vote.OptionID,
		}
	}

	if err = v.repo.UpdateUserVote(ctx, userVote); err != nil {
		return nil, err
	}

	return &dto.VoteResponse{
		Status: dto.StatusUpdated,
	}, nil
}

func (v *voteImpl) DeleteUserVote(ctx context.Context, vote *dto.Voters) (*dto.VoteResponse, error) {
	if err := v.validate.Struct(vote); err != nil {
		return nil, err
	}

	voteInfo, err := v.repo.GetVoteInfo(ctx, vote.ID)
	if err != nil {
		return nil, err
	}

	var userVote *models.UserVotes
	uv, err := v.repo.GetUserVote(ctx, vote.VoteID, vote.UserID)
	if !errors.Is(err, postgres.ErrNoUserVoteFound) {
		return nil, err
	}

	if uv != nil {
		if voteInfo.Type == models.VoteQuiz || voteInfo.Type == models.VoteQuizAnon {
			return nil, ErrCantDeleteVote
		}
		uv.OptionID = vote.OptionID
		userVote = uv
	} else {
		userVote = &models.UserVotes{
			ID: uv.ID,
		}
	}

	if err = v.repo.DeleteUserVote(ctx, userVote); err != nil {
		return nil, err
	}

	return &dto.VoteResponse{
		Status: dto.StatusDeleted,
	}, nil
}
