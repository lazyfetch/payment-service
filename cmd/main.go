package main

import (
	"log/slog"
	"os"
	"payment/internal/config"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

func main() {

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	// TODO: INIT LOGGER

	// TODO: SETUP APP (db, kafka, grpc, webhook)

	// TODO: START SERVER

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
