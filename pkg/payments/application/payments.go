package application

import (
	payments "github.com/gpioblink/go-stripe-book-seller/pkg/payments/domain"
	"log"
	"time"

	"github.com/gpioblink/go-stripe-book-seller/pkg/common/price"
)

type ordersService interface {
	MarkOrderAsPaid(orderID string) error
}

type providerService interface {
	InitPaymentProvider(orderID string, name string, price price.Price) error
}

type PaymentsService struct {
	ordersService   ordersService
	providerService providerService
	repository      payments.Repository
}

func NewPaymentsService(ordersService ordersService, providerService providerService, repository payments.Repository) PaymentsService {
	return PaymentsService{ordersService, providerService, repository}
}

func (s PaymentsService) InitializeOrderPayment(orderID string, name string, price price.Price) error {
	// ...
	log.Printf("initializing payment for order %s", orderID)

	go func() {
		time.Sleep(time.Millisecond * 500)
		if err := s.providerService.InitPaymentProvider(orderID, name, price); err != nil {
			log.Printf("cannot post order payment: %s", err)
		}
	}()

	// simulating payments provider delay
	//time.Sleep(time.Second)

	return nil
}

func (s PaymentsService) PostOrderPayment(paymentID string) error {
	// convert PaymentID to OrderID, then send paid info to the order service
	payment, err := s.repository.ByPaymentID(paymentID)
	if err != nil {
		log.Printf("cannot find paymentID: %s", err)
		return err
	}
	log.Printf("payment for payment %s (%s) done, marking order as paid", paymentID, payment.OrderId())

	_ = s.repository.DeleteByPaymentID(paymentID)

	return s.ordersService.MarkOrderAsPaid(payment.OrderId())
}

// TODO: こんなユースケースあっていいのか？httpのインフラでそのまま使えばいいのでは？
func (s PaymentsService) GetPaymentID(orderId string) (string, error) {
	payment, err := s.repository.ByOrderID(orderId)
	if err != nil {
		return "", err
	}

	log.Printf(payment.PaymentId(), payment.OrderId())
	return payment.PaymentId(), nil
}
