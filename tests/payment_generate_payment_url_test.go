package tests

import (
	"payment/proto/gen/payment"
	"payment/tests/suite"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type DefaultForm struct {
	Name          string
	Description   string
	Amount        int64
	PaymentMethod string
	UserID        string
}

var (
	defaultForm = DefaultForm{
		Name:          "14.1",
		Description:   "Privet, papich! How its going?",
		Amount:        5000,
		PaymentMethod: "Govnokassa",
		UserID:        "stray228",
	}
)

func TestGeneratePaymentUrl_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	respUrl, err := st.AuthClient.GetPaymentUrl(ctx, &payment.GetPaymentUrlRequest{
		Name:          defaultForm.Name,
		Description:   defaultForm.Description,
		Amount:        defaultForm.Amount,
		PaymentMethod: defaultForm.PaymentMethod,
		UserId:        defaultForm.UserID,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respUrl.PaymentUrl)

}

// not working lol
func TestGeneratePaymentUrl_BadPath(t *testing.T) {
	url := &payment.GetPaymentUrlRequest{
		Name:          defaultForm.Name,
		Description:   defaultForm.Description,
		Amount:        defaultForm.Amount,
		PaymentMethod: defaultForm.PaymentMethod,
		UserId:        defaultForm.UserID,
	}

	ctx, st := suite.New(t)

	url.Description = 
	// respUrl, err := st.AuthClient.GetPaymentUrl(ctx, url)
}
