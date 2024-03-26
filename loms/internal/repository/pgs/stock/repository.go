package stock_pgs_repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"
	"route256.ozon.ru/project/loms/internal/model"
	"route256.ozon.ru/project/loms/internal/repository/pgs"
	"route256.ozon.ru/project/loms/internal/repository/pgs/queries"
	"strings"
)

type StocksPgRepository struct {
	DbPool *pgs.DB
}

var ErrInsufficientStocks = errors.New("insufficient stocks")

func NewStocksPgRepository(dbPool *pgs.DB) *StocksPgRepository {
	return &StocksPgRepository{
		DbPool: dbPool,
	}
}

func (repo *StocksPgRepository) Reserve(ctx context.Context, items []model.Item) error {
	tx, err := repo.DbPool.Begin(ctx)
	if err != nil {
		log.Err(err).Msg("cannot begin transaction for reserving stocks")
		return fmt.Errorf("cannot begin transaction for reserving stocks: %w", err)
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Err(err).Msg("cannot rollback transaction")
		}
	}(tx, ctx)

	q := queries.New(tx)

	for _, prod := range items {
		params := queries.UpdateReserveBySkuParams{
			Reserved: int64(prod.Count),
			Sku:      int64(prod.SKU),
		}
		res, err := q.UpdateReserveBySku(ctx, params)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if strings.Contains(pgErr.Message, "nonnegative") {
					return ErrInsufficientStocks
				}
				return err
			}
		}
		if res.RowsAffected() != 1 {
			return ErrInsufficientStocks
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Err(err).Msg("cannot commit transaction for reserving stocks")
		return fmt.Errorf("cannot commit transaction for reserving stocks: %w", err)
	}

	return nil
}
func (repo *StocksPgRepository) UnReserve(ctx context.Context, items []model.Item) error {
	minusItem := make([]model.Item, len(items))
	copy(minusItem, items)
	for i := range minusItem {
		minusItem[i].Count = -minusItem[i].Count
	}
	return repo.Reserve(ctx, minusItem)
}

func (repo *StocksPgRepository) GetCountBySku(ctx context.Context, sku uint32) (uint64, error) {
	q := repo.DbPool.GetSelectQueries(ctx)
	stockInfo, err := q.GetBySku(ctx, int64(sku))
	if err != nil {
		return 0, err
	}

	return uint64(stockInfo.Count - stockInfo.Reserved), nil
}

func (repo *StocksPgRepository) GetBySku(ctx context.Context, sku uint32) (*model.ProductStock, error) {
	q := repo.DbPool.GetSelectQueries(ctx)
	stockInfo, err := q.GetBySku(ctx, int64(sku))
	if err != nil {
		return nil, err
	}

	return &model.ProductStock{
		SKU:        sku,
		TotalCount: uint64(stockInfo.Count),
		Reserved:   uint64(stockInfo.Reserved),
	}, nil
}
