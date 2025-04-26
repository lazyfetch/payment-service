package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"payment/internal/config"
	"payment/internal/domain/models"
	"payment/internal/lib/logger/sl"
	"payment/internal/storage"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	log    *slog.Logger
	client *redis.Client
	// lockTTL time.Duration
	// reqTTL  time.Duration // look like a shit... temp
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

func buildKey(parts ...string) string {
	var b strings.Builder

	for i, part := range parts {
		if i > 0 {
			b.WriteByte(':')
		}
		b.WriteString(part)
	}

	return b.String()
}

// temp
func (r *Redis) Allow(ctx context.Context, ip string) (bool, error) {

	banKey := buildKey(KeyBan, ip)
	countKey := buildKey(KeyRequests, ip)
	levelKey := buildKey(KeyBanLevel, ip)

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

// temp потому что отстой, чисто через список делаем, не надежно выглядит, мб для теста ток...
func (r *Redis) SendEvent(ctx context.Context, payload models.Event) error {
	// выглядит как полный кал, temp
	data, err := json.Marshal(payload.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	key := buildKey(KeyEvents, payload.Type)
	err = r.client.RPush(ctx, key, data).Err()
	if err != nil {
		return fmt.Errorf("failed to send event to redis: %w", err)
	}

	return nil
}

func (r *Redis) GetMinAmount(ctx context.Context, userID string) (int64, error) {
	const op = "Redis.GetMinAmount"

	log := r.log.With(
		slog.String("op", op),
		slog.String("user_id", userID),
	)

	log.Info("start check")
	key := buildKey(KeyMinAmount, userID)

	value, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			log.Warn("user_id not exists")
			return 0, storage.ErrUserIDNotExists
		}
		log.Error("failed to check", sl.Err(err))
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	minAmount, err := strconv.ParseInt(value, 10, 64)

	if err != nil {
		log.Error("error to parse integer", sl.Err(err))
	}

	return minAmount, nil

}

// temp todo: провести конфиг для TTL нормально
func (r *Redis) SetMinAmount(ctx context.Context, userID string, amount int64) error {
	const op = "Redis.SetMinAmount"
	log := r.log.With(
		slog.String("op", op),
	)

	log.Info("start set")

	key := buildKey(KeyMinAmount, userID)
	// temp ttl value
	if err := r.client.Set(ctx, key, amount, 2).Err(); err != nil {
		log.Error("failed to set", sl.Err(err))
	}

	log.Info("success")

	return nil
}

/*
func (r *Redis) TryMinAmountLock(ctx context.Context, userID, lockID string) (bool, error) {
	const op = "Redis.TryMinAmountLock"
	log := r.log.With(
		slog.String("op", op),
	)

	log.Info("start try to lock")

	key := buildKey(KeyMinAmountLock, userID)

	// HARDCODE!!!! WARNING TEMP, 3 = SHIT DELETE ITS
	ok, err := r.client.SetNX(ctx, key, lockID, 3).Result()

	if err != nil {
		log.Error("failed to lock", sl.Err(err))
		return false, fmt.Errorf("%s:%w", op, err)
	}

	if !ok {
		log.Warn("key already created")
		return false, nil
	}

	log.Info("success lock")

	return true, nil
}

func (r *Redis) ReleaseMinAmountLock(ctx context.Context, userID, lockID string) (bool, error) {
	const op = "Redis.ReleaseMinAmountLock"
	log := r.log.With(
		slog.String("op", op),
	)

	log.Info("start try to unlock")

	key := buildKey(KeyMinAmountLock, userID)

	res, err := r.client.Eval(ctx, UnlockScript, []string{key}, lockID).Result()
	if err != nil {
		return false, fmt.Errorf("%s:%w", op, err)
	}

	deleted, ok := res.(int64)
	if !ok && deleted < 1 {
		log.Error("no performed on fields")
		return false, fmt.Errorf("no performed on fields")
	}

	log.Info("success unlock")

	return true, nil

}
*/
