package app

import (
	"context"
	"log/slog"
	grpcapp "sso/internal/app/grpc"
	storageapp "sso/internal/app/storage"
	"sso/internal/services/auth"
	"time"
)

type App struct {
	storageApp *storageapp.App
	gRPCServer *grpcapp.App
}

func New(ctx context.Context, log *slog.Logger, grpcPort int, connectionString string, tokenTTL time.Duration) *App {
	storageApp := storageapp.New(ctx, log, connectionString)

	authService := auth.New(log, storageApp.Storage, storageApp.Storage, storageApp.Storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		storageApp: storageApp,
		gRPCServer: grpcApp,
	}
}

func (a *App) MustRun() {
	a.gRPCServer.MustRun()
}

func (a *App) Stop() {
	a.gRPCServer.Stop()
	a.storageApp.MustStop()
}
