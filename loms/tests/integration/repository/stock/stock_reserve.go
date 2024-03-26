package stock

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"route256.ozon.ru/project/loms/internal/model"
	"testing"
)

func (testSuit *StockPgRepositoryTestSuite) TestStockReserve() {
	sku := uint32(1234567892)
	count := uint64(10)
	reserved := 0
	testCases := []struct {
		name      string
		orderItem []model.Item
	}{
		{
			name:      "Reserve item from stock 1",
			orderItem: []model.Item{{SKU: 1234567892, Count: 1}},
		},
		{
			name:      "Reserve item from stock 2",
			orderItem: []model.Item{{SKU: 1234567892, Count: 3}},
		},
		{
			name:      "Reserve item from stock 3",
			orderItem: []model.Item{{SKU: 1234567892, Count: 3}},
		},
	}
	_, err := testSuit.repo.DbPool.Exec(testSuit.ctx, "INSERT INTO stocks (sku, count, reserved) VALUES ($1, $2, $3)", sku, count, reserved)
	require.NoError(testSuit.T(), err)
	defer func() { testSuit.clearStocks(testSuit.ctx, sku) }()

	for _, tc := range testCases {
		testSuit.T().Run(tc.name, func(t *testing.T) {
			err := testSuit.repo.Reserve(testSuit.ctx, tc.orderItem)
			assert.NoError(t, err)
			dbCount, err := testSuit.repo.GetCountBySku(testSuit.ctx, sku)
			assert.Less(t, dbCount, count)
		})
	}
}
