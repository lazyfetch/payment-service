package models

type Payment struct {
	Name          string
	Description   string
	Amount        int64
	PaymentMethod string
	UserID        string
}
