package dto

type LoginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

func NewLoginRequest() *LoginRequest {
	return &LoginRequest{}
}
