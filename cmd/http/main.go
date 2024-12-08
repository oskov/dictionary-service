package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/oskov/dictionary-service/db/migrations"
	"github.com/oskov/dictionary-service/internal/api/http"
	"github.com/oskov/dictionary-service/internal/application"
	"golang.org/x/sync/errgroup"
)

func main() {
	if err := run(); err != nil {
		log.Printf("Application error: %v", err)
		return
	}
	log.Println("Application stopped")
}

func run() error {
	mainCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	cfg, err := application.NewConfig()
	if err != nil {
		return err
	}

	err = migrations.RunSQLLite(cfg.DB.SQLite.Path)
	if err != nil {
		return err
	}

	db, err := application.NewDB(cfg)
	if err != nil {
		return err
	}

	app := application.NewApp(db)

	server := http.NewServer(mainCtx, cfg.HTTP.Port, *app)

	g, gCtx := errgroup.WithContext(mainCtx)
	g.Go(func() error {
		return server.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		return server.Shutdown(context.Background())
	})

	log.Printf("Application started on port %d", cfg.HTTP.Port)

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}
