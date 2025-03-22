package models

import (
	"time"
)

type GRPCPayment struct {
	Name          string
	Description   string
	Amount        int64
	PaymentMethod string
	UserID        string
}

type DBPayment struct {
	Name           string
	Description    string
	Amount         int64
	UserID         string
	IdempotencyKey string
	Status         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func MapGRPCToDB(input *GRPCPayment, idemKey string) *DBPayment {
	return &DBPayment{
		Name:           input.Name,
		Description:    input.Description,
		Amount:         input.Amount,
		UserID:         input.UserID,
		IdempotencyKey: idemKey,
		Status:         "in_progress",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}
