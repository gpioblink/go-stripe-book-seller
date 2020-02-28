package intraprocess

import "github.com/gpioblink/go-stripe-book-seller/pkg/printer/application"

type PrinterInterface struct {
	service application.PrinterService
}

func NewPrinterInterface(service application.PrinterService) PrinterInterface {
	return PrinterInterface{service}
}

func (p PrinterInterface) PrintReceipt(orderID string, name string, price int, thumbnail string, isbn string) error {
	return p.service.PrintReceipt(orderID, name, price, thumbnail, isbn)
}
