package order_usecase

import (
	"context"
	"errors"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"route256.ozon.ru/project/loms/internal/model"
	stock_repository "route256.ozon.ru/project/loms/internal/repository/stock"
	order_usecase "route256.ozon.ru/project/loms/internal/service/loms/mock"
	"testing"
)

func TestService_OrderCancel(t *testing.T) {
	ctx := context.Background()
	type fields struct {
		mockOrderRepository *order_usecase.OrderRepositoryMock
		mockStockRepository *order_usecase.StockRepositoryMock
	}
	testCases := []struct {
		name          string
		orderID       int64
		mockSetup     func(s *fields, orderID int64)
		expectedError error
	}{
		{
			name:    "Invalid Order ID",
			orderID: 0,
			mockSetup: func(s *fields, orderID int64) {
			},
			expectedError: ErrOrderInvalid,
		},
		{
			name:    "Order Not Found",
			orderID: 1,
			mockSetup: func(s *fields, orderID int64) {
				s.mockOrderRepository.OrderInfoMock.Expect(ctx, orderID).Return(nil, errors.New("order not found"))
			},
			expectedError: ErrOrderNotFound,
		},
		{
			name:    "Order Already Canceled",
			orderID: 2,
			mockSetup: func(s *fields, orderID int64) {
				order := &model.Order{Status: model.Cancelled}
				s.mockOrderRepository.OrderInfoMock.Expect(ctx, orderID).Return(order, nil)
			},
			expectedError: ErrOrderAlreadyCanceled,
		},
		{
			name:    "Unable to Unreserve Stock",
			orderID: 3,
			mockSetup: func(s *fields, orderID int64) {
				order := &model.Order{Id: orderID, Items: []model.Item{{SKU: 1, Count: 5}}}
				s.mockStockRepository.GetBySkuMock.Expect(ctx, 1).Return(stock_repository.Product{SKU: 1, TotalCount: 5, Reserved: 0}, nil)
				s.mockOrderRepository.OrderInfoMock.Expect(ctx, orderID).Return(order, nil)
			},
			expectedError: ErrOrderCantUnReserve,
		},
		{
			name:    "Successful Order Cancellation",
			orderID: 4,
			mockSetup: func(s *fields, orderID int64) {
				order := &model.Order{Id: 4, Items: []model.Item{{SKU: 2, Count: 5}}}
				s.mockStockRepository.GetBySkuMock.Expect(ctx, 2).Return(stock_repository.Product{SKU: 2, TotalCount: 5, Reserved: 6}, nil)
				s.mockStockRepository.UnReserveMock.Expect(ctx, order.Items).Return(nil)

				s.mockOrderRepository.OrderInfoMock.Expect(ctx, orderID).Return(order, nil)
				s.mockOrderRepository.OrderCancelMock.Expect(ctx, orderID).Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			mockOrderRepository := order_usecase.NewOrderRepositoryMock(mc)
			mockStockRepository := order_usecase.NewStockRepositoryMock(mc)
			f := &fields{
				mockOrderRepository,
				mockStockRepository,
			}
			t.Parallel()
			tc.mockSetup(f, tc.orderID)
			service := NewService(mockOrderRepository, mockStockRepository)
			err := service.OrderCancel(ctx, tc.orderID)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
