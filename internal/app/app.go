package app

import (
	"context"
	"log/slog"
	grpcapp "payment/internal/app/grpc"
	webhookapp "payment/internal/app/webhook"
	"payment/internal/config"
	"payment/internal/govnokassa"
	confirmsrv "payment/internal/service/confirm"
	eventsender "payment/internal/service/event_sender"
	generatesrv "payment/internal/service/generate"
	"payment/internal/storage/postgres"
	Redis "payment/internal/storage/redis"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
	Webhook    *webhookapp.App
	Redis      *Redis.Redis
	Storage    *postgres.Storage
}

func New(log *slog.Logger, config *config.Config) *App {

	// VERY temp

	t := time.Second * 5
	sender := eventsender.Sender{Log: log}
	sender.StartProcessEvents(context.Background(), t)

	// payment service
	gvkassa := &govnokassa.Govnokassa{}

	// init db
	storage := postgres.New(log, config.Postgres)

	// init redis
	redis := Redis.New(log, config.Redis)

	// init gen service
	generateService := generatesrv.New(log, storage, storage, gvkassa)

	// init webhook service
	confirmService := confirmsrv.New(log, storage, storage, gvkassa)

	// init grpc
	grpcApp := grpcapp.New(log, generateService, config)

	// init webhook
	webhookApp := webhookapp.New(log, confirmService, config.Webhook.Port)

	return &App{GRPCServer: grpcApp, Webhook: webhookApp, Redis: redis, Storage: storage}
}
