package robokassa

import (
	"bytes"
	"io"
	"net/http"
)

func GeneratePaymentURL() (string, error) {

	baseURL := "https://services.robokassa.ru/InvoiceServiceWebApi/api/CreateInvoice"

	token := GenerateJWT()

	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer([]byte(token)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	url, _ := io.ReadAll(resp.Body)

	return string(url), nil
}

func CheckResultURL() {

}
