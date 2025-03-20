package models

import "time"

type GRPCPayment struct {
	Name          string
	Description   string
	Amount        int64
	PaymentMethod string
	UserID        string
}

type DBPayment struct {
	GRPCPayment
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
