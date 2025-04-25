package Redis

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
	log     *slog.Logger
	client  *redis.Client
	lockTTL time.Duration
	reqTTL  time.Duration // look like a shit... temp
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

// TODO реализовать два метода этих, и жить спокойно <3 <3

const (
	keyBan      = "ban"
	keyRequests = "requests"
	keyBanLevel = "banlevel"

	keyMinAmount = "min_amount"

	keyEvents = "events"

	keyLock = "lock"
)

func (r *Redis) Allow(ctx context.Context, ip string) (bool, error) {
	// здесь ужасные аллокации, будут срать нам GC
	// обязательно выносим все эти билдеры в отдельную функцию,
	// и туда протягиваем конфиг, или шаманим как угодно temp
	banKey := buildKey(keyBan, ip)
	countKey := buildKey(keyRequests, ip)
	levelKey := buildKey(keyBanLevel, ip)

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

	key := buildKey(keyEvents, payload.Type)
	err = r.client.RPush(ctx, key, data).Err()
	if err != nil {
		return fmt.Errorf("failed to send event to redis: %w", err)
	}

	return nil
}

// Need to realize, so now u can hardcode TTL, but later... NO!111!
// А вообще мы не хотим давать право TTLить каким-то другим компонентам наш Redis
// лучше его добавлять из конфига в структуру, и не парится
// Лучше еще default value поставить, чтоб не было проблем
func (r *Redis) GetMinAmount(ctx context.Context, userID string) (int64, error) {
	const op = "Redis.GetMinAmount"

	log := r.log.With(
		slog.String("op", op),
		slog.String("user_id", userID),
	)

	log.Info("start check")
	key := buildKey(keyMinAmount, userID)

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

func (r *Redis) SetMinAmount(ctx context.Context, userID string, amount int) error {
	const op = "Redis.SetMinAmount"
	log := r.log.With(
		slog.String("op", op),
	)

	log.Info("start set")

	key := buildKey(keyMinAmount, userID)

	// EXTRA TEMP!!! SUPER TEMP!!! temp mean time.Duration()
	if err := r.client.Set(ctx, key, amount, 2).Err(); err != nil {
		log.Error("failed to set", sl.Err(err))
	}

	return nil
}

func (r *Redis) TryAcquireMinAmountLock(ctx context.Context, userID string) (bool, error) {
	const op = "Redis.TryAcquireMinAmountLock"
	log := r.log.With(
		slog.String("op", op),
	)

	log.Info("start try to lock")

	// key := buildKey(keyLock, userID)

	return false, nil
}

func (r *Redis) ReleaseMinAmountLock(ctx context.Context, userID string) error {
	const op = "Redis.ReleaseMinAmountLock"
	log := r.log.With(
		slog.String("op", op),
	)

	return nil
}
