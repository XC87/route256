package service

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
	"route256.ozon.ru/project/cart/internal/service/mock"
	"testing"
)

func TestCartService_DeleteItem(t *testing.T) {
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
			name: "Successful deletion of an product from the cart",
			prepare: func(f *fields, args args) {
				f.repository.DeleteItemMock.Expect(ctx, args.userId, args.skuId).Return(nil)
			},
			args: args{
				userId: 31337,
				skuId:  773297411,
			},
			wantErr: nil,
		},
		{
			name: "Bad user",
			prepare: func(f *fields, args args) {
			},
			args: args{
				userId: 0,
				skuId:  773297411,
			},
			wantErr: ErrUserInvalid,
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
			err := cartService.DeleteItem(ctx, tt.args.userId, tt.args.skuId)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
