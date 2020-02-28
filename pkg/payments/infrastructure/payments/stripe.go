package payments

import (
	"fmt"
	"github.com/gpioblink/go-stripe-book-seller/pkg/common/price"
	payments "github.com/gpioblink/go-stripe-book-seller/pkg/payments/domain"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/checkout/session"
)

type ConfirmedPayment struct {
	orderID string
	Price   price.Price
}

type StripeService struct {
	repository payments.Repository
	// TODO: product情報を取得して決済画面で商品の詳細を表示できるようにする
}

func NewStripeService(apiKey string, repository payments.Repository) StripeService {
	stripe.Key = apiKey
	return StripeService{repository}
}

func (s StripeService) InitPaymentProvider(orderID string, price price.Price) error {
	//TODO: need editing

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Name:        stripe.String(fmt.Sprintf("Item %s", orderID)),
				Description: stripe.String("Sample product"),
				Amount:      stripe.Int64(int64(price.Cents())),
				Currency:    stripe.String(price.Currency()),
				Quantity:    stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String("https://stripe-sample.gpioblink.now.sh/success.html?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String("https://stripe-sample.gpioblink.now.sh/error.html"),
	}
	res, err := session.New(params)
	if err != nil {
		return err
	}

	// save paymentID and orderID as a table to convert order
	payment, err := payments.NewPayment(orderID, res.ID)
	if err != nil {
		return err
	}
	err = s.repository.Save(&payment)
	if err != nil {
		return err
	}

	return nil
}
