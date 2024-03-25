package stock_pgs_repository

import (
	"context"
	"route256.ozon.ru/project/loms/internal/config"
	"route256.ozon.ru/project/loms/internal/model"
	pgs "route256.ozon.ru/project/loms/internal/repository/pgs"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type StockPgRepositoryTestSuite struct {
	suite.Suite
	repo   *StocksPgRepository
	dbPool *pgs.DB
	ctx    context.Context
}

func TestStockPgRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(StockPgRepositoryTestSuite))
}

func (t *StockPgRepositoryTestSuite) TestRepository() {
	sku := uint32(12345678)
	orderItem := []model.Item{{
		SKU:   sku,
		Count: 1,
	}}

	t.clearStocks(t.ctx, sku)
	var initialStockCount uint64 = 10

	t.prepareStocksForSku(t.ctx, sku, initialStockCount, 0)

	stockInfo, err := t.repo.GetBySku(t.ctx, sku)
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), initialStockCount, stockInfo.TotalCount)

	err = t.repo.Reserve(t.ctx, orderItem)
	assert.NoError(t.T(), err)

	count, err := t.repo.GetCountBySku(t.ctx, sku)
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), initialStockCount-1, count)

	err = t.repo.UnReserve(t.ctx, orderItem)
	assert.NoError(t.T(), err)

	count, err = t.repo.GetCountBySku(t.ctx, sku)
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), initialStockCount, count)

	err = t.repo.UnReserve(t.ctx, orderItem)
	assert.ErrorIs(t.T(), err, ErrInsufficientStocks)
}

func (t *StockPgRepositoryTestSuite) clearStocks(ctx context.Context, sku uint32) {
	_, err := t.repo.dbPool.Exec(ctx, "DELETE FROM stocks WHERE sku = $1", sku)
	require.NoError(t.T(), err)
}

func (t *StockPgRepositoryTestSuite) prepareStocksForSku(ctx context.Context, sku uint32, count uint64, reserved uint64) {
	_, err := t.repo.dbPool.Exec(ctx, "INSERT INTO stocks (sku, count, reserved) VALUES ($1, $2, $3)", sku, count, reserved)
	require.NoError(t.T(), err)
}

func (t *StockPgRepositoryTestSuite) SetupSuite() {
	ctx := context.Background()
	lomsConfig, err := config.GetConfig(ctx)
	require.NoError(t.T(), err)

	dbPool, err := pgs.ConnectToPgsDb(ctx, lomsConfig, true)
	require.NoError(t.T(), err)

	repo := NewStocksPgRepository(dbPool)
	require.NoError(t.T(), err)

	t.repo = repo
	t.dbPool = dbPool
	t.ctx = ctx
}

func (t *StockPgRepositoryTestSuite) TearDownSuite() {
	t.dbPool.Close()
}
