package robokassa

import (
	"bytes"
	"io"
	"net/http"
	"payment/internal/domain/models"
)

type Robokassa struct {
	login    string
	password string
}

func New(login string, password string) *Robokassa {
	return &Robokassa{
		login:    login,
		password: password,
	}
}

func (r *Robokassa) GeneratePaymentURL(payment models.GRPCPayment) (string, error) {

	data := JWT{
		MerchantLogin:  r.login,
		InvoiceType:    "OneTime",
		OutSum:         float64(payment.Amount) / 100,
		ShpUsername:    payment.Name,
		ShpUserID:      payment.UserID,
		ShpDescription: payment.Description,
	}

	baseURL := "https://services.robokassa.ru/InvoiceServiceWebApi/api/CreateInvoice"

	token, err := GenerateJWT(r.password, data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer([]byte(token)))
	if err != nil {
		return "", err // temp
	}
	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err // temp ?
	}

	defer resp.Body.Close()
	url, _ := io.ReadAll(resp.Body)

	return string(url), nil
}

func (r *Robokassa) CheckResultURL() {

}
