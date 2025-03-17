package app

import (
	grpcapp "payment/internal/app/grpc"
	paymentsrv "payment/internal/service/payment"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(port int) *App {

	paymentService := paymentsrv.New()
	grpcApp := grpcapp.New(paymentService, port)

	return &App{GRPCServer: grpcApp}
}
