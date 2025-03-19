package robokassa

import (
	"bytes"
	"io"
	"net/http"
	"payment/internal/domain/models"
)

func GeneratePaymentURL(payment models.Payment) (string, error) {

	data := Payload{
		MerchantLogin:  "Some",
		InvoiceType:    "",
		OutSum:         float64(payment.Amount) / 100,
		ShpUsername:    payment.Name,
		ShpUserID:      payment.UserID,
		ShpDescription: payment.Description,
	}

	baseURL := "https://services.robokassa.ru/InvoiceServiceWebApi/api/CreateInvoice"

	token, err := GenerateJWT("imp_me", data)
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

	url, _ := io.ReadAll(resp.Body)

	return string(url), nil
}

func CheckResultURL() {

}
