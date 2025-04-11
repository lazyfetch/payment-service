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

func TestGeneratePaymentUrl_BadPath(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name string
		form DefaultForm
	}{
		{
			name: "empty name",
			form: DefaultForm{
				Name:          "",
				Description:   defaultForm.Description,
				Amount:        defaultForm.Amount,
				PaymentMethod: defaultForm.PaymentMethod,
				UserID:        defaultForm.UserID,
			},
		},
		{
			name: "empty description",
			form: DefaultForm{
				Name:          defaultForm.Name,
				Description:   "",
				Amount:        defaultForm.Amount,
				PaymentMethod: defaultForm.PaymentMethod,
				UserID:        defaultForm.UserID,
			},
		},
		{
			name: "zero amount",
			form: DefaultForm{
				Name:          defaultForm.Name,
				Description:   defaultForm.Description,
				Amount:        0,
				PaymentMethod: defaultForm.PaymentMethod,
				UserID:        defaultForm.UserID,
			},
		},
		{
			name: "empty payment method",
			form: DefaultForm{
				Name:          defaultForm.Name,
				Description:   defaultForm.Description,
				Amount:        defaultForm.Amount,
				PaymentMethod: "",
				UserID:        defaultForm.UserID,
			},
		},
		{
			name: "empty user id",
			form: DefaultForm{
				Name:          defaultForm.Name,
				Description:   defaultForm.Description,
				Amount:        defaultForm.Amount,
				PaymentMethod: defaultForm.PaymentMethod,
				UserID:        "",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := st.AuthClient.GetPaymentUrl(ctx, &payment.GetPaymentUrlRequest{
				Name:          tc.form.Name,
				Description:   tc.form.Description,
				Amount:        tc.form.Amount,
				PaymentMethod: tc.form.PaymentMethod,
				UserId:        tc.form.UserID,
			})
			require.Error(t, err, "expected error but got none")
		})
	}
}
