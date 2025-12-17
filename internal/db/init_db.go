package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func CreateDB(pathToDB string) (*sqlx.DB, error) {
	connect, err := sqlx.Connect("postgres", pathToDB)
	if err != nil {
		return nil, err
	}
	return connect, nil
}
