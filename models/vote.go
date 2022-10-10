package models

import "database/sql/driver"

type (
	VoteType   uint8
	VoteStatus uint8

	Vote struct {
		ID          uint          `json:"id,omitempty"`
		UserID      uint          `json:"user_id,omitempty"`
		Title       string        `json:"title"`
		Slug        string        `json:"slug"`
		Type        VoteType      `json:"type"`
		Status      VoteStatus    `json:"status"`
		Votes       []UserVotes   `json:"user_ids"`
		VoteOptions []VoteOptions `json:"options"`
	}

	VoteOptions struct {
		ID           uint
		VoteID       uint
		Option       string `json:"option"`
		Count        uint   `json:"count"`
		IsQuizAnswer bool   `json:"is_answer"`
	}

	UserVotes struct {
		ID       uint
		UserID   uint
		Username string
		VoteID   uint
		OptionID uint
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

func (v *VoteType) Scan(value any) error {
	*v = VoteType(value.(int64))
	return nil
}

func (v *VoteType) Value() (driver.Value, error) {
	return int64(*v), nil
}

const (
	_ VoteStatus = iota
	VoteClose
	VoteOpen
)

func (v *VoteStatus) Scan(value any) error {
	*v = VoteStatus(value.(int64))
	return nil
}

func (v *VoteStatus) Value() (driver.Value, error) {
	return int64(*v), nil
}
