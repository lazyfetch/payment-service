package models

type PaymentRequest struct {
	Name        string
	Description string
	Amount      string
	UserID      string
	AiResponse  bool
}
