package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type (
	UserRole uint8

	User struct {
		ID       primitive.ObjectID `bson:"_id,omitempty"`
		FullName string             `bson:"fullname"`
		UserName string             `bson:"username"`
		Email    string             `bson:"email"`
		UserRole UserRole           `bson:"role"`
	}
)

const (
	_ UserRole = iota
	NormalUser
	AdminUser
)
