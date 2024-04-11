package tests

import (
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"
	"route256.ozon.ru/project/loms/tests/integration/repository/order"
	"route256.ozon.ru/project/loms/tests/integration/repository/stock"
	"testing"
)

func TestOrderSuite(t *testing.T) {
	defer goleak.VerifyNone(t)
	t.Parallel()
	suite.Run(t, new(order.OrderPgRepositoryTestSuite))
}

func TestStockSuite(t *testing.T) {
	defer goleak.VerifyNone(t)
	t.Parallel()
	suite.Run(t, new(stock.StockPgRepositoryTestSuite))
}
