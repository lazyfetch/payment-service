package govnokassa

type GovnoPayment struct {
	IdempotencyKey string `json:"idempotency_key"`
	UserID         string `json:"user_id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Amount         int64  `json:"amount"`
}
