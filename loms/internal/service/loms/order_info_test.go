package order_usecase

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"route256.ozon.ru/project/loms/internal/model"
	order_usecase "route256.ozon.ru/project/loms/internal/service/loms/mock"
	"testing"
)

func TestService_OrderInfo(t *testing.T) {
	ctx := context.Background()
	type fields struct {
		mockOrderRepository *order_usecase.OrderRepositoryMock
		mockStockRepository *order_usecase.StockRepositoryMock
	}
	testCases := []struct {
		name          string
		orderID       int64
		mockSetup     func(f *fields, orderID int64)
		expectedError error
	}{
		{
			name:    "Valid Order ID",
			orderID: 1,
			mockSetup: func(f *fields, orderID int64) {
				f.mockOrderRepository.OrderInfoMock.Expect(ctx, orderID).Return(&model.Order{}, nil)
			},
			expectedError: nil,
		},
		{
			name:    "Invalid Order ID",
			orderID: 0,
			mockSetup: func(f *fields, orderID int64) {
			},
			expectedError: ErrOrderInvalid,
		},
		{
			name:    "Order not found",
			orderID: 1,
			mockSetup: func(f *fields, orderID int64) {
				f.mockOrderRepository.OrderInfoMock.Expect(ctx, orderID).Return(&model.Order{}, ErrOrderNotFound)
			},
			expectedError: ErrOrderNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mc := minimock.NewController(t)
			mockOrderRepository := order_usecase.NewOrderRepositoryMock(mc)
			mockStockRepository := order_usecase.NewStockRepositoryMock(mc)
			mockEventsManagerRepository := order_usecase.NewEventManagersMock(mc)

			f := &fields{
				mockOrderRepository,
				mockStockRepository,
			}
			tc.mockSetup(f, tc.orderID)
			service := NewService(mockOrderRepository, mockStockRepository, mockEventsManagerRepository)
			_, err := service.OrderInfo(ctx, tc.orderID)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
