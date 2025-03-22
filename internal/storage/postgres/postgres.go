package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"payment/internal/config"
	"payment/internal/domain/models"
	"payment/internal/lib/logger/sl"
	"payment/internal/storage"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	log  *slog.Logger
	Conn *pgxpool.Pool
}

func BuildDSN(c config.PostgresConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.DBname)
}

func New(log *slog.Logger, config config.PostgresConfig) *Storage {

	conn, err := pgxpool.New(context.Background(), BuildDSN(config)) // протянуть конфиг бд, temp
	if err != nil {
		panic(err) // temp
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := conn.Ping(ctx); err != nil {
		panic(err)
	}

	log.Info("postgres connection is successful")

	return &Storage{log: log, Conn: conn}
}

func (s *Storage) Stop() {
	s.Conn.Close()
}

// IsIdempotencyKey returns true of false.
// True is their find same idempotency key, or false if not.
func (s *Storage) IsIdempotencyKey(ctx context.Context, data *models.DBPayment) (bool, error) {
	op := "Storage.IsIdempotencyKey"

	log := s.log.With(
		slog.String("op", op),
		slog.String("user_id", data.UserID),
	)

	log.Info("start update payment in postgres")
	return false, nil // temp
}

func (s *Storage) CreatePayment(ctx context.Context, data *models.DBPayment) error {
	op := "Storage.CreatePayment"

	log := s.log.With(
		slog.String("op", op),
		slog.String("user_id", data.UserID),
	)

	log.Info("start create payment in postgres")

	return nil // temp
}

func (s *Storage) UpdatePayment(ctx context.Context, data *models.DBPayment) error {
	op := "Storage.CreatePayment"

	log := s.log.With(
		slog.String("op", op),
		slog.String("user_id", data.UserID),
	)

	log.Info("start update payment in postgres")
	return nil // temp
}

func (s *Storage) GetMinAmountByUser(ctx context.Context, userID string) (int64, error) {
	var minAmount int64

	op := "Storage.User"

	log := s.log.With(
		slog.String("op", op),
		slog.String("user_id", userID),
	)

	log.Info("start check min_amount of user")

	err := s.Conn.QueryRow(ctx, "SELECT min_amount FROM users WHERE user_id=$1", userID).Scan(&minAmount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn("user_id not found")
			return 0, storage.ErrUserIDNotFound
		}
		log.Warn("failed to check user", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("min_amount exists")
	return minAmount, nil // temp
}
