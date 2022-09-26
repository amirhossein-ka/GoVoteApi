package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	VoteType   uint8
	VoteStatus uint8

	Vote struct {
		ID          primitive.ObjectID
		Title       string        `bson:"title"`
		Type        VoteType      `bson:"type"`
		Status      VoteStatus    `bson:"status"`
		VoteOptions []VoteOptions `bson:"options"`
	}

	VoteOptions struct {
		Option       string `bson:"option"`
		Count        string `bson:"cout"`
		IsQuizAnswer bool   `bson:"is_answer"`
	}
)

const (
	_ VoteType = iota
	VoteQuiz
	VoteAnon
	VoteMulti
)

const (
	_ VoteStatus = iota
	VoteClose
	VoteOpen
)
