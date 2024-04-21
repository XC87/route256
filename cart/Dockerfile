FROM golang:1.22-alpine as builder

WORKDIR /app

COPY cart/go.mod go.mod
COPY cart/go.sum go.sum
COPY ./pkg/ ../pkg

RUN go mod download

COPY ./cart/ .

RUN go build -o bin/cart-server ./cmd/cart-server

FROM alpine

COPY --from=builder /app/bin/cart-server /cart-server

CMD [ "/cart-server" ]
