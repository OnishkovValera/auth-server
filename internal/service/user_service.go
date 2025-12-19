package service

import (
	"auth-server/internal/model/entity"
	"auth-server/internal/repository"
	"context"
	"errors"
	"fmt"
)

type UserService struct {
	r           repository.UserRepository
	authService *AuthorizationService
}

func NewUserService(r repository.UserRepository, authService *AuthorizationService) *UserService {
	return &UserService{r: r, authService: authService}
}

func (s *UserService) CreateUser(ctx context.Context, user *entity.User) error {
	var err error
	if userExists, err := s.r.GetUserByLogin(ctx, user.Login); err == nil || userExists != nil {
		LoginAlreadyUsed := errors.New("user already exists")
		return fmt.Errorf("service: %w", LoginAlreadyUsed)
	}
	err = s.r.CreateUser(ctx, user)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (s *UserService) GetUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	return s.r.GetUser(ctx, user.Id)
}

func (s *UserService) GetUserInfo(ctx context.Context, token string) (*entity.User, error) {
	claims, err := s.authService.GetClaimsFromToken(token)
	if err != nil {
		return nil, err
	}
	user, err := s.r.GetUserByLogin(ctx, claims.Login)
	if err != nil {
		return nil, err
	}
	return user, nil
}
