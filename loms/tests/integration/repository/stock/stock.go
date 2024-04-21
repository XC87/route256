package stock

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"route256.ozon.ru/project/loms/internal/config"
	"route256.ozon.ru/project/loms/internal/model"
	pgs "route256.ozon.ru/project/loms/internal/repository/pgs"
	stock_pgs_repository "route256.ozon.ru/project/loms/internal/repository/pgs/stock"
)

type StockPgRepositoryTestSuite struct {
	suite.Suite
	repo   *stock_pgs_repository.StocksPgRepository
	DbPool *pgs.DB
	ctx    context.Context
}

func (t *StockPgRepositoryTestSuite) TestRepository() {
	sku := uint32(1234567891)
	orderItem := []model.Item{{
		SKU:   sku,
		Count: 1,
	}}
	defer func() { t.clearStocks(t.ctx, sku) }()

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
	assert.ErrorIs(t.T(), err, stock_pgs_repository.ErrInsufficientStocks)
}

func (t *StockPgRepositoryTestSuite) clearStocks(ctx context.Context, sku uint32) {
	_, err := t.repo.DbPool.Exec(ctx, "DELETE FROM stocks WHERE sku = $1", sku)
	require.NoError(t.T(), err)
}

func (t *StockPgRepositoryTestSuite) prepareStocksForSku(ctx context.Context, sku uint32, count uint64, reserved uint64) {
	_, err := t.repo.DbPool.Exec(ctx, "INSERT INTO stocks (sku, count, reserved) VALUES ($1, $2, $3)", sku, count, reserved)
	require.NoError(t.T(), err)
}

func (t *StockPgRepositoryTestSuite) SetupSuite() {
	ctx := context.Background()
	lomsConfig, err := config.GetConfig(ctx)
	require.NoError(t.T(), err)

	dbPool, err := pgs.ConnectToPgsDb(ctx, lomsConfig, true, nil)
	require.NoError(t.T(), err)

	repo := stock_pgs_repository.NewStocksPgRepository(dbPool)
	require.NoError(t.T(), err)

	t.repo = repo
	t.DbPool = dbPool
	t.ctx = ctx
}

func (t *StockPgRepositoryTestSuite) TearDownSuite() {
	t.DbPool.Close()
}
