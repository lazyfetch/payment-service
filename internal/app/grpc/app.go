package grpcapp

import (
	"fmt"
	"log/slog"
	"net"
	"payment/internal/app/grpc/interceptors"
	"payment/internal/config"
	paymentgrpc "payment/internal/grpc/payment"

	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, paymentService paymentgrpc.PaymentService, RateLimiter interceptors.RateLimiter, cfg *config.Config) *App {
	gRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			interceptors.ValidationInterceptor(&cfg.Internal),
		),
		grpc.UnaryInterceptor(
			interceptors.LimiterInterceptor(&cfg.Internal, RateLimiter),
		),
	)

	paymentgrpc.Register(gRPCServer, paymentService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       cfg.GRPC.Port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.log.Info("grpc server started", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {

}
