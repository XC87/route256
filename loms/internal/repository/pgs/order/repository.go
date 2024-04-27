package order_pgs_repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
	"route256.ozon.ru/project/loms/internal/model"
	pgs "route256.ozon.ru/project/loms/internal/repository/pgs"
	"route256.ozon.ru/project/loms/internal/repository/pgs/queries"
	"time"
)

type OrderPgsRepository struct {
	ShardsPool *pgs.Manager
}

func NewOrderPgsRepository(dbPool *pgs.Manager) *OrderPgsRepository {
	return &OrderPgsRepository{
		ShardsPool: dbPool,
	}
}

var (
	ErrOrderNotFound = errors.New("order not found")
)

func (repo *OrderPgsRepository) OrderCreate(ctx context.Context, order *model.Order) (int64, error) {
	shardIndex := repo.ShardsPool.AutoPickIndex(order.User)
	dbPool, err := repo.ShardsPool.Pick(shardIndex)
	if err != nil {
		return 0, fmt.Errorf("cant pick db: %w", err)
	}

	tx, err := dbPool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			zap.L().Info("cannot rollback transaction")
		}
	}(tx, ctx)

	q := queries.New(tx)

	orderCreateParams := queries.CreateOrderParams{
		Id: shardIndex,
		CreatedAt: pgtype.Timestamp{
			Time:  time.Now(),
			Valid: true,
		},
		UserID: order.User,
		Status: model.MapStatusToId(order.Status),
	}
	orderId, err := q.CreateOrder(ctx, orderCreateParams)
	if err != nil {
		return 0, fmt.Errorf("cannot create order: %w", err)
	}

	for _, item := range order.Items {
		createOrderItemParam := queries.CreateOrderItemsParams{
			OrderID: orderId,
			Sku:     int64(item.SKU),
			Count:   int64(item.Count),
		}
		err = q.CreateOrderItems(ctx, createOrderItemParam)
		if err != nil {
			return 0, fmt.Errorf("cannot create order item: %w", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, fmt.Errorf("cannot commit order create transaction: %w", err)
	}

	return orderId, nil
}

func (repo *OrderPgsRepository) OrderUpdate(ctx context.Context, order *model.Order) error {
	shardIndex := repo.ShardsPool.AutoPickIndex(order.User)
	dbPool, err := repo.ShardsPool.Pick(shardIndex)
	if err != nil {
		return fmt.Errorf("cant pick db: %w", err)
	}

	tx, err := dbPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			zap.L().Info("cannot rollback transaction")
		}
	}(tx, ctx)

	q := queries.New(tx)

	params := queries.UpdateOrderParams{
		Status:    model.MapStatusToId(order.Status),
		UserID:    order.User,
		CreatedAt: pgtype.Timestamp{Time: order.CreatedAt, Valid: true},
		UpdatedAt: pgtype.Timestamp{Time: order.UpdatedAt, Valid: true},
		ID:        order.Id,
	}
	err = q.UpdateOrder(ctx, params)
	if err != nil {
		return fmt.Errorf("cannot set update order: %w", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("cannot commit update of order transaction: %w", err)
	}
	return nil
}

func (repo *OrderPgsRepository) SetStatus(ctx context.Context, id int64, userId int64, status model.OrderStatus) error {
	shardIndex := repo.ShardsPool.AutoPickIndex(userId)
	dbPool, err := repo.ShardsPool.Pick(shardIndex)
	if err != nil {
		return fmt.Errorf("cant pick db: %w", err)
	}

	q := dbPool.GetUpdateQueries(ctx)
	params := queries.UpdateOrderStatusParams{
		Status: model.MapStatusToId(status),
		ID:     id,
	}
	err = q.UpdateOrderStatus(ctx, params)
	if err != nil {
		return fmt.Errorf("cannot set status of order: %w", err)
	}

	return nil
}

func (repo *OrderPgsRepository) OrderInfo(ctx context.Context, id int64, userId int64) (*model.Order, error) {
	shardIndex := repo.ShardsPool.AutoPickIndex(userId)
	dbPool, err := repo.ShardsPool.Pick(shardIndex)
	if err != nil {
		return nil, fmt.Errorf("cant pick db: %w", err)
	}

	q := dbPool.GetSelectQueries(ctx)
	getOrderByIdParams := queries.GetOrderByIdParams{
		ID:     id,
		UserID: userId,
	}
	orderWithItems, err := q.GetOrderById(ctx, getOrderByIdParams)
	if err != nil || orderWithItems == nil {
		return nil, ErrOrderNotFound
	}
	orderInfo := orderWithItems[0].Order
	order := &model.Order{
		CreatedAt: orderInfo.CreatedAt.Time,
		Items:     make([]model.Item, 0, len(orderWithItems)),
		Id:        orderInfo.ID,
		User:      orderInfo.UserID,
		Status:    model.MapIdToStatus(orderInfo.Status),
	}
	for _, orderItemFromDb := range orderWithItems {
		orderItem := model.Item{
			SKU:   uint32(orderItemFromDb.OrderItem.Sku),
			Count: uint64(orderItemFromDb.OrderItem.Count),
		}
		order.Items = append(order.Items, orderItem)
	}

	return order, nil
}

func (repo *OrderPgsRepository) OrderInfoAll(ctx context.Context) ([]*model.Order, error) {
	dbPoolList := repo.ShardsPool.GerAllShards()
	orderList := make([]*model.Order, 0)
	for _, dbPool := range dbPoolList {
		q := dbPool.GetSelectQueries(ctx)
		orderLostWithItems, err := q.GetOrderAll(ctx)
		if err != nil || orderLostWithItems == nil {
			return nil, ErrOrderNotFound
		}
		lastId := int64(0)
		for _, orderWithItems := range orderLostWithItems {
			orderInfo := orderWithItems.Order

			order := &model.Order{
				CreatedAt: orderInfo.CreatedAt.Time,
				Items:     make([]model.Item, 0),
				Id:        orderInfo.ID,
				User:      orderInfo.UserID,
				Status:    model.MapIdToStatus(orderInfo.Status),
			}
			orderItem := model.Item{
				SKU:   uint32(orderWithItems.OrderItem.Sku),
				Count: uint64(orderWithItems.OrderItem.Count),
			}
			order.Items = append(order.Items, orderItem)

			if lastId != orderInfo.ID {
				orderList = append(orderList, order)
			}
			lastId = orderInfo.ID
		}
	}

	return orderList, nil
}

func (repo *OrderPgsRepository) OrderPay(ctx context.Context, id int64, userId int64) error {
	return repo.SetStatus(ctx, id, userId, model.Paid)
}

func (repo *OrderPgsRepository) OrderCancel(ctx context.Context, id int64, userId int64) error {
	return repo.SetStatus(ctx, id, userId, model.Cancelled)
}
