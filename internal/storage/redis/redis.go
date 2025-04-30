package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"payment/internal/domain/models"
	"payment/internal/lib/logger/sl"
	"payment/internal/storage"
	"payment/internal/telemetry/tracing"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/attribute"
)

type RedisOpts struct {
	Host string
	Port int
	DB   int
}

type Redis struct {
	log    *slog.Logger
	client *redis.Client
	// lockTTL time.Duration
	// reqTTL  time.Duration // look like a shit... temp
}

func New(log *slog.Logger, redisOpts RedisOpts) *Redis {
	socket := fmt.Sprintf("%s:%v", redisOpts.Host, redisOpts.Port)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     socket,
		Password: "", // temp for production i think we need pass
		DB:       redisOpts.DB,
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

// false - бан есть, true - все чисто
func (r *Redis) Allow(ctx context.Context, ip string, window time.Duration, maxRequests int, banDurations []time.Duration) (bool, error) {

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
		// Устанавливаем window TTL на первый запрос
		r.client.Expire(ctx, countKey, window)
	}

	// если запросов больше чем максимально допустимо
	// Получаем уровень бана, каждый индекс бана = большему времени бана
	// уровень бана не должен привышать максимальное значение len(banDurations - 1)

	if count > int64(maxRequests) {
		// Получаем текущий уровень бана
		banLevel, err := r.client.Get(ctx, levelKey).Int()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				r.client.Set(ctx, levelKey, banLevel, 0)
				// temp, hardcode
				r.client.Expire(ctx, levelKey, time.Duration(time.Minute*30))
			}
		}

		var banDuration time.Duration

		if banLevel >= len(banDurations)-1 {
			banDuration = banDurations[len(banDurations)-1]
		} else {
			banDuration = banDurations[banLevel]
			r.client.Incr(ctx, levelKey).Result()
		}

		r.client.Set(ctx, banKey, "1", banDuration)

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

	ctx, span := tracing.StartSpan(ctx, "Redis GetMinAmount", attribute.String("user_id", userID))
	defer span.End()

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
		span.RecordError(err)
		log.Error("failed to check", sl.Err(err))
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	minAmount, err := strconv.ParseInt(value, 10, 64)

	if err != nil {
		span.RecordError(err)
		log.Error("error to parse integer", sl.Err(err))
	}

	return minAmount, nil

}

// temp todo: провести конфиг для TTL нормально
func (r *Redis) SetMinAmount(ctx context.Context, userID string, amount int64, userTTL time.Duration) error {
	const op = "Redis.SetMinAmount"

	ctx, span := tracing.StartSpan(ctx, "Redis SetMinAmount",
		attribute.String("user_id", userID),
		attribute.Int64("amount", amount))
	span.End()

	log := r.log.With(
		slog.String("op", op),
	)

	log.Info("start set")

	key := buildKey(KeyMinAmount, userID)
	// temp ttl value
	if err := r.client.Set(ctx, key, amount, userTTL).Err(); err != nil {
		span.RecordError(err)
		log.Error("failed to set", sl.Err(err))
	}

	log.Info("success")

	return nil
}
