package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	application "payment/internal/app"
	"payment/internal/config"
	"payment/internal/lib/logger/sl"
	"payment/internal/telemetry"
	tc "payment/internal/telemetry/config"
	"syscall"
	"time"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

func setupTelemetryOpts(cfg *config.Config) []tc.Option {
	var opts []tc.Option
	// global
	opts = append(opts, tc.WithService(cfg.Telemetry.ServiceName))
	opts = append(opts, tc.WithInsecure(cfg.Telemetry.Insecure))

	// metrics
	opts = append(opts, tc.MetricsWithEndpoint(cfg.Telemetry.Metrics.Endpoint))
	opts = append(opts, tc.MetricsWithInterval(cfg.Telemetry.Metrics.Interval))
	opts = append(opts, tc.MetricsWithTimeout(cfg.Telemetry.Metrics.Timeout))
	//traces
	opts = append(opts, tc.TracesWithEndpoint(cfg.Telemetry.Traces.Endpoint))
	opts = append(opts, tc.TracesWithTimeout(cfg.Telemetry.Traces.Timeout))
	opts = append(opts, tc.TracesWithSampler(cfg.Telemetry.Traces.Sampler))
	opts = append(opts, tc.TracesWithRatio(cfg.Telemetry.Traces.SamplerRatio))

	return opts
}

func main() {

	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	app := application.New(log, cfg)

	// telemtry start
	telShutdown, err := telemetry.New(setupTelemetryOpts(cfg)...)
	if err != nil {
		panic(err) // nobrain, so dumai sam temp
	}
	// start

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

	// stop telemetry
	if err := telShutdown(ctx); err != nil {
		log.Error("telemetry shutdown error", sl.Err(err))
	}
	app.GRPCServer.Stop()

	if err := app.Webhook.Stop(ctx); err != nil {
		log.Error("Shutdown error", sl.Err(err))
	}
	app.Storage.Stop()

	app.Redis.Close()

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
