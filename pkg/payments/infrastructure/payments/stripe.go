package payments

import (
	"fmt"
	"github.com/gpioblink/go-stripe-book-seller/pkg/common/price"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/checkout/session"
)

type ConfirmedPayment struct {
	orderID string
	Price   price.Price
}

type StripeService struct {
	// TODO: product情報を取得して決済画面で商品の詳細を表示できるようにする
}

func NewStripeService(apiKey string) StripeService {
	stripe.Key = apiKey
	return StripeService{}
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
		SuccessURL: stripe.String("https://example.com/success?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String("https://example.com/cancel"),
	}
	params.SetIdempotencyKey(orderID)
	params.AddMetadata("order_id", orderID)

	_, err := session.New(params)
	if err != nil {
		return err
	}

	return nil
}
