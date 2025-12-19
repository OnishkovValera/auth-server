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
	RotateTokenWithTransaction(ctx context.Context, tokenToRevoke int, newToken entity.RefreshToken) error
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
	refreshToken := &entity.RefreshToken{}

	err := p.db.GetContext(ctx, refreshToken, `SELECT * FROM refresh_tokens WHERE token_hash = $1 AND revoked <> true`, hash)

	if err != nil {
		return nil, err
	}

	return refreshToken, nil
}

func (p PostgresRefreshTokenRepository) RotateTokenWithTransaction(ctx context.Context, tokenToRevoke int, newToken entity.RefreshToken) error {
	tx, err := p.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	_, err = tx.ExecContext(ctx, `
		UPDATE refresh_tokens
		SET revoked = true
		WHERE id = $1 AND revoked = false
	`, tokenToRevoke)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO refresh_tokens (
			user_id, token_hash, expired_at, revoked
		) VALUES ($1, $2, $3, $4)
	`, newToken.UserId, newToken.TokenHash, newToken.ExpiredAt, newToken.Revoked)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}
