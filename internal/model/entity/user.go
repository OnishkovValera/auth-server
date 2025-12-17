package entity

import "time"

type User struct {
	Id           int       `db:"id"`
	Login        string    `db:"login"`
	Name         string    `db:"name"`
	Surname      string    `db:"surname"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
	IsActive     bool      `db:"is_active"`
}
