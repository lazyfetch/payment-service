package app

import (
	"log/slog"
	grpcapp "payment/internal/app/grpc"
	webhookapp "payment/internal/app/webhook"
	"payment/internal/lib/robokassa"
	generatesrv "payment/internal/service/grpc/generate"
	confirmsrv "payment/internal/service/webhook/confirm"
	"payment/internal/storage/postgres"
)

type App struct {
	GRPCServer *grpcapp.App
	Webhook    *webhookapp.App
}

func New(log *slog.Logger, webhookPort int, grpcPort int, login, password string) *App {

	// init db
	storage := postgres.New()

	// init robokassa
	robokassa := robokassa.New(login, password)

	// init gen service
	generateService := generatesrv.New(log, storage, robokassa) // сюда передаем Storage структуру

	// init webhook service
	confirmService := confirmsrv.New(log, storage)

	// init grpc
	grpcApp := grpcapp.New(generateService, grpcPort)

	// init webhook
	webhookApp := webhookapp.New(confirmService, webhookPort)

	return &App{GRPCServer: grpcApp, Webhook: webhookApp}
}
