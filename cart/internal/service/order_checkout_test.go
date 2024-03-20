package service

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	product2 "route256.ozon.ru/project/cart/internal/clients/http/product"
	"route256.ozon.ru/project/cart/internal/domain"
	"route256.ozon.ru/project/cart/internal/service/mock"
	"testing"
)

func TestCartService_OrderCheckout(t *testing.T) {
	ctx := context.Background()
	type fields struct {
		repository     *mock.RepositoryMock
		productService *mock.ProductServiceMock
		lomsService    *mock.LomsServiceMock
	}
	type args struct {
		userId int64
		skuId  int64
	}
	tests := []struct {
		name    string
		prepare func(f *fields, args args)
		args    args
		wantErr error
	}{
		{
			name: "Cant get nil cart",
			prepare: func(f *fields, args args) {
				f.repository.GetItemsByUserIdMock.Expect(args.userId).Return(nil, ErrCartCantGet)
			},
			args: args{
				userId: 31337,
			},
			wantErr: ErrCartCantGet,
		},

		{
			name: "Get empty cart",
			prepare: func(f *fields, args args) {
				f.repository.GetItemsByUserIdMock.Expect(args.userId).Return(domain.ItemsMap{}, nil)
			},
			args: args{
				userId: 31337,
			},
			wantErr: ErrCartIsEmpty,
		},
		{
			name: "Successful checkout",
			prepare: func(f *fields, args args) {
				f.repository.GetItemsByUserIdMock.Expect(args.userId).Return(domain.ItemsMap{
					31337: {
						Sku_id: 31337,
						Count:  1,
					},
				}, nil)
				f.productService.GetProductMock.Expect(31337).Return(&product2.ProductGetProductResponse{}, nil)
				items := []domain.Item{
					{
						Sku_id: 31337,
						Count:  1,
					},
				}
				f.lomsService.CreateOrderMock.Expect(ctx, args.userId, items).Return(1, nil)
				f.repository.DeleteItemsByUserIdMock.Expect(args.userId).Return(nil)
			},
			args: args{
				userId: 31337,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mc := minimock.NewController(t)

			f := fields{
				productService: mock.NewProductServiceMock(mc),
				lomsService:    mock.NewLomsServiceMock(mc),
				repository:     mock.NewRepositoryMock(mc),
			}

			cartService := NewCartService(f.repository, f.productService, f.lomsService)
			tt.prepare(&f, tt.args)
			_, err := cartService.OrderCheckout(ctx, tt.args.userId)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
