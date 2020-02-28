package http

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	common_http "github.com/gpioblink/go-stripe-book-seller/pkg/common/http"
	"github.com/gpioblink/go-stripe-book-seller/pkg/payments/application"
	"net/http"
)

func AddRoutes(router *chi.Mux, service application.PaymentsService) {
	resource := paymentResource{service}
	router.Get("/payments/{id}", resource.Get)
}

type paymentResource struct {
	service application.PaymentsService
}

type PaymentView struct {
	PaymentID string `json:"PaymentID"`
}

func (p paymentResource) Get(w http.ResponseWriter, r *http.Request) {
	paymentId, err := p.service.GetPaymentID(chi.URLParam(r, "id"))
	if err != nil {
		_ = render.Render(w, r, common_http.ErrBadRequest(err))
		return
	}

	render.Respond(w, r, PaymentView{paymentId})
}
