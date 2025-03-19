package app

import (
	"log/slog"
	grpcapp "payment/internal/app/grpc"
	webhookapp "payment/internal/app/webhook"
	paymentsrv "payment/internal/service/grpc/generate"
	// "payment/internal/storage/postgres"
)

type App struct {
	GRPCServer *grpcapp.App
	Webhook    *webhookapp.App
}

func New(log *slog.Logger, webHookPort int, grpcPort int) *App {

	//	storage := postgres.New()

	// init service layer
	paymentService := paymentsrv.New(log) // сюда передаем Storage структуру

	// init grpc
	grpcApp := grpcapp.New(paymentService, grpcPort)

	// init webhook component
	webhookApp := webhookapp.New(webHookPort)

	return &App{GRPCServer: grpcApp, Webhook: webhookApp}
}
