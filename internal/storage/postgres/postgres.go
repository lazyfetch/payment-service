package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"payment/internal/config"
	"payment/internal/domain/models"
	"payment/internal/lib/logger/sl"
	"payment/internal/storage"
	"time"

	"github.com/jackc/pgx/v5"
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

	conn, err := pgxpool.New(context.Background(), BuildDSN(config))
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // temp, loks like hardcode
	defer cancel()

	if err := conn.Ping(ctx); err != nil {
		panic(err)
	}

	log.Info("postgres connection is successful")

	return &Storage{log: log, Conn: conn}
}

func (s *Storage) Stop() {
	s.Conn.Close()
	s.log.Info("postgres close connection")
}

func (s *Storage) CreatePayment(ctx context.Context, data *models.DBPayment) error {
	const op = "Storage.CreatePayment"

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
		log.Error("unexpected error", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	if num := cmd.RowsAffected(); num != 1 {
		log.Error("failed to affect on row")
		return fmt.Errorf("unexpected number of rows affected: %d", num)
	}

	return nil
}

func (s *Storage) UpdatePayment(ctx context.Context, idemKey string) error {
	const op = "Storage.Update"

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
		log.Error("failed to update payment", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	if num := cmd.RowsAffected(); num != 1 {
		log.Error("failed to affect on row")
		return fmt.Errorf("unexpected number of rows affected: %d", num)
	}
	return nil
}

func (s *Storage) IdempotencyAndStatus(ctx context.Context, idempotencyKey string) error {
	const op = "Storage.IsIdempotencyKey"

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
		if errors.Is(err, sql.ErrNoRows) {
			log.Info("in progress idempotency_key not found")
			return sql.ErrNoRows
		}
		log.Error("unexpected error", sl.Err(err))
		return err
	}
	return nil

}

func (s *Storage) GetMinAmountByUser(ctx context.Context, userID string) (int64, error) {
	const op = "Storage.User"
	var minAmount int64

	log := s.log.With(
		slog.String("op", op),
		slog.String("user_id", userID),
	)

	log.Info("start check min_amount of user")

	err := s.Conn.QueryRow(ctx, "SELECT min_amount FROM users WHERE user_id=$1", userID).Scan(&minAmount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Info("user_id not found")
			return 0, storage.ErrUserIDNotFound
		}
		log.Error("failed to check user", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("min_amount exists")
	return minAmount, nil
}

func (s *Storage) CreateEvent(ctx context.Context, payload any) error {
	const op = "Storage.CreateEvent"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("start create event")

	pl, err := json.Marshal(payload)
	if err != nil {
		log.Error("failed to marshal payload", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	cmd, err := s.Conn.Exec(ctx, `INSERT INTO events (event_type, payload, status) VALUES
	($1, $2, $3)`, "payments.success", pl, "new")

	if err != nil {
		log.Error("unexpected error", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	if num := cmd.RowsAffected(); num != 1 {
		log.Error("failed to affect on row")
		return fmt.Errorf("unexpected number of rows affected: %d", num)
	}

	return nil

}

func (s *Storage) OutboxUpdatePayment(ctx context.Context, idemKey string, payload any) error {
	const op = "Storage.OutboxUpdatePayment"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("start outbox transaction")

	tx, err := s.Conn.Begin(ctx)
	if err != nil {
		log.Error("failed to start tx", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	// defer tx.Rollback
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Error("rollback failed", sl.Err(err))
		}
	}()

	// main operation in transaction
	s.CreateEvent(ctx, payload)
	s.UpdatePayment(ctx, idemKey)

	// commit
	err = tx.Commit(ctx)
	if err != nil {
		log.Error("failed to commit tx", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

/* func (s *Storage) GetNewEvent(ctx context.Context) (models.Event, error) {
	const op = "Storage.GetNewEvent"

	var event models.Event

	err := s.Conn.QueryRow(ctx, "").Scan()
	if err != nil {
		return models.Event{}, fmt.Errorf("%s: %w", op, err)
	}

} */
