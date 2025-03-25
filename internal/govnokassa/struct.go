package govnokassa

type GovnoPayment struct {
	IdempotencyKey string `json:"idempotency_key"`
	UserID         string `json:"user_id"`
}
