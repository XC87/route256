package service

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
	product2 "route256.ozon.ru/project/cart/internal/clients/http/product"
	"route256.ozon.ru/project/cart/internal/domain"
	"route256.ozon.ru/project/cart/internal/service/mock"
	"testing"
)

func TestCartService_AddItem(t *testing.T) {
	type fields struct {
		productService *mock.ProductServiceMock
		lomsService    *mock.LomsServiceMock
		repository     *mock.RepositoryMock
	}
	type args struct {
		userId int64
		sku    domain.Item
	}

	ctx := context.Background()
	tests := []struct {
		name    string
		prepare func(f *fields, args args)
		args    args
		wantErr error
	}{
		{
			name: "Check success cart add",
			prepare: func(f *fields, args args) {
				f.lomsService.GetStockInfoMock.Return(5, nil)
				f.productService.GetProductMock.Return(&domain.Product{}, nil)
				f.repository.AddItemMock.Return(nil)
			},
			args: args{
				userId: 31337,
				sku: domain.Item{
					Sku_id: 773297411,
					Count:  3,
				},
			},
			wantErr: nil,
		},
		{
			name:    "Check count error",
			prepare: func(f *fields, args args) {},
			args: args{
				userId: 31337,
				sku: domain.Item{
					Sku_id: 773297411,
					Count:  0,
				},
			},
			wantErr: ErrProductCountInvalid,
		},

		{
			name: "Check not enough error",
			prepare: func(f *fields, args args) {
				f.lomsService.GetStockInfoMock.Return(5, nil)
				f.productService.GetProductMock.Return(&domain.Product{}, nil)
			},
			args: args{
				userId: 31337,
				sku: domain.Item{
					Sku_id: 773297411,
					Count:  1000,
				},
			},
			wantErr: ErrStockInsufficient,
		},

		{
			name:    "Adding an service to the cart with an invalid user ID",
			prepare: func(f *fields, args args) {},
			args: args{
				userId: 0,
				sku: domain.Item{
					Sku_id: 773297411,
					Count:  0,
				},
			},
			wantErr: ErrUserInvalid,
		},

		{
			name: "Adding an service to the cart with an invalid SKU ID",
			prepare: func(f *fields, args args) {
				f.productService.GetProductMock.Return(nil, product2.ErrProductNotFound)
			},
			args: args{
				userId: 31337,
				sku: domain.Item{
					Sku_id: 0,
					Count:  3,
				},
			},
			wantErr: ErrProductNotFound,
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

			err := cartService.AddItem(ctx, tt.args.userId, tt.args.sku)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
