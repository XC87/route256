FROM golang:1.22-alpine as builder

WORKDIR /app

COPY notifier/go.mod go.mod
COPY notifier/go.sum go.sum
COPY ./pkg/ ../pkg

RUN go mod download

COPY ./notifier/ .

RUN go build -o bin/notifier ./cmd/notifier

FROM alpine

COPY --from=builder /app/bin/notifier /notifier

CMD [ "/notifier" ]
