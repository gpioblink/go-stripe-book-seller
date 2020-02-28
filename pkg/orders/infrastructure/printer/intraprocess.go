package printer

import (
	"github.com/gpioblink/go-stripe-book-seller/pkg/orders/domain/orders"
	"github.com/gpioblink/go-stripe-book-seller/pkg/printer/interfaces/private/intraprocess"
	productInterprocess "github.com/gpioblink/go-stripe-book-seller/pkg/shop/interfaces/private/intraprocess"
)

type IntraprocessService struct {
	intraprocessInterface        intraprocess.PrinterInterface
	productIntraprocessInterface productInterprocess.ProductInterface
}

func NewIntraprocessService(intraprocessInterface intraprocess.PrinterInterface, productIntraprocessInterface productInterprocess.ProductInterface) IntraprocessService {
	return IntraprocessService{intraprocessInterface, productIntraprocessInterface}
}

func (i IntraprocessService) PrintReceipt(order orders.Order) error {
	shopProduct, err := i.productIntraprocessInterface.ProductByID(string(order.Product().ID()))
	if err != nil {
		return err
	}

	return i.intraprocessInterface.PrintReceipt(string(order.ID()), order.Product().Name(), int(order.Product().Price().Cents()), shopProduct.ThumbnailUrl, shopProduct.Isbn)
}
