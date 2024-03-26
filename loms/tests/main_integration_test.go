package tests

import (
	"github.com/stretchr/testify/suite"
	"route256.ozon.ru/project/loms/tests/integration/repository/order"
	"route256.ozon.ru/project/loms/tests/integration/repository/stock"
	"testing"
)

func TestOrderSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(order.OrderPgRepositoryTestSuite))
}

func TestStockSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(stock.StockPgRepositoryTestSuite))
}
