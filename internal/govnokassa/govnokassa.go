package govnokassa

import (
	"fmt"
	"payment/internal/domain/models"
)

type Govnokassa struct{}

func (g *Govnokassa) GeneratePaymentURL(data *models.DBPayment) (string, error) {

	url := fmt.Sprintf("https://govnokassa.local/pay/inv_id=%s?description=%s", data.IdempotencyKey, data.Description)

	return url, nil
}

func (g *Govnokassa) ProcessWebhook() error {

	return nil
}
