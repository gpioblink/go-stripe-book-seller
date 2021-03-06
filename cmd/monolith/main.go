package main

// yep, it's a bit ugly :(
import (
	orders_infra_printer "github.com/gpioblink/go-stripe-book-seller/pkg/orders/infrastructure/printer"
	payments_infra_database "github.com/gpioblink/go-stripe-book-seller/pkg/payments/infrastructure/database"
	payments_infra_stripe "github.com/gpioblink/go-stripe-book-seller/pkg/payments/infrastructure/payments"
	payments_interfaces_http "github.com/gpioblink/go-stripe-book-seller/pkg/payments/interfaces/public/http"
	payments_interfaces_stripe_http "github.com/gpioblink/go-stripe-book-seller/pkg/payments/interfaces/stripe/http"
	printer_app "github.com/gpioblink/go-stripe-book-seller/pkg/printer/application"
	printer_infra_lineprinter "github.com/gpioblink/go-stripe-book-seller/pkg/printer/infrastructure/printer"
	printer_interfaces_intraprocess "github.com/gpioblink/go-stripe-book-seller/pkg/printer/interfaces/private/intraprocess"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/gpioblink/go-stripe-book-seller/pkg/common/cmd"
	orders_app "github.com/gpioblink/go-stripe-book-seller/pkg/orders/application"
	orders_infra_orders "github.com/gpioblink/go-stripe-book-seller/pkg/orders/infrastructure/orders"
	orders_infra_payments "github.com/gpioblink/go-stripe-book-seller/pkg/orders/infrastructure/payments"
	orders_infra_product "github.com/gpioblink/go-stripe-book-seller/pkg/orders/infrastructure/shop"
	orders_interfaces_intraprocess "github.com/gpioblink/go-stripe-book-seller/pkg/orders/interfaces/private/intraprocess"
	orders_interfaces_http "github.com/gpioblink/go-stripe-book-seller/pkg/orders/interfaces/public/http"
	payments_app "github.com/gpioblink/go-stripe-book-seller/pkg/payments/application"
	payments_infra_orders "github.com/gpioblink/go-stripe-book-seller/pkg/payments/infrastructure/orders"
	payments_interfaces_intraprocess "github.com/gpioblink/go-stripe-book-seller/pkg/payments/interfaces/intraprocess"
	"github.com/gpioblink/go-stripe-book-seller/pkg/shop"
	shop_app "github.com/gpioblink/go-stripe-book-seller/pkg/shop/application"
	shop_infra_product "github.com/gpioblink/go-stripe-book-seller/pkg/shop/infrastructure/products"
	shop_interfaces_intraprocess "github.com/gpioblink/go-stripe-book-seller/pkg/shop/interfaces/private/intraprocess"
	shop_interfaces_http "github.com/gpioblink/go-stripe-book-seller/pkg/shop/interfaces/public/http"
)

func main() {
	log.Println("Starting monolith")
	ctx := cmd.Context()

	ordersToPay := make(chan payments_interfaces_intraprocess.OrderToProcess)
	router, paymentsInterface := createMonolith(ordersToPay)
	go paymentsInterface.Run()

	server := &http.Server{Addr: os.Getenv("SHOP_MONOLITH_BIND_ADDR"), Handler: router}
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()
	log.Printf("Monolith is listening on %s", server.Addr)

	<-ctx.Done()
	log.Println("Closing monolith")

	if err := server.Close(); err != nil {
		panic(err)
	}

	close(ordersToPay)
	paymentsInterface.Close()
}

func createMonolith(ordersToPay chan payments_interfaces_intraprocess.OrderToProcess) (*chi.Mux, payments_interfaces_intraprocess.PaymentsInterface) {
	shopProductRepo := shop_infra_product.NewMemoryRepository()
	shopProductsService := shop_app.NewProductsService(shopProductRepo, shopProductRepo)
	shopProductIntraprocessInterface := shop_interfaces_intraprocess.NewProductInterface(shopProductRepo)

	printerLinePrinterRepo := printer_infra_lineprinter.NewPrinterInterface()
	printerService := printer_app.NewPrinterService(printerLinePrinterRepo)
	printerIntraprocessInterface := printer_interfaces_intraprocess.NewPrinterInterface(printerService)

	ordersRepo := orders_infra_orders.NewMemoryRepository()
	orderService := orders_app.NewOrdersService(
		orders_infra_product.NewIntraprocessService(shopProductIntraprocessInterface),
		orders_infra_payments.NewIntraprocessService(ordersToPay),
		orders_infra_printer.NewIntraprocessService(printerIntraprocessInterface, shopProductIntraprocessInterface),
		ordersRepo,
	)
	ordersIntraprocessInterface := orders_interfaces_intraprocess.NewOrdersInterface(orderService)

	paymentsRepo := payments_infra_database.NewMemoryRepository()
	paymentsService := payments_app.NewPaymentsService(
		payments_infra_orders.NewIntraprocessService(ordersIntraprocessInterface),
		payments_infra_stripe.NewStripeService(os.Getenv("SHOP_MONOLITH_STRIPE_SECRET_KEY"), paymentsRepo),
		paymentsRepo,
	)
	paymentsIntraprocessInterface := payments_interfaces_intraprocess.NewPaymentsInterface(ordersToPay, paymentsService)

	if err := shop.LoadShopFixtures(shopProductsService); err != nil {
		panic(err)
	}

	r := cmd.CreateRouter()

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

	shop_interfaces_http.AddRoutes(r, shopProductRepo)
	orders_interfaces_http.AddRoutes(r, orderService, ordersRepo)
	payments_interfaces_http.AddRoutes(r, paymentsService)
	payments_interfaces_stripe_http.NewCheckoutInterface(r, paymentsService, os.Getenv("SHOP_MONOLITH_STRIPE_WEBHOOK_SECRET"))

	return r, paymentsIntraprocessInterface
}
