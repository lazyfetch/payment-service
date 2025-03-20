package main

import (
	"log/slog"
	"os"
	"os/signal"
	application "payment/internal/app"
	"payment/internal/config"
	"syscall"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

func main() {

	// TODO: INIT CONFIG
	cfg := config.MustLoad()

	// TODO: INIT LOGGER
	log := setupLogger(cfg.Env)

	// TODO: SETUP APP (db, kafka, grpc, webhook)
	app := application.New(log, cfg.Webhook.Port, cfg.GRPC.Port, cfg.RoboKassa.MerchantLogin, cfg.RoboKassa.Password)

	// TODO: START SERVER
	app.GRPCServer.MustRun()
	app.Webhook.MustRun()

	// TODO: Graceful shutdown for server & kafka, db, other shit

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	app.GRPCServer.Stop()
	app.Webhook.Stop()

	log.Info("Gracefully stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envDev:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
