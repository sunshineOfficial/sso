package app

import (
	"context"
	"log/slog"
	grpcapp "sso/internal/app/grpc"
	"sso/internal/services/auth"
	"sso/internal/storage/postgres"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, connectionString string, tokenTTL time.Duration) *App {
	storage, err := postgres.New(context.Background(), connectionString)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
