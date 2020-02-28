package payments

type Payment struct {
	orderID   string
	paymentID string
}

func (p Payment) OrderId() string {
	return p.orderID
}

func (p Payment) PaymentId() string {
	return p.paymentID
}

func NewPayment(orderId string, paymentId string) (Payment, error) {
	return Payment{orderId, paymentId}, nil
}
