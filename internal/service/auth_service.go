package service

import (
	"auth-server/internal/model/dto"
	"auth-server/internal/model/entity"
	"auth-server/internal/repository"
	"auth-server/internal/security"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

type AuthorizationService struct {
	jwtMaker *security.JwtMaker
	r        repository.RefreshTokenRepository
	userRepo repository.UserRepository
}

func NewAuthorizationService(maker *security.JwtMaker, r repository.RefreshTokenRepository, userRepo repository.UserRepository) *AuthorizationService {
	return &AuthorizationService{jwtMaker: maker, r: r, userRepo: userRepo}
}

func (h *AuthorizationService) Login(ctx context.Context, login *dto.LoginRequest) (*dto.JWTDto, error) {
	user, err := h.userRepo.GetUserByLogin(ctx, login.Login)
	if err != nil {
		return nil, err
	}

	isCorrect := CheckPasswordHash(login.Password, user.PasswordHash)

	if !isCorrect {
		WrongPassword := errors.New("wrong password")
		return nil, WrongPassword
	}

	token, err := h.CreateToken(ctx, login.Login)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return token, nil
}

func (h *AuthorizationService) CheckJWT(token string) (bool, error) {
	return h.jwtMaker.ParseJWT(token)
}

func (h *AuthorizationService) GetClaimsFromToken(token string) (*security.Claims, error) {
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

	refreshEntity := entity.NewRefreshToken(user.Id, h.getTokenHash(refreshToken), time.Now().Add(90*24*time.Hour))
	err = h.r.CreateToken(ctx, refreshEntity)
	return dto.NewJWTDto(accessToken, refreshToken), err
}

func (h *AuthorizationService) RefreshToken(ctx context.Context, refreshToken string) (*dto.JWTDto, error) {
	tokenHash := h.getTokenHash(refreshToken)
	token, err := h.r.GetTokenByHash(ctx, tokenHash)
	if err != nil {
		return nil, err
	}

	if token.ExpiredAt.Before(time.Now()) {
		TokenExpired := errors.New("token is expired")
		return nil, TokenExpired
	}

	user, err := h.userRepo.GetUser(ctx, token.UserId)
	if err != nil {
		return nil, err
	}

	accessToken, err := h.jwtMaker.CreateAccessToken(user.Login)
	if err != nil {
		return nil, err
	}
	refreshToken, err = h.rotateRefreshToken(ctx, token.Id, user.Id)

	return dto.NewJWTDto(accessToken, refreshToken), nil
}

func (h *AuthorizationService) rotateRefreshToken(ctx context.Context, tokenToRevoke int, userId int) (string, error) {
	tokenString := h.jwtMaker.CreateRefreshToken()
	refreshEntity := entity.NewRefreshToken(userId, h.getTokenHash(tokenString), time.Now().Add(90*24*time.Hour))

	err := h.r.RotateTokenWithTransaction(ctx, tokenToRevoke, *refreshEntity)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (h *AuthorizationService) getTokenHash(token string) string {
	refreshHash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(refreshHash[:])
}
