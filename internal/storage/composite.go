package storage

import (
	"context"
	"time"
)

type DBProvider interface {
	GetMinAmountByUser(ctx context.Context, userID string) (int64, error)
}

type CacheProvider interface {
	GetMinAmountByUser(ctx context.Context, userID string) (int64, error)
	SetMinAmountByUser(ctx context.Context, userID string, amount int, TTL time.Duration) error

	// For distributed singleflight pattern (singleflight pattern)
	// So we can say, it's like mutex, just for cache :)
	TryAcquireMinAmountLock(ctx context.Context, userID string, TTL time.Duration) (bool, error)
	ReleaseMinAmountLock(ctx context.Context, userID string) error
}

type Composite struct {
	DBProvider    DBProvider
	CacheProvider CacheProvider
}

func (c *Composite) GetMinAmountByUserWithCache(ctx context.Context, userID string) (int64, error) {

	// Эта херня называется distributed singleflight

	// Сначала проверяем есть ли у нас в redis такая штука, то есть
	// ключ в редисе +- будет такой <min_amount:userID>:value

	// Если нету, значит блокируем через <lock:min_amount:userID>
	// SET lock:min_amount:<userID> "1" NX EX 5 где NX - если не существует, EX 5 - истекает через 5 минут
	// if !TryAcquireMinAmountLock {return err} ну тут надо будет обернуть все правильно, чтобы без херни

	// DBProvider.GetMinAmountByUser()

	// CacheProvider.SetMinAmountByUser()

	// } return data

	// return data

	return 0, nil // temp
}
