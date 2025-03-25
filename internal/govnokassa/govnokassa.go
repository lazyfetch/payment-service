package govnokassa

import (
	"encoding/json"
	"fmt"
	"payment/internal/domain/models"
)

type Govnokassa struct{}

func (g *Govnokassa) GeneratePaymentURL(data *models.DBPayment) (string, error) { // conditional realization

	url := fmt.Sprintf("https://govnokassa.local/pay/inv_id=%s?user_id=%s", data.IdempotencyKey, data.UserID)

	return url, nil
}

func (g *Govnokassa) ValidateData(rawData []byte) (*GovnoPayment, error) {

	var p GovnoPayment

	err := json.Unmarshal(rawData, &p)
	if err != nil {
		return nil, err
	}

	if p.IdempotencyKey == "" || p.UserID == "" {
		return nil, fmt.Errorf("field's is empty")
	}

	return &p, nil
}
