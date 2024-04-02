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
				f.mockStockRepository.GetCountBySkuMock.Expect(ctx, order.Items[0].SKU).Return(1, nil)
				f.mockStockRepository.GetBySkuMock.Expect(ctx, order.Items[0].SKU).Return(&model.ProductStock{SKU: order.Items[0].SKU, TotalCount: 5, Reserved: 0}, nil)
				f.mockStockRepository.ReserveMock.Expect(ctx, order.Items).Return(nil)

				f.mockOrderRepository.OrderCreateMock.Expect(ctx, order).Return(1, nil)
				f.mockOrderRepository.OrderUpdateMock.Expect(ctx, order).Return(nil)

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
				f.mockStockRepository.GetCountBySkuMock.Expect(ctx, order.Items[0].SKU).Return(1, nil)
				f.mockStockRepository.GetBySkuMock.Expect(ctx, order.Items[0].SKU).Return(&model.ProductStock{SKU: order.Items[0].SKU, TotalCount: 5, Reserved: 0}, nil)

				f.mockOrderRepository.OrderCreateMock.Expect(ctx, order).Return(1, nil)
				f.mockOrderRepository.OrderUpdateMock.Expect(ctx, order).Return(nil)

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
			f := &fields{
				mockOrderRepository,
				mockStockRepository,
			}
			tc.mockSetup(f, tc.order)
			service := NewService(mockOrderRepository, mockStockRepository)
			_, err := service.OrderCreate(ctx, tc.order)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
