package models

import "database/sql/driver"

// TODO
// 1. complete vote repo,
// 2. mayyybe change to psql
// 3. idk,

type (
	VoteType   uint8
	VoteStatus uint8

	Vote struct {
		ID          uint          `json:"id,omitempty"`
		Title       string        `json:"title"`
		Slug        string        `json:"slug"`
		Type        VoteType      `json:"type"`
		Status      VoteStatus    `json:"status"`
		UserIDs     []string      `json:"user_ids"`
		VoteOptions []VoteOptions `json:"options"`
	}

	VoteOptions struct {
		Option       string `json:"option"`
		Count        uint   `json:"cout"`
		IsQuizAnswer bool   `json:"is_answer"`
	}
)

const (
	_ VoteType = iota
	VoteQuiz
	VoteAnon
	VoteMulti
)

func (v *VoteType) Scan(value any) error {
	*v = VoteType(value.(uint8))
	return nil
}

func (v VoteType) Value() (driver.Value, error) {
	return int64(v), nil
}

const (
	_ VoteStatus = iota
	VoteClose
	VoteOpen
)
