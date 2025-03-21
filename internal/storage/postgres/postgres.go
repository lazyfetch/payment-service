package postgres

import (
	"context"
	"fmt"
	"payment/internal/config"
	"payment/internal/domain/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	Conn *pgxpool.Pool
}

func BuildDSN(c config.PostgresConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		c.User, c.Password, c.Host, c.Port, c.DBname)
}

func New(config config.PostgresConfig) *Storage {

	conn, err := pgxpool.New(context.Background(), BuildDSN(config)) // протянуть конфиг бд, temp
	if err != nil {
		panic(err) // temp
	}

	return &Storage{Conn: conn}
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
