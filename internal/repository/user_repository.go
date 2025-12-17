package repository

import (
	"auth-server/internal/model/entity"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	UpdateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, userId int) error
	GetUser(ctx context.Context, userId int) (*entity.User, error)
	GetAllUsers(ctx context.Context) ([]*entity.User, error)
}

type PostgresUserRepository struct {
	db *sqlx.DB
}

func NewPostgresUserRepository(db *sqlx.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	result, err := r.db.NamedExecContext(ctx, `
								INSERT INTO users (login, name, surname, password_hash, is_active)
								VALUES (:login, :name, :surname, :password_hash, :is_active)
								`, &user)
	fmt.Println(result)
	return err
}
func (r *PostgresUserRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	return nil
}

func (r *PostgresUserRepository) DeleteUser(ctx context.Context, userId int) error {
	return nil
}

func (r *PostgresUserRepository) GetUser(ctx context.Context, userId int) (*entity.User, error) {
	return nil, nil
}

func (r *PostgresUserRepository) GetAllUsers(ctx context.Context) ([]*entity.User, error) {
	return nil, nil
}
