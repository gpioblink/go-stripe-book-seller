package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	common_http "github.com/gpioblink/go-stripe-book-seller/pkg/common/http"
	"github.com/gpioblink/go-stripe-book-seller/pkg/common/price"
	"github.com/gpioblink/go-stripe-book-seller/pkg/shop/domain/products"
)

func AddRoutes(router *chi.Mux, productsReadModel productsReadModel) {
	resource := productsResource{productsReadModel}
	router.Get("/products", resource.GetAll)
}

type productsReadModel interface {
	AllProducts() ([]products.Product, error)
}

type productView struct {
	ID string `json:"id"`

	Name         string `json:"name"`
	Description  string `json:"description"`
	ThumbnailUrl string `json:"thumbnail"`
	Isbn         string `json:"isbn"`

	Price priceView `json:"price"`
}

type priceView struct {
	Cents    uint   `json:"cents"`
	Currency string `json:"currency"`
}

func priceViewFromPrice(p price.Price) priceView {
	return priceView{p.Cents(), p.Currency()}
}

type productsResource struct {
	readModel productsReadModel
}

func (p productsResource) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := p.readModel.AllProducts()
	if err != nil {
		_ = render.Render(w, r, common_http.ErrInternal(err))
		return
	}

	view := []productView{}
	for _, product := range products {
		view = append(view, productView{
			string(product.ID()),
			product.Name(),
			product.Description(),
			product.ThumbnailUrl(),
			product.Isbn(),
			priceViewFromPrice(product.Price()),
		})
	}

	render.Respond(w, r, view)
}
