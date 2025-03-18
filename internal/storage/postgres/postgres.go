package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
}

func MustRun() *Storage {

	conn, err := pgxpool.New()
	if err != nil {
		panic(err) // rework soon
	}

	return &Storage{}
}

func (s *Storage) Stop() {

}

func (s *Storage) HealthCheck() {

}

func (s *Storage) CreatePayment(ctx context.Context) {

}
