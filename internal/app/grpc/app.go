package grpcapp

import (
	paymentgrpc "payment/internal/grpc/payment"

	"google.golang.org/grpc"
)

type App struct {
	gRPCServer *grpc.Server
	port       int
}

func New(paymentService paymentgrpc.PaymentService, port int) *App {
	gRPCServer := grpc.NewServer()

	paymentgrpc.Register(gRPCServer, paymentService)

	return &App{
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {

	return nil // temp
}
