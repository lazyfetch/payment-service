package Redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"payment/internal/config"
	"payment/internal/domain/models"
	"payment/internal/lib/logger/sl"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	log    *slog.Logger
	client *redis.Client
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
		client: redisClient,
	}
}

func (r *Redis) Close() {
	// const op = "redis.Close"

	if err := r.client.Close(); err != nil {
		r.log.Error("redis closing connection err:", sl.Err(err))
	}

	r.log.Info("redis close connection")

}

// TODO реализовать два метода этих, и жить спокойно <3 <3

func (r *Redis) Allow(ctx context.Context, ip string) (bool, error) {
	// здесь ужасные аллокации, будут срать нам GC
	// обязательно выносим все эти билдеры в отдельную функцию,
	// и туда протягиваем конфиг, или шаманим как угодно temp
	banKey := fmt.Sprintf("ban:%s", ip)
	countKey := fmt.Sprintf("requests:%s", ip)
	levelKey := fmt.Sprintf("banlevel:%s", ip)

	// 1. Проверяем бан
	isBanned, err := r.client.Exists(ctx, banKey).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check ban: %w", err)
	}
	if isBanned > 0 {
		return false, nil
	}

	// 2. Инкремент запроса
	count, err := r.client.Incr(ctx, countKey).Result()
	if err != nil {
		return false, fmt.Errorf("failed to increment counter: %w", err)
	}

	if count == 1 {
		// Устанавливаем TTL на первый запрос
		r.client.Expire(ctx, countKey, 60*time.Second)
	}

	if count > 9 {
		// Получаем текущий уровень бана
		banLevel, _ := r.client.Get(ctx, levelKey).Int()
		if banLevel == 0 {
			banLevel = 1
		} else {
			banLevel *= 3 // экспоненциальный рост
		}

		banDuration := time.Duration(banLevel) * time.Minute
		r.client.Set(ctx, banKey, "1", banDuration)
		r.client.Set(ctx, levelKey, banLevel, 0)

		return false, nil
	}

	return true, nil
}

func (r *Redis) SendEvent(ctx context.Context, payload models.Event) error {
	// выглядит как полный кал, temp
	data, err := json.Marshal(payload.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// здесь плохая аллокация, будет срать нам GC
	key := fmt.Sprintf("events:%s", payload.Type)
	err = r.client.RPush(ctx, key, data).Err()
	if err != nil {
		return fmt.Errorf("failed to send event to redis: %w", err)
	}

	return nil
}
