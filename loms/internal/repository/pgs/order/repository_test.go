package order_pgs_repository

import (
	"context"
	"route256.ozon.ru/project/loms/internal/config"
	"route256.ozon.ru/project/loms/internal/model"
	pgs "route256.ozon.ru/project/loms/internal/repository/pgs"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type OrderPgRepositoryTestSuite struct {
	suite.Suite
	repo   *OrderPgsRepository
	dbPool *pgs.DB
	ctx    context.Context
}

func TestOrderPgRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(OrderPgRepositoryTestSuite))
}

func (t *OrderPgRepositoryTestSuite) TestRepository() {
	order := &model.Order{
		CreatedAt: time.Now(),
		Items: []model.Item{{
			SKU:   77329741,
			Count: 123,
		}},
		User:   1,
		Status: model.New,
	}

	id, err := t.repo.OrderCreate(t.ctx, order)
	assert.NoError(t.T(), err)

	err = t.repo.OrderCancel(t.ctx, id)
	assert.NoError(t.T(), err)

	dbOrder, err := t.repo.OrderInfo(t.ctx, id)
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), model.Cancelled, dbOrder.Status)

	_, err = t.repo.dbPool.Exec(t.ctx, "DELETE FROM orders WHERE id = $1", id)
	require.NoError(t.T(), err)

	_, err = t.repo.OrderInfo(t.ctx, id)
	assert.ErrorIs(t.T(), err, ErrOrderNotFound)
}

func (t *OrderPgRepositoryTestSuite) SetupSuite() {
	ctx := context.Background()
	lomsConfig, err := config.GetConfig(ctx)
	require.NoError(t.T(), err)

	dbPool, err := pgs.ConnectToPgsDb(ctx, lomsConfig, false)
	require.NoError(t.T(), err)

	repo := NewOrderPgsRepository(dbPool)
	require.NoError(t.T(), err)

	t.repo = repo
	t.dbPool = dbPool
	t.ctx = ctx
}

func (t *OrderPgRepositoryTestSuite) TearDownSuite() {
	t.dbPool.Close()
}
