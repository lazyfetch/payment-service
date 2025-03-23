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
func (s *Storage) IdempotencyAndStatus(ctx context.Context, idempotencyKey string) bool {
	op := "Storage.IsIdempotencyKey"

	log := s.log.With(
		slog.String("op", op),
		slog.String("idempotency_key", idempotencyKey),
	)

	log.Info("start update payment in postgres")

	var exists int

	err := s.Conn.QueryRow(ctx, `
    SELECT 1
    FROM payments
    WHERE idempotency_key = $1 AND status = 'in_progress'
    LIMIT 1
	`, idempotencyKey).Scan(&exists)

	if err != nil {
		return false // можно доработать temp, но не нужно пока
	}
	return true

}

func (s *Storage) CreatePayment(ctx context.Context, data *models.DBPayment) error {
	op := "Storage.CreatePayment"

	log := s.log.With(
		slog.String("op", op),
		slog.String("user_id", data.UserID),
	)

	log.Info("start create payment in postgres")

	cmd, err := s.Conn.Exec(ctx, `
	INSERT INTO payments (idempotency_key, name, description, amount, user_id, status, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, data.IdempotencyKey, data.Name, data.Description, data.Amount, data.UserID, data.Status, data.CreatedAt, data.UpdatedAt)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if num := cmd.RowsAffected(); num != 1 {
		log.Error("failed to affect on row") // temp
		return fmt.Errorf("unexpected number of rows affected: %d", num)
	}

	return nil
}

func (s *Storage) UpdatePayment(ctx context.Context, idemKey string) error {
	op := "Storage.Update"

	log := s.log.With(
		slog.String("op", op),
		slog.String("idempotency_key", idemKey),
	)

	log.Info("start update payment in postgres")

	cmd, err := s.Conn.Exec(ctx, `
    UPDATE payments
    SET status = 'success',
        updated_at = $2
    WHERE idempotency_key = $1
	`, idemKey, time.Now())

	if err != nil {
		return err
	}
	if cmd.RowsAffected() != 1 {
		return fmt.Errorf("update failed, rows affected: %d", cmd.RowsAffected())
	}
	return nil
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

func (s *Storage) UpdateAndOutboxPattern() {

}
