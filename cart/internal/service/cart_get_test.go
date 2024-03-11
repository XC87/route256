package service

import (
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"route256.ozon.ru/project/cart/internal/clients/product"
	"route256.ozon.ru/project/cart/internal/domain"
	"testing"
)

func TestCartService_GetItemsByUserId(t *testing.T) {
	type fields struct {
		productService *ProductServiceMock
		repository     *RepositoryMock
	}
	type args struct {
		userId int64
	}
	tests := []struct {
		name    string
		prepare func(f *fields, args args, want *CartResponse)
		args    args
		want    *CartResponse
		wantErr error
	}{
		{
			name: "Successful get items from cart",
			prepare: func(f *fields, args args, want *CartResponse) {
				maps := domain.ItemsMap{
					want.Items[0].SkuId: domain.Item{
						Sku_id: want.Items[0].SkuId,
						Count:  want.Items[0].Count,
					},
				}
				f.productService.GetProductMock.Expect(want.Items[0].SkuId).Return(&product.ProductGetProductResponse{
					Name:  want.Items[0].Name,
					Price: want.Items[0].Price,
				}, nil)
				f.repository.GetItemsByUserIdMock.Expect(args.userId).Return(maps, nil)
			},
			args: args{
				userId: 31337,
			},
			want: &CartResponse{
				Items: []CartItem{
					{
						SkuId: 773297411,
						Name:  "Кроссовки NIKE JORDAN",
						Count: 10,
						Price: 2202,
					},
				},
				TotalPrice: 22020,
			},
			wantErr: nil,
		},
		{
			name: "Successful get empty cart",
			prepare: func(f *fields, args args, want *CartResponse) {
				f.repository.GetItemsByUserIdMock.Expect(args.userId).Return(domain.ItemsMap{}, nil)
			},
			args: args{
				userId: 31337,
			},
			want:    &CartResponse{},
			wantErr: nil,
		},
		{
			name: "Bad user",
			prepare: func(f *fields, args args, want *CartResponse) {
			},
			args: args{
				userId: 0,
			},
			want:    nil,
			wantErr: ErrUserInvalid,
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
			tt.prepare(&f, tt.args, tt.want)
			res, err := cartService.GetItemsByUserId(tt.args.userId)
			assert.Equal(t, tt.want, res)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
