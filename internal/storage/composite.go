package storage

import (
	"context"
	"log/slog"
)

type DBProvider interface {
	GetMinAmount(ctx context.Context, userID string) (int64, error)
}

type CacheProvider interface {
	GetMinAmount(ctx context.Context, userID string) (int64, error)
	SetMinAmount(ctx context.Context, userID string, amount int) error

	// For distributed singleflight pattern (singleflight pattern)
	// So we can say, it's like mutex, just for cache :)
	TryAcquireMinAmountLock(ctx context.Context, userID string) (bool, error)
	ReleaseMinAmountLock(ctx context.Context, userID string) error
}

type Composite struct {
	Log           *slog.Logger
	DBProvider    DBProvider
	CacheProvider CacheProvider
}

func (c *Composite) GetMinAmountWithCache(ctx context.Context, userID string) (int64, error) {
	const op = "Composite.GetMinAmountWithCache"

	log := c.Log.With(
		slog.String("op", op),
	)

	// Эта херня называется distributed singleflight
	minAmount, err := c.CacheProvider.GetMinAmount(ctx, userID)
	if err != nil {
		log.Error("temp")
	}
	// Сначала проверяем есть ли у нас в redis такая штука, то есть
	// ключ в редисе +- будет такой <min_amount:userID>:value

	// Если нету, значит блокируем через <lock:min_amount:userID>
	// SET lock:min_amount:<userID> "1" NX EX 5 где NX - если не существует, EX 5 - истекает через 5 минут
	// if !TryAcquireMinAmountLock {return err} ну тут надо будет обернуть все правильно, чтобы без херни

	// DBProvider.GetMinAmount()

	// CacheProvider.SetMinAmount()

	// DEL lock:min_amount:<userID>
	// т.е ReleaseMinAmountLock()

	// } return data тут типа рекурсия должна быть, пон.

	// return data

	return minAmount, nil // temp
}
