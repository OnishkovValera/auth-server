package repository

import (
	"auth-server/internal/model/entity"
	"context"

	"github.com/jmoiron/sqlx"
)

type RefreshTokenRepository interface {
	CreateToken(ctx context.Context, token *entity.RefreshToken) error
	UpdateToken(ctx context.Context, id *entity.RefreshToken) error
	DeleteToken(ctx context.Context, id int) error
	GetToken(ctx context.Context, id int) (*entity.RefreshToken, error)
	GetTokenByHash(ctx context.Context, hash string) (*entity.RefreshToken, error)
}

type PostgresRefreshTokenRepository struct {
	db *sqlx.DB
}

func NewRefreshTokenRepository(db *sqlx.DB) *PostgresRefreshTokenRepository {
	return &PostgresRefreshTokenRepository{db: db}
}

func (p PostgresRefreshTokenRepository) CreateToken(ctx context.Context, token *entity.RefreshToken) error {
	_, err := p.db.NamedExecContext(ctx, `
								INSERT INTO refresh_tokens (user_id, token_hash, expired_at, revoked)
								VALUES (:user_id, :token_hash, :expired_at, :revoked)
								`, token)
	return err
}

func (p PostgresRefreshTokenRepository) UpdateToken(ctx context.Context, id *entity.RefreshToken) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresRefreshTokenRepository) DeleteToken(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresRefreshTokenRepository) GetToken(ctx context.Context, id int) (*entity.RefreshToken, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresRefreshTokenRepository) GetTokenByHash(ctx context.Context, hash string) (*entity.RefreshToken, error) {
	//TODO implement me
	panic("implement me")
}
