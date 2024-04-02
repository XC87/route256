package service

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
	"route256.ozon.ru/project/cart/internal/domain"
	"route256.ozon.ru/project/cart/internal/service/mock"
	"testing"
)

func TestCartService_GetItemsByUserId(t *testing.T) {
	type fields struct {
		productService *mock.ProductServiceMock
		lomsService    *mock.LomsServiceMock
		repository     *mock.RepositoryMock
	}
	type args struct {
		userId int64
	}
	ctx := context.Background()
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
				f.productService.GetProductListMock.Expect(ctx, []int64{want.Items[0].SkuId}).Return([]*domain.Product{
					{
						Sku:   want.Items[0].SkuId,
						Name:  want.Items[0].Name,
						Price: want.Items[0].Price,
					},
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

	defer goleak.VerifyNone(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)

			f := fields{
				productService: mock.NewProductServiceMock(mc),
				lomsService:    mock.NewLomsServiceMock(mc),
				repository:     mock.NewRepositoryMock(mc),
			}

			cartService := NewCartService(f.repository, f.productService, f.lomsService)
			tt.prepare(&f, tt.args, tt.want)
			res, err := cartService.GetItemsByUserId(ctx, tt.args.userId)
			assert.Equal(t, tt.want, res)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
