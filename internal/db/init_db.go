package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func CreateDB(pathToDB string) (*sqlx.DB, error) {
	return sqlx.Connect("postgres", pathToDB)
}
