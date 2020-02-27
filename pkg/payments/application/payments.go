package application

import (
	"log"
	"time"

	"github.com/gpioblink/go-stripe-book-seller/pkg/common/price"
)

type ordersService interface {
	MarkOrderAsPaid(orderID string) error
}

type providerService interface {
	InitPaymentProvider(orderID string, price price.Price) error
}

type PaymentsService struct {
	ordersService   ordersService
	providerService providerService
}

func NewPaymentsService(ordersService ordersService, providerService providerService) PaymentsService {
	return PaymentsService{ordersService, providerService}
}

func (s PaymentsService) InitializeOrderPayment(orderID string, price price.Price) error {
	// ...
	log.Printf("initializing payment for order %s", orderID)

	go func() {
		time.Sleep(time.Millisecond * 500)
		if err := s.providerService.InitPaymentProvider(orderID, price); err != nil {
			log.Printf("cannot post order payment: %s", err)
		}
	}()

	// simulating payments provider delay
	//time.Sleep(time.Second)

	return nil
}

func (s PaymentsService) PostOrderPayment(orderID string) error {
	log.Printf("payment for order %s done, marking order as paid", orderID)

	return s.ordersService.MarkOrderAsPaid(orderID)
}
