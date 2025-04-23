package Redis

import (
	"context"
	"fmt"
	"log/slog"
	"payment/internal/config"
	"payment/internal/lib/logger/sl"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	log    *slog.Logger
	Client *redis.Client
}

func New(log *slog.Logger, cfg config.RedisConfig) *Redis {
	socket := fmt.Sprintf("%s:%v", cfg.Host, cfg.Port)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     socket,
		Password: "", // temp for production i think we need pass
		DB:       0,  // dunno? ??? temp
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // temp, loks like hardcode
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Error("redis not responding", sl.Err(err))
		panic("failed to connection to Redis")
	}

	log.Info("redis connection is successful")

	return &Redis{
		log:    log,
		Client: redisClient,
	}
}

func (r *Redis) Close() {
	// const op = "redis.Close"

	if err := r.Client.Close(); err != nil {
		r.log.Error("redis closing connection err:", sl.Err(err))
	}

	r.log.Info("redis close connection")

}
