package dto

type JWTDto struct {
	Jwt string `json:"jwt"`
}

func NewJWTDto(token string) *JWTDto {
	return &JWTDto{token}
}
