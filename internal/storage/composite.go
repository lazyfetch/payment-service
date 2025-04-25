package storage

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"payment/internal/lib/uuid"
	"time"
)

type DBProvider interface {
	GetMinAmount(ctx context.Context, userID string) (int64, error)

	UserExists(ctx context.Context, userID string) (bool, error)
}

type CacheProvider interface {
	UserExists(ctx context.Context, userID string) (bool, error)

	GetMinAmount(ctx context.Context, userID string) (int64, error)
	SetMinAmount(ctx context.Context, userID string, amount int) error

	// For distributed singleflight pattern (singleflight pattern)
	// So we can say, it's like mutex, just for cache :)
	TryMinAmountLock(ctx context.Context, userID, lockID string) (bool, error)
	ReleaseMinAmountLock(ctx context.Context, userID, lockID string) (bool, error)
}

type Composite struct {
	Log           *slog.Logger
	DBProvider    DBProvider
	CacheProvider CacheProvider
}

func (c *Composite) UserExists(ctx context.Context, userID string) (bool, error) {
	return false, nil // temp
}

// ЗДЕСЬ ПЛАН ТАКОЙ, СНАЧАЛА ПРОВЕРЯЕМ IsExistsWithCache
// если ЮЗЕРА ВООБЩЕ НЕТУ, ТО МЫ И НЕ ПРОДОЛЖАЕМ ДАЛЬШЕ
func (c *Composite) GetMinAmountWithCache(ctx context.Context, userID string) (int64, error) {
	const op = "Composite.GetMinAmountWithCache"

	log := c.Log.With(
		slog.String("op", op),
		slog.String("user_id", userID),
	)

	log.Info("start operation")

	log.Info("check if user_id exists")

	ok, err := c.UserExists(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	// Здесь короче каша выходит, давай завтра доделаем, седня и так +400 строк где-то, окей нормально
	if !ok {
		return 0, ErrUserIDNotExists
	}

	// Эта херня называется distributed singleflight
	minAmount, err := c.CacheProvider.GetMinAmount(ctx, userID)
	if err != nil {
		if errors.Is(err, ErrUserIDNotExists) {

			log.Warn("min_amount not found")

			lockID := uuid.UUID()
			ok, err := c.CacheProvider.TryMinAmountLock(ctx, userID, lockID)

			if err != nil {
				return 0, fmt.Errorf("%s:%w", op, err)
			}

			if !ok {
				time.Sleep(time.Millisecond * 250) // temp, look like a shit
				return c.GetMinAmountWithCache(ctx, userID)
			}

		}
		return 0, fmt.Errorf("%s:%w", op, err)
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
