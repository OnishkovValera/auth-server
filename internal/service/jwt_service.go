package service

import (
	"auth-server/internal/model/dto"
	"auth-server/internal/model/entity"
	"auth-server/internal/repository"
	"auth-server/internal/security"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthorizationService struct {
	jwtMaker *security.JwtMaker
	r        repository.RefreshTokenRepository
	userRepo repository.UserRepository
}

func NewAuthorizationService(maker *security.JwtMaker, r repository.RefreshTokenRepository, userRepo repository.UserRepository) *AuthorizationService {
	return &AuthorizationService{jwtMaker: maker, r: r, userRepo: userRepo}
}

func (h *AuthorizationService) CheckJWT(token string) (bool, error) {
	return h.jwtMaker.ParseJWT(token)
}

func (h *AuthorizationService) GetClaimsFromToken(token string) (jwt.Claims, error) {
	return h.jwtMaker.GetClaimsFromToken(token)
}

func (h *AuthorizationService) CreateToken(ctx context.Context, login string) (*dto.JWTDto, error) {
	accessToken, err := h.jwtMaker.CreateAccessToken(login)
	if err != nil {
		return nil, err
	}
	user, err := h.userRepo.GetUserByLogin(ctx, login)
	if err != nil {
		return nil, err
	}
	refreshToken := h.jwtMaker.CreateRefreshToken()
	refreshHash := sha256.Sum256([]byte(refreshToken))
	encodedHash := hex.EncodeToString(refreshHash[:])
	refreshEntity := entity.NewRefreshToken(user.Id, encodedHash, time.Now().Add(90*24*time.Hour))
	err = h.r.CreateToken(ctx, refreshEntity)
	return dto.NewJWTDto(accessToken, refreshToken), err
}

func (h *AuthorizationService) saveRefreshToken(refreshToken *entity.RefreshToken) error {
	return nil
}

func (h *AuthorizationService) RefreshToken(refreshToken string) (*dto.JWTDto, error) {
	return nil, nil
}

func (h *AuthorizationService) RotateToken(token string) (*dto.JWTDto, error) {
	return nil, nil
}
