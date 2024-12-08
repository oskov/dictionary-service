package application

import (
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func NewDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite", cfg.DB.SQLite.Path)
	if err != nil {
		return nil, err
	}
	return db, nil
}
