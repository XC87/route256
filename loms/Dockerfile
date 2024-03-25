FROM golang:1.22-alpine as builder

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .

RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN go build -o bin/loms-server ./cmd/loms-server

FROM alpine

COPY --from=builder /app/migrations /loms/migrations
COPY --from=builder /app/start.sh /start.sh
COPY --from=builder /app/bin/loms-server /loms-server
COPY --from=builder /go/bin/goose /bin/goose
#COPY --from=builder /app/bin/swagger /swagger
COPY --from=builder /app/stock-data.json /stock-data.json
COPY --from=builder /app/swagger-ui /swagger-ui
COPY --from=builder /app/api/openapiv2/loms.swagger.json /api/openapiv2/loms.swagger.json

ENTRYPOINT ["/start.sh" ]
