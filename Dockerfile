FROM golang:1.9

RUN go get github.com/golang/dep/cmd/dep

WORKDIR /go/src/github.com/gpioblink/go-stripe-book-seller
COPY . .

RUN dep ensure
RUN go get github.com/cespare/reflex
