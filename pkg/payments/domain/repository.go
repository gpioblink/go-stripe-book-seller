package payments

import "errors"

var ErrNotFound = errors.New("payment not found")

type Repository interface {
	Save(*Payment) error
	DeleteByPaymentID(string) error
	ByOrderID(string) (*Payment, error)
	ByPaymentID(string) (*Payment, error)
}
