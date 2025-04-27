package storage

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"payment/internal/lib/logger/sl"
	"time"

	"golang.org/x/sync/singleflight"
)

type DBProvider interface {
	GetMinAmount(ctx context.Context, userID string) (int64, error)
}

type CacheProvider interface {
	GetMinAmount(ctx context.Context, userID string) (int64, error)
	SetMinAmount(ctx context.Context, userID string, amount int64, userTTL time.Duration) error
}

type Composite struct {
	Log           *slog.Logger
	DBProvider    DBProvider
	CacheProvider CacheProvider
	sfGroup       singleflight.Group
	UserTTL       time.Duration
}

func New(log *slog.Logger, dbProvider DBProvider, chProvider CacheProvider, userTTL time.Duration) *Composite {

	return &Composite{
		Log:           log,
		DBProvider:    dbProvider,
		CacheProvider: chProvider,
		UserTTL:       userTTL,
	}

}

func (c *Composite) GetMinAmountWithCache(ctx context.Context, userID string) (int64, error) {
	const op = "Composite.GetMinAmountWithCache"

	log := c.Log.With(
		slog.String("op", op),
		slog.String("user_id", userID),
	)

	log.Info("start")

	minAmount, err := c.CacheProvider.GetMinAmount(ctx, userID)

	// good path
	if err == nil {
		if minAmount == 0 {
			log.Warn("user_id not exists")
			return 0, ErrUserIDNotExists
		}
		log.Info("success", slog.Int64("min_amount", minAmount))
		return minAmount, nil
	}
	// unexpected
	if err != ErrUserIDNotExists {
		log.Info("unexpected error", sl.Err(err))
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	log.Warn("not found user_id, start check db")

	// Здесь код дублируется, нужно исправить дубликацию (чтобы красиво было ;3)
	v, err, _ := c.sfGroup.Do("min_amount_"+userID, func() (interface{}, error) {
		amount, err := c.DBProvider.GetMinAmount(ctx, userID)
		if err != nil {
			if errors.Is(err, ErrUserIDNotFound) {
				log.Info("user_id not exists")
				_ = c.CacheProvider.SetMinAmount(ctx, userID, amount, c.UserTTL)
				return int64(0), ErrUserIDNotExists
			}
			log.Info("unexpected error", sl.Err(err))
			return int64(0), fmt.Errorf("%s:%w", op, err)
		}

		_ = c.CacheProvider.SetMinAmount(ctx, userID, amount, c.UserTTL)
		return amount, nil

	})

	if err != nil {
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	return v.(int64), nil
}
