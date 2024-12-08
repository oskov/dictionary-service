package application

import (
	"github.com/jmoiron/sqlx"
	"github.com/oskov/dictionary-service/internal/core"
)

type App struct {
	Core *core.Core
}

func NewApp(db *sqlx.DB) *App {
	return &App{
		Core: core.NewCore(db),
	}
}
