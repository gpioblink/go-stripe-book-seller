package shop

import (
	shop_app "github.com/gpioblink/go-stripe-book-seller/pkg/shop/application"
	"github.com/stripe/stripe-go"
)

func LoadShopFixtures(productsService shop_app.ProductsService) error {
	err := productsService.AddProduct(shop_app.AddProductCommand{
		ID:            "1",
		Name:          "ボクの電子工作ノート",
		Description:   "貴重な実測データと使えるアイデアを満載",
		ThumbnailUrl:  "http://books.google.com/books/content?id=ikVgLwEACAAJ&printsec=frontcover&img=1&zoom=1&source=gbs_api",
		Isbn:          "9784899773078",
		PriceCents:    2457,
		PriceCurrency: string(stripe.CurrencyJPY),
	})
	if err != nil {
		return err
	}

	return productsService.AddProduct(shop_app.AddProductCommand{
		ID:            "2",
		Name:          "Raspberry Piクックブック",
		Description:   "Raspberry Piのすべてを使いこなす216のレシピ。初期設定、LinuxとPythonの基礎、ネットワークの設定、GPIOの使い方などの基本的な情報から、各種センサー、モーターとの組み合わせなどを豊富なサンプルコードと合わせて解説。Arduinoとの連携についても詳しく紹介した決定版!",
		ThumbnailUrl:  "http://books.google.com/books/content?id=Y46uoQEACAAJ&printsec=frontcover&img=1&zoom=1&source=gbs_api",
		Isbn:          "9784873116907",
		PriceCents:    3400,
		PriceCurrency: string(stripe.CurrencyJPY),
	})
}
