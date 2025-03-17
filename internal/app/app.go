package app

import (
	grpcapp "payment/internal/app/grpc"
	"payment/internal/rest/webhook"
	paymentsrv "payment/internal/service/payment"
)

type App struct {
	GRPCServer *grpcapp.App
	Webhook    *webhook.App
}

func New(webHook int, grpcPort int) *App {

	// init service layer
	paymentService := paymentsrv.New()

	// init grpc
	grpcApp := grpcapp.New(paymentService, grpcPort)

	// init webhook component
	webhookApp := webhook.New(webHook)

	return &App{GRPCServer: grpcApp, Webhook: webhookApp}
}
