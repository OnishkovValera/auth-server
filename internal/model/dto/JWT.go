package dto

type JWTDto struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func NewJWTDto(accessToken string, refreshToken string) *JWTDto {
	return &JWTDto{accessToken, refreshToken}
}
