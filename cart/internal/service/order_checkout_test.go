package service

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
	"math/rand/v2"
	"route256.ozon.ru/project/cart/internal/domain"
	"route256.ozon.ru/project/cart/internal/service/mock"
	"testing"
)

func TestCartService_OrderCheckout(t *testing.T) {
	type fields struct {
		repository     *mock.RepositoryMock
		productService *mock.ProductServiceMock
		lomsService    *mock.LomsServiceMock
	}
	type args struct {
		userId int64
		skuId  int64
	}
	ctx := context.Background()
	tests := []struct {
		name    string
		prepare func(f *fields, args args)
		args    args
		wantErr error
	}{
		{
			name: "Cant get nil cart",
			prepare: func(f *fields, args args) {
				f.repository.GetItemsByUserIdMock.Return(nil, ErrCartCantGet)
			},
			args: args{
				userId: rand.Int64(),
			},
			wantErr: ErrCartCantGet,
		},

		{
			name: "Get empty cart",
			prepare: func(f *fields, args args) {
				f.repository.GetItemsByUserIdMock.Return(domain.ItemsMap{}, nil)
			},
			args: args{
				userId: rand.Int64(),
			},
			wantErr: ErrCartIsEmpty,
		},
		{
			name: "Successful checkout",
			prepare: func(f *fields, args args) {
				f.repository.GetItemsByUserIdMock.Return(domain.ItemsMap{
					31337: {
						Sku_id: 31337,
						Count:  1,
					},
				}, nil)
				Product := []*domain.Product{
					{
						Sku: 31337,
					},
				}
				f.productService.GetProductListMock.Return(Product, nil)
				f.lomsService.CreateOrderMock.Return(1, nil)
				f.repository.DeleteItemsByUserIdMock.Return(nil)
			},
			args: args{
				userId: rand.Int64(),
			},
			wantErr: nil,
		},
	}
	defer goleak.VerifyNone(t)
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
