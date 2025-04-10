package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	application "payment/internal/app"
	"payment/internal/config"
	"payment/internal/lib/logger/sl"
	"syscall"
	"time"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

func main() {

	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	app := application.New(log, cfg)

	go app.GRPCServer.MustRun()
	go func() {
		if err := app.Webhook.Run(); err != nil {
			log.Error("Webhook server error: ", sl.Err(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	log.Info("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	app.GRPCServer.Stop()
	if err := app.Webhook.Stop(ctx); err != nil {
		log.Error("Shutdown error")
	}
	app.Storage.Stop()

	log.Info("Server gracefully stopped")
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
