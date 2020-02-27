package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/gpioblink/go-stripe-book-seller/pkg/payments/application"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/webhook"
	"io/ioutil"
	"net/http"
	"os"
)

func NewCheckoutInterface(router *chi.Mux, service application.PaymentsService, secretKey string) {
	resource := checkoutResource{secretKey, service}
	router.Post("/webhook", resource.ProcessEvent)
}

type checkoutResource struct {
	secretKey string
	service   application.PaymentsService
}

func (c checkoutResource) ProcessEvent(w http.ResponseWriter, req *http.Request) {
	const MaxBodyBytes = int64(65536)
	req.Body = http.MaxBytesReader(w, req.Body, MaxBodyBytes)
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// Pass the request body & Stripe-Signature header to ConstructEvent, along with the webhook signing key
	// You can find your endpoint's secret in your webhook settings
	endpointSecret := c.secretKey
	event, err := webhook.ConstructEvent(body, req.Header.Get("Stripe-Signature"), endpointSecret)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
		w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}

	// Handle the checkout.session.completed event
	if event.Type == "checkout.session.completed" {
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Fulfill the purchase...
		err = c.service.PostOrderPayment(session.ID)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error Confirming Order: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
