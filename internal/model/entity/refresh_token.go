package entity

import "time"

type RefreshToken struct {
	Id        int       `db:"id"`
	UserId    int       `db:"user_id"`
	TokenHash string    `db:"token_hash"`
	ExpiredAt time.Time `db:"expired_at"`
	Revoked   bool      `db:"revoked"`
}

func NewRefreshToken(userId int, tokenHash string, expiredAt time.Time) *RefreshToken {
	return &RefreshToken{UserId: userId, TokenHash: tokenHash, ExpiredAt: expiredAt, Revoked: false}
}
