package main

import (
	"log/slog"
	"os"
	application "payment/internal/app"
	"payment/internal/config"
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
	app := application.New(log, cfg.Webhook.Port, cfg.GRPC.Port)

	// TODO: START SERVER
	app.GRPCServer.MustRun()
	app.Webhook.MustRun()

	// TODO: Graceful shutdown for server & kafka, db, other shit
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
