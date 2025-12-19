package repository

import (
	"auth-server/internal/model/entity"
	"context"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	UpdateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, userId int) error
	GetUser(ctx context.Context, userId int) (*entity.User, error)
	GetAllUsers(ctx context.Context) ([]*entity.User, error)
	GetUserByLogin(ctx context.Context, login string) (*entity.User, error)
}

type PostgresUserRepository struct {
	db *sqlx.DB
}

func NewPostgresUserRepository(db *sqlx.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	_, err := r.db.NamedExecContext(ctx, `
								INSERT INTO users (login, name, surname, password_hash, is_active)
								VALUES (:login, :name, :surname, :password_hash, :is_active)
								`, user)
	return err
}
func (r *PostgresUserRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	return nil
}

func (r *PostgresUserRepository) DeleteUser(ctx context.Context, userId int) error {
	return nil
}

func (r *PostgresUserRepository) GetUser(ctx context.Context, userId int) (*entity.User, error) {
	var user entity.User
	err := r.db.GetContext(ctx, &user, `SELECT * FROM users WHERE id=$1`, userId)
	return &user, err
}

func (r *PostgresUserRepository) GetAllUsers(ctx context.Context) ([]*entity.User, error) {
	return nil, nil
}

func (r *PostgresUserRepository) GetUserByLogin(ctx context.Context, login string) (*entity.User, error) {
	var user entity.User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE login=$1", login)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
