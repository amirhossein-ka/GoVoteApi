package dto

type (
	Status string

	UserRequest struct {
		ID       string `json:"id"`
		Username string `json:"username" validate:"required,min=6"`
		Fullname string `json:"fullname" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}

	UserResponse struct {
		Status Status `json:"status"`
		ID     string `json:"id,omitempty"`
		Token  string `json:"token,omitempty"`
		Data   any    `json:"data,omitempty"`
	}
)

const (
	StatusFound   Status = "found"
	StatusCreated Status = "created"
	StatusError   Status = "error"
)
