package database

import payments "github.com/gpioblink/go-stripe-book-seller/pkg/payments/domain"

type MemoryRepository struct {
	orders []payments.Payment
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{[]payments.Payment{}}
}

func (m *MemoryRepository) Save(orderToSave *payments.Payment) error {
	m.orders = append(m.orders, *orderToSave)
	return nil
}

func (m MemoryRepository) ByOrderId(id string) (*payments.Payment, error) {
	for _, p := range m.orders {
		if p.OrderId() == id {
			return &p, nil
		}
	}

	return nil, payments.ErrNotFound
}

func (m MemoryRepository) ByPaymentId(id string) (*payments.Payment, error) {
	for _, p := range m.orders {
		if p.PaymentId() == id {
			return &p, nil
		}
	}

	return nil, payments.ErrNotFound
}

func (m MemoryRepository) DeleteByPaymentId(id string) error {
	var res []payments.Payment
	flag := false
	for _, p := range m.orders {
		if p.PaymentId() == id {
			flag = true
			continue
		}
		res = append(res, p)
	}
	if !flag {
		return payments.ErrNotFound
	}
	m.orders = append(res)
	return nil
}
