package order

import (
	"github.com/stretchr/testify/require"
	"math/rand/v2"
	"route256.ozon.ru/project/loms/internal/model"
	"testing"
	"time"
)

func (testSuit *OrderPgRepositoryTestSuite) TestOrderCreate() {
	tests := []struct {
		name  string
		order *model.Order
		clear func(t *OrderPgRepositoryTestSuite, id int64)
	}{
		{
			name: "Order create 1",
			order: &model.Order{
				CreatedAt: time.Now(),
				Items: []model.Item{{
					SKU:   1234567891,
					Count: 1,
				}},
				User:   rand.Int64(),
				Status: model.New,
			},
		},
		{
			name: "Order create 2",
			order: &model.Order{
				CreatedAt: time.Now(),
				Items: []model.Item{{
					SKU:   1234567891,
					Count: 1,
				}},
				User:   rand.Int64(),
				Status: model.New,
			},
		},
		{
			name: "Order create 3",
			order: &model.Order{
				CreatedAt: time.Now(),
				Items: []model.Item{{
					SKU:   1234567891,
					Count: 1,
				}},
				User:   rand.Int64(),
				Status: model.New,
			},
		},
	}

	for _, tc := range tests {
		testSuit.T().Run(tc.name, func(t *testing.T) {
			t.Parallel()
			id, err := testSuit.repo.OrderCreate(testSuit.ctx, tc.order)
			require.NoError(t, err)
			require.NotEmpty(t, id)
			index := testSuit.repo.ShardsPool.AutoPickIndex(tc.order.User)
			dbPool, _ := testSuit.repo.ShardsPool.Pick(index)
			dbPool.Exec(testSuit.ctx, "DELETE FROM orders WHERE id = $1", id)
			require.NoError(t, err)
		})
	}
}
