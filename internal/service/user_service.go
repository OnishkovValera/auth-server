package service

import (
	"auth-server/internal/model/entity"
	"auth-server/internal/repository"
	"context"
	"fmt"
)

type UserService struct {
	r repository.UserRepository
}

func NewUserService(r repository.UserRepository) *UserService {
	return &UserService{r: r}
}

func (s *UserService) CreateUser(ctx context.Context, user *entity.User) error {
	err := s.r.CreateUser(ctx, user)
	if err != nil {
		fmt.Println(err)
	}
	return err
}
