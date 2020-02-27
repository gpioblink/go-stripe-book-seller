package payments

type Payment struct {
	orderID   string
	paymentID string
}

func (p Payment) OrderId() string {
	return p.OrderId()
}

func (p Payment) PaymentId() string {
	return p.PaymentId()
}

func NewPayment(orderId string, paymentId string) (Payment, error) {
	return Payment{orderId, paymentId}, nil
}
