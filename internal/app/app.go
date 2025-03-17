package app

import (
	grpcapp "payment/internal/app/grpc"
	paymentsrv "payment/internal/service/payment"
	webhookapp "payment/internal/webhook"
)

type App struct {
	GRPCServer *grpcapp.App
	Webhook    *webhookapp.App
}

func New(webHookPort int, grpcPort int) *App {

	// init service layer
	paymentService := paymentsrv.New()

	// init grpc
	grpcApp := grpcapp.New(paymentService, grpcPort)

	// init webhook component
	webhookApp := webhookapp.New(webHookPort)

	return &App{GRPCServer: grpcApp, Webhook: webhookApp}
}
