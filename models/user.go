package models

import "database/sql/driver"

type (
	UserRole uint8

	User struct {
		ID       uint     `json:"id,omitempty"`
		FullName string   `json:"fullname"`
		UserName string   `json:"username"`
		Email    string   `json:"email"`
		Password string   `json:"password,omitempty"`
		UserRole UserRole `json:"role,omitempty"`
	}
)

const (
	_ UserRole = iota
	NormalUser
	AdminUser
)

// Scan implements sql.Scanner
func (u *UserRole) Scan(value any) error {
	*u = UserRole(value.(int64))
	return nil
}

func (u *UserRole) Value() (driver.Value, error) {
	return int64(*u), nil
}
