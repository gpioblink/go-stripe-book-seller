package application

import (
	"log"

	"github.com/gpioblink/go-stripe-book-seller/pkg/common/price"
	"github.com/gpioblink/go-stripe-book-seller/pkg/orders/domain/orders"
	"github.com/pkg/errors"
)

type printerService interface {
	PrintReceipt(order orders.Order) error
}

type productsService interface {
	ProductByID(id orders.ProductID) (orders.Product, error)
}

type paymentsService interface {
	InitializeOrderPayment(id orders.ID, name string, price price.Price) error
}

type OrdersService struct {
	productsService productsService
	paymentsService paymentsService
	printerService  printerService

	ordersRepository orders.Repository
}

func NewOrdersService(productsService productsService, paymentsService paymentsService, printerService printerService, ordersRepository orders.Repository) OrdersService {
	return OrdersService{productsService, paymentsService, printerService, ordersRepository}
}

type PlaceOrderCommandAddress struct {
	Name     string
	Street   string
	City     string
	PostCode string
	Country  string
}

type PlaceOrderCommand struct {
	OrderID   orders.ID
	ProductID orders.ProductID

	Address PlaceOrderCommandAddress
}

func (s OrdersService) PlaceOrder(cmd PlaceOrderCommand) error {
	address, err := orders.NewAddress(
		cmd.Address.Name,
		cmd.Address.Street,
		cmd.Address.City,
		cmd.Address.PostCode,
		cmd.Address.Country,
	)
	if err != nil {
		return errors.Wrap(err, "invalid address")
	}

	product, err := s.productsService.ProductByID(cmd.ProductID)
	if err != nil {
		return errors.Wrap(err, "cannot get product")
	}

	newOrder, err := orders.NewOrder(cmd.OrderID, product, address)
	if err != nil {
		return errors.Wrap(err, "cannot create order")
	}

	if err := s.ordersRepository.Save(newOrder); err != nil {
		return errors.Wrap(err, "cannot save order")
	}

	if err := s.paymentsService.InitializeOrderPayment(newOrder.ID(), newOrder.Product().Name(), newOrder.Product().Price()); err != nil {
		return errors.Wrap(err, "cannot initialize payment")
	}

	log.Printf("order %s placed", cmd.OrderID)

	return nil
}

type MarkOrderAsPaidCommand struct {
	OrderID orders.ID
}

func (s OrdersService) MarkOrderAsPaid(cmd MarkOrderAsPaidCommand) error {
	o, err := s.ordersRepository.ByID(cmd.OrderID)
	if err != nil {
		return errors.Wrapf(err, "cannot get order %s", cmd.OrderID)
	}

	err = s.printerService.PrintReceipt(*o)
	if err != nil {
		return errors.Wrapf(err, "cannot print receipt %s", cmd.OrderID)
	}

	o.MarkAsPaid()

	if err := s.ordersRepository.Save(o); err != nil {
		return errors.Wrap(err, "cannot save order")
	}

	log.Printf("marked order %s as paid", cmd.OrderID)

	return nil
}

func (s OrdersService) OrderByID(id orders.ID) (orders.Order, error) {
	o, err := s.ordersRepository.ByID(id)
	if err != nil {
		return orders.Order{}, errors.Wrapf(err, "cannot get order %s", id)
	}

	return *o, nil
}
