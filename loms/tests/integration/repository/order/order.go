package order

import (
	"context"
	"math/rand/v2"
	"route256.ozon.ru/project/loms/internal/config"
	"route256.ozon.ru/project/loms/internal/model"
	pgs "route256.ozon.ru/project/loms/internal/repository/pgs"
	order_pgs_repository "route256.ozon.ru/project/loms/internal/repository/pgs/order"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type OrderPgRepositoryTestSuite struct {
	suite.Suite
	repo   *order_pgs_repository.OrderPgsRepository
	dbPool *pgs.DB
	ctx    context.Context
}

func (t *OrderPgRepositoryTestSuite) TestFullCycle() {
	order := &model.Order{
		CreatedAt: time.Now(),
		Items: []model.Item{{
			SKU:   1234567891,
			Count: 1,
		}},
		User:   rand.Int64(),
		Status: model.New,
	}

	id, err := t.repo.OrderCreate(t.ctx, order)
	assert.NoError(t.T(), err)

	err = t.repo.OrderCancel(t.ctx, id)
	assert.NoError(t.T(), err)

	dbOrder, err := t.repo.OrderInfo(t.ctx, id)
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), model.Cancelled, dbOrder.Status)

	_, err = t.repo.DbPool.Exec(t.ctx, "DELETE FROM orders WHERE id = $1", id)
	require.NoError(t.T(), err)

	_, err = t.repo.OrderInfo(t.ctx, id)
	assert.ErrorIs(t.T(), err, order_pgs_repository.ErrOrderNotFound)
}

func (t *OrderPgRepositoryTestSuite) SetupSuite() {
	ctx := context.Background()
	lomsConfig, err := config.GetConfig(ctx)
	require.NoError(t.T(), err)

	dbPool, err := pgs.ConnectToPgsDb(ctx, lomsConfig, true, nil)
	require.NoError(t.T(), err)

	repo := order_pgs_repository.NewOrderPgsRepository(dbPool)
	require.NoError(t.T(), err)

	t.repo = repo
	t.dbPool = dbPool
	t.ctx = ctx
}

func (t *OrderPgRepositoryTestSuite) TearDownSuite() {
	t.dbPool.Close()
}
