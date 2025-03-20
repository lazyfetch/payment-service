package postgres

import (
	"context"
	"payment/internal/domain/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	Conn *pgxpool.Pool
}

func New() *Storage {

	conn, err := pgxpool.New(context.Background(), "") // протянуть конфиг бд, temp
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
