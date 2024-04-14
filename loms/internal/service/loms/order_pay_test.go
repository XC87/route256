package order_usecase

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"route256.ozon.ru/project/loms/internal/model"
	order_usecase "route256.ozon.ru/project/loms/internal/service/loms/mock"
	"testing"
)

func TestService_OrderPay(t *testing.T) {
	ctx := context.Background()
	type fields struct {
		mockOrderRepository *order_usecase.OrderRepositoryMock
		mockStockRepository *order_usecase.StockRepositoryMock
		mockEventManager    *order_usecase.EventManagersMock
	}
	testCases := []struct {
		name          string
		orderID       int64
		mockSetup     func(f *fields, orderID int64)
		expectedError error
	}{
		{
			name:    "Valid Order ID && status",
			orderID: 1,
			mockSetup: func(f *fields, orderID int64) {
				f.mockOrderRepository.OrderInfoMock.Expect(ctx, orderID).Return(&model.Order{Status: model.AwaitingPayment}, nil)
				f.mockOrderRepository.OrderPayMock.Expect(ctx, orderID).Return(nil)
				f.mockEventManager.PublishMock.Return(nil)
			},
			expectedError: nil,
		},
		{
			name:    "Valid order ID && invalid status",
			orderID: 1,
			mockSetup: func(f *fields, orderID int64) {
				f.mockOrderRepository.OrderInfoMock.Expect(ctx, orderID).Return(&model.Order{Status: model.New}, nil)
			},
			expectedError: ErrOrderCantBePaid,
		},
		{
			name:    "Valid order ID && already paid status",
			orderID: 1,
			mockSetup: func(f *fields, orderID int64) {
				f.mockOrderRepository.OrderInfoMock.Expect(ctx, orderID).Return(&model.Order{Status: model.Paid}, nil)
			},
			expectedError: ErrOrderAlreadyPaid,
		},
		{
			name:    "Invalid order ID",
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
				mockEventsManagerRepository,
			}
			tc.mockSetup(f, tc.orderID)
			service := NewService(mockOrderRepository, mockStockRepository, mockEventsManagerRepository)
			err := service.OrderPay(ctx, tc.orderID)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
