package dto

import "encoding/json"

type (
	VoteType   uint8
	VoteStatus uint8

	VoteRequest struct {
		UserID uint     `json:"user_id,omitempty"`
		Title  string   `json:"title" validate:"gt=6,lt=256,required"`
		Slug   string   `json:"slug"`
		Type   VoteType `json:"type" validate:"lte=3,required"`
		//VoteOptions []struct {
		//	Option       string `json:"option" validate:"required"`
		//	IsQuizAnswer bool   `json:"is_answer,omitempty"`
		//} `json:"vote_options" validate:"required"`
		VoteOptions []VoteOptions `json:"vote_options" validate:"required"`
	}

	VoteResponse struct {
		ID          uint          `json:"id,omitempty"`
		Status      Status        `json:"status"`
		VoteID      uint          `json:"vote_id,omitempty"`
		UserID      uint          `json:"creator_id,omitempty"`
		Title       string        `json:"title,omitempty"`
		Slug        string        `json:"slug,omitempty"`
		Type        VoteType      `json:"type,omitempty"`
		VoteStatus  VoteStatus    `json:"vote_status,omitempty"`
		VoteOptions []VoteOptions `json:"vote_options,omitempty"`
		Voters      []Voters      `json:"voters,omitempty"`
		Data        any           `json:"data,omitempty"`
	}

	VoteOptions struct {
		ID           uint   `json:"option_id"`
		Option       string `json:"option,omitempty"`
		Count        uint   `json:"count,omitempty"`
		IsQuizAnswer bool   `json:"is_answer,omitempty"`
	}

	Voters struct {
		ID       uint
		VoteID   uint
		UserID   uint
		OptionID uint `json:"option_id"`
		Username string
	}
)

const (
	_ VoteType = iota
	VoteQuiz
	VoteAnon
	VoteMulti
	VoteMultiAnon
	VoteQuizAnon
)

const (
	_ VoteStatus = iota
	VoteClose
	VoteOpen
)

func (r VoteResponse) String() string {
	res, _ := json.MarshalIndent(r, "", "    ")
	return string(res)
}

func (r Voters) String() string {
	res, _ := json.MarshalIndent(r, "", "    ")
	return string(res)
}
