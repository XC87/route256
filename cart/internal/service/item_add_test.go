package service

import (
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"route256.ozon.ru/project/cart/internal/clients/product"
	"route256.ozon.ru/project/cart/internal/domain"
	"testing"
)

func TestCartService_AddItem(t *testing.T) {
	type fields struct {
		productService *ProductServiceMock
		repository     *RepositoryMock
	}
	type args struct {
		userId int64
		sku    domain.Item
	}
	tests := []struct {
		name    string
		prepare func(f *fields, args args)
		args    args
		wantErr error
	}{
		{
			name: "Check success cart add",
			prepare: func(f *fields, args args) {
				f.productService.GetProductMock.Expect(args.sku.Sku_id).Return(&product.ProductGetProductResponse{}, nil)
				f.repository.AddItemMock.Expect(args.userId, args.sku).Return(nil)
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
				f.productService.GetProductMock.Expect(args.sku.Sku_id).Return(nil, product.ErrProductNotFound)
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)

			f := fields{
				productService: NewProductServiceMock(mc),
				repository:     NewRepositoryMock(mc),
			}

			cartService := NewCartService(f.repository, f.productService)
			tt.prepare(&f, tt.args)

			err := cartService.AddItem(tt.args.userId, tt.args.sku)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
