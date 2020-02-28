package application

import "fmt"

type printService interface {
	PrintTitle(query string) error
	PrintLine(query string) error
	PrintPicture(srcUrl string) error
	PrintBarCode(query string) error
	Cut() error
}

type PrinterService struct {
	printService printService
}

func NewPrinterService(service printService) PrinterService {
	return PrinterService{service}
}

func (p PrinterService) PrintReceipt(orderID string, name string, price int,
	thumbnail string, isbn string) error {

	err := p.printService.PrintTitle("受け取り票")
	if err != nil {
		return err
	}

	err = p.printService.PrintTitle(fmt.Sprintf("購入番号: %s", orderID))
	if err != nil {
		return err
	}
	err = p.printService.PrintTitle(fmt.Sprintf("商品名: %s", name))
	if err != nil {
		return err
	}
	err = p.printService.PrintTitle(fmt.Sprintf("金額: %d円", price))
	if err != nil {
		return err
	}

	err = p.printService.PrintPicture(thumbnail)
	if err != nil {
		return err
	}

	err = p.printService.PrintBarCode(isbn)
	if err != nil {
		return err
	}

	return nil
}
