package service

import (
	"auth-server/internal/model/dto"
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

func (s *UserService) Login(ctx context.Context, login *dto.LoginRequest) (*dto.JWTDto, error) {
	user, err := s.r.GetUserByLogin(ctx, login.Login)
	if err != nil {
		return nil, err
	}

	isCorrect := CheckPasswordHash(login.Password, user.PasswordHash)

	if !isCorrect {
		WrongPassword := errors.New("wrong password")
		return nil, WrongPassword
	}

	token, err := s.authService.CreateToken(ctx, login.Login)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return token, nil

}

func (s *UserService) GetDataFromUser(token string) (map[string]any, error) {
	return nil, nil
}
