FROM golang:1.22-alpine as builder

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .

RUN go build -o bin/cart-server ./cmd/cart-server

FROM alpine

COPY --from=builder /app/bin/cart-server /cart-server

CMD [ "/cart-server" ]
