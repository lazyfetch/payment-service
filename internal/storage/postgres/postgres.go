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

type Postgres struct {
	log  *slog.Logger
	Conn *pgxpool.Pool
}

func BuildDSN(c config.PostgresConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.DBname)
}

func New(log *slog.Logger, config config.PostgresConfig) *Postgres {

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

	return &Postgres{log: log, Conn: conn}
}

func (s *Postgres) Stop() {
	s.Conn.Close()
	s.log.Info("postgres close connection")
}

func (s *Postgres) CreatePayment(ctx context.Context, data *models.DBPayment) error {
	const op = "Postgres.CreatePayment"

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

func (s *Postgres) UpdatePayment(ctx context.Context, idemKey string) error {
	const op = "Postgres.Update"

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

func (s *Postgres) IdempotencyAndStatus(ctx context.Context, idempotencyKey string) error {
	const op = "Postgres.IsIdempotencyKey"

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

func (s *Postgres) GetMinAmountByUser(ctx context.Context, userID string) (int64, error) {
	const op = "Postgres.User"
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

func (s *Postgres) CreateEvent(ctx context.Context, payload any) error {
	const op = "Postgres.CreateEvent"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("start create event")

	// hardcode, we dont like this, make utils or something, dont do this
	// temp !11!11
	pl, err := json.Marshal(payload)
	if err != nil {
		log.Error("failed to marshal payload", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	cmd, err := s.Conn.Exec(ctx, `INSERT INTO events (event_type, payload, status) VALUES
	($1, $2, $3)`, "payments.success", pl, "new") // здесь операцию payment.success мы хардкодим, вообще так делать не надо как будто
	// temp

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

func (s *Postgres) OutboxUpdatePaymentTx(ctx context.Context, idemKey string, payload any) error {
	const op = "Postgres.OutboxUpdatePaymentTx"

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

// Это временная шляпа, нужно использовать транзакцию что снизу
func (s *Postgres) GetNewEvent(ctx context.Context) (models.Event, error) {
	const op = "Postgres.GetNewEvent"

	// Здесь мы должны получить payload + id, где status new или in_progress + ttl < now - 3 minut (или какой там ТТЛ)

	var event models.Event

	err := s.Conn.QueryRow(ctx, `SELECT payload from events WHERE status = 'new'`).Scan()
	if err != nil {
		return models.Event{}, fmt.Errorf("%s: %w", op, err)
	}

	return event, nil
}

func (s *Postgres) UpdateEventTTL(ctx context.Context, id int) error {
	// Здесь мы обновляем статус с new -> in_progress + TTL время которое = now + 3, условно говоря
	return nil // temp
}

func (s *Postgres) UpdateEventStatus(ctx context.Context, id int) error {
	// Здесь в рамках транзакции мы ставим complete, ok, done или любой другой статус
	return nil // temp
}

func (s *Postgres) EventTx(ctx context.Context) error {

	/*
		BEGIN;

		-- 1. Выбрать и залочить одну задачу с нужным статусом или просроченной in_progress
		SELECT id, payload
		FROM events
		WHERE
		  (status = 'new' OR (status = 'in_progress' AND updated_at < NOW() - INTERVAL '3 minutes'))
		ORDER BY updated_at ASC
		LIMIT 1
		FOR UPDATE SKIP LOCKED;

		-- 2. Обновить статус на in_progress и обновить updated_at
		UPDATE events
		SET status = 'in_progress', updated_at = NOW()
		WHERE id = <ID_из_предыдущего_запроса>;

		COMMIT;
	*/

	return nil

}
