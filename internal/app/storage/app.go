package storageapp

import (
	"context"
	"fmt"
	"log/slog"
	"sso/internal/storage/postgres"
)

type App struct {
	log     *slog.Logger
	Storage *postgres.Storage
}

func New(ctx context.Context, log *slog.Logger, connectionString string) *App {
	storage, err := postgres.New(ctx, connectionString)
	if err != nil {
		panic(err)
	}

	return &App{
		log:     log,
		Storage: storage,
	}
}

func (a *App) MustStop() {
	if err := a.Stop(); err != nil {
		panic(err)
	}
}

func (a *App) Stop() error {
	const op = "storageapp.Stop"

	a.log.With(slog.String("op", op)).Info("stopping storage")

	if err := a.Storage.Stop(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
