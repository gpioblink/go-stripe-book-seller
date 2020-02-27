package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	common_http "github.com/gpioblink/go-stripe-book-seller/pkg/common/http"
	"github.com/gpioblink/go-stripe-book-seller/pkg/orders/application"
	"github.com/gpioblink/go-stripe-book-seller/pkg/orders/domain/orders"
)

func AddRoutes(router *chi.Mux, service application.OrdersService, repository orders.Repository) {
	resource := ordersResource{service, repository}
	router.Post("/orders/{id}/paid", resource.PostPaid)
}

type ordersResource struct {
	service application.OrdersService

	repository orders.Repository
}

func (o ordersResource) PostPaid(w http.ResponseWriter, r *http.Request) {
	cmd := application.MarkOrderAsPaidCommand{
		OrderID: orders.ID(chi.URLParam(r, "id")),
	}

	if err := o.service.MarkOrderAsPaid(cmd); err != nil {
		_ = render.Render(w, r, common_http.ErrInternal(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
