package app

import (
	"log/slog"
	grpcapp "payment/internal/app/grpc"
	webhookapp "payment/internal/app/webhook"
	"payment/internal/config"
	"payment/internal/govnokassa"
	confirmsrv "payment/internal/service/confirm"
	generatesrv "payment/internal/service/generate"
	"payment/internal/storage"
	"payment/internal/storage/postgres"
	Redis "payment/internal/storage/redis"
)

type App struct {
	GRPCServer *grpcapp.App
	Webhook    *webhookapp.App
	Redis      *Redis.Redis
	Storage    *postgres.Postgres
}

func New(log *slog.Logger, config *config.Config) *App {

	// payment service mock
	gvkassa := &govnokassa.Govnokassa{}

	// init db with opts
	db := postgres.New(log, postgres.PostgresOpts{
		User:     config.Postgres.User,
		Password: config.Postgres.Password,
		Host:     config.Postgres.Host,
		Port:     config.Postgres.Port,
		DBname:   config.Postgres.DBname,
	})

	// init redis with opts
	cache := Redis.New(log, Redis.RedisOpts{
		Host: config.Redis.Host,
		Port: config.Redis.Port,
		DB:   config.Redis.DB,
	})

	// init compositor
	composite := storage.New(log, db, cache, config.Internal.UserTTL)

	// init generate service
	generateService := generatesrv.New(log, db, composite, gvkassa)

	// init configrm service
	confirmService := confirmsrv.New(log, db, db, gvkassa)

	// init grpc
	grpcApp := grpcapp.New(log, generateService, cache, config)

	// init webhook
	webhookApp := webhookapp.New(log, confirmService, config.Webhook.Port)

	return &App{GRPCServer: grpcApp, Webhook: webhookApp, Redis: cache, Storage: db}
}
