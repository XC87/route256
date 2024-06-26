// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package queries

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
)

type Querier interface {
	CreateOrder(ctx context.Context, arg CreateOrderParams) (int64, error)
	CreateOrderItems(ctx context.Context, arg CreateOrderItemsParams) error
	GetBySku(ctx context.Context, sku int64) (GetBySkuRow, error)
	GetOrderAll(ctx context.Context) ([]GetOrderAllRow, error)
	GetOrderById(ctx context.Context, arg GetOrderByIdParams) ([]GetOrderByIdRow, error)
	UpdateCountBySku(ctx context.Context, arg UpdateCountBySkuParams) (pgconn.CommandTag, error)
	UpdateOrder(ctx context.Context, arg UpdateOrderParams) error
	UpdateOrderStatus(ctx context.Context, arg UpdateOrderStatusParams) error
	UpdateReserveBySku(ctx context.Context, arg UpdateReserveBySkuParams) (pgconn.CommandTag, error)
}

var _ Querier = (*Queries)(nil)
