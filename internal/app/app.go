package app

import (
	"log/slog"
	grpcapp "payment/internal/app/grpc"
	webhookapp "payment/internal/app/webhook"
	"payment/internal/config"
	"payment/internal/lib/robokassa"
	confirmsrv "payment/internal/service/confirm"
	generatesrv "payment/internal/service/generate"
	"payment/internal/storage/postgres"
)

type App struct {
	GRPCServer *grpcapp.App
	Webhook    *webhookapp.App
	Storage    *postgres.Storage
}

func New(log *slog.Logger, config *config.Config) *App {

	// init db
	storage := postgres.New(config.Postgres)

	// init robokassa
	robokassa := robokassa.New(config.RoboKassa.MerchantLogin, config.RoboKassa.Password)

	// init gen service
	generateService := generatesrv.New(log, storage, storage, robokassa) // сюда передаем Storage структуру

	// init webhook service
	confirmService := confirmsrv.New(log, storage)

	// init grpc
	grpcApp := grpcapp.New(log, generateService, config.GRPC.Port)

	// init webhook
	webhookApp := webhookapp.New(confirmService, config.Webhook.Port)

	return &App{GRPCServer: grpcApp, Webhook: webhookApp, Storage: storage}
}
