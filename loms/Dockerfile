FROM golang:1.22-alpine as builder

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .

RUN go build -o bin/loms-server ./cmd/loms-server
#RUN GOBIN=/app/bin go install github.com/go-swagger/go-swagger/cmd/swagger@v0.30.5

FROM alpine

COPY --from=builder /app/bin/loms-server /loms-server
#COPY --from=builder /app/bin/swagger /swagger
COPY --from=builder /app/stock-data.json /stock-data.json
COPY --from=builder /app/swagger-ui /swagger-ui
COPY --from=builder /app/api/openapiv2/loms.swagger.json /api/openapiv2/loms.swagger.json

#CMD ["/swagger", "serve", "/api/openapiv2/loms.swagger.json", "-F", "swagger", "-p", "8085", "--no-open" ]
ENTRYPOINT [ "/loms-server" ]
