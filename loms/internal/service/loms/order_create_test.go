package order_usecase

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"route256.ozon.ru/project/loms/internal/model"
	order_usecase "route256.ozon.ru/project/loms/internal/service/loms/mock"
	"testing"
)

func TestService_OrderCreate(t *testing.T) {
	ctx := context.Background()
	type fields struct {
		mockOrderRepository *order_usecase.OrderRepositoryMock
		mockStockRepository *order_usecase.StockRepositoryMock
		mockEventManager    *order_usecase.EventManagersMock
	}
	testCases := []struct {
		name          string
		order         *model.Order
		mockSetup     func(f *fields, order *model.Order)
		expectedError error
	}{
		{
			name: "Successful order creation with reserve",
			order: &model.Order{
				Status: model.New,
				User:   1,
				Items:  []model.Item{{SKU: 2, Count: 5}},
			},
			mockSetup: func(f *fields, order *model.Order) {
				f.mockStockRepository.GetCountBySkuMock.Return(1, nil)
				f.mockStockRepository.GetBySkuMock.Return(&model.ProductStock{SKU: order.Items[0].SKU, TotalCount: 5, Reserved: 0}, nil)
				f.mockStockRepository.ReserveMock.Return(nil)

				f.mockOrderRepository.OrderCreateMock.Return(1, nil)
				f.mockOrderRepository.OrderUpdateMock.Return(nil)

				f.mockEventManager.TriggerMock.Return(nil)

			},
			expectedError: nil,
		},
		{
			name: "Successful order creation with cant reserve",
			order: &model.Order{
				Status: model.New,
				User:   1,
				Items:  []model.Item{{SKU: 2, Count: 55}},
			},
			mockSetup: func(f *fields, order *model.Order) {
				f.mockStockRepository.GetCountBySkuMock.Return(1, nil)
				f.mockStockRepository.GetBySkuMock.Expect(ctx, order.Items[0].SKU).Return(&model.ProductStock{SKU: order.Items[0].SKU, TotalCount: 5, Reserved: 0}, nil)

				f.mockOrderRepository.OrderCreateMock.Return(1, nil)
				f.mockOrderRepository.OrderUpdateMock.Return(nil)

				f.mockEventManager.TriggerMock.Return(nil)

			},
			expectedError: ErrOrderCantReserve,
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
			tc.mockSetup(f, tc.order)
			service := NewService(mockOrderRepository, mockStockRepository, mockEventsManagerRepository)
			_, err := service.OrderCreate(ctx, tc.order)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
