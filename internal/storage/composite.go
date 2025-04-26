package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

type DBProvider interface {
	GetMinAmount(ctx context.Context, userID string) (int64, error)
}

type CacheProvider interface {
	GetMinAmount(ctx context.Context, userID string) (int64, error)
	SetMinAmount(ctx context.Context, userID string, amount int64) error
}

type Composite struct {
	Log           *slog.Logger
	DBProvider    DBProvider
	CacheProvider CacheProvider
	sfGroup       singleflight.Group // в app.go протянуть обязательно, хотя вроде и nil нормально будет пахать
}

// GetMinAmountWithCache получает минимальную сумму через кэш или базу
func (c *Composite) GetMinAmountWithCache(ctx context.Context, userID string) (int64, error) {
	const op = "Composite.GetMinAmountWithCache"

	minAmount, err := c.CacheProvider.GetMinAmount(ctx, userID)

	// good path
	if err == nil {
		if minAmount == 0 {
			return 0, ErrUserIDNotExists
		}
		return minAmount, nil
	}
	// unexpected
	if err != redis.Nil {
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	v, err, _ := c.sfGroup.Do("min_amount_"+userID, func() (interface{}, error) {
		amount, err := c.DBProvider.GetMinAmount(ctx, userID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {

				return int64(0), ErrUserIDNotExists
			}
			return int64(0), fmt.Errorf("%s:%w", op, err)
		}

		_ = c.CacheProvider.SetMinAmount(ctx, userID, amount)
		return amount, nil
	})
	if err != nil {
		return 0, err
	}
	return v.(int64), nil
}
