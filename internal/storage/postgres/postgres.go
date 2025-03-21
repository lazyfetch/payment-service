package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"payment/internal/config"
	"payment/internal/domain/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	Log  *slog.Logger
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

	return &Storage{Log: log, Conn: conn}
}

func (s *Storage) Stop() {
	s.Conn.Close()
}

func (s *Storage) CreatePayment(ctx context.Context, data models.DBPayment) error {
	return nil // temp
}

func (s *Storage) UpdatePayment(ctx context.Context, data models.DBPayment) error {

	return nil // temp
}

func (s *Storage) User(ctx context.Context, userID string) (models.User, error) {

	return models.User{}, nil // temp
}
