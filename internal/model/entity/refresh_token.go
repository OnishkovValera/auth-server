package entity

import "time"

type RefreshToken struct {
	id        int       `db:"id"`
	userId    int       `db:"user_id"`
	tokenHash string    `db:"token_hash"`
	expiredAt time.Time `db:"expired_at"`
	revoked   bool      `db:"revoked_at"`
}
