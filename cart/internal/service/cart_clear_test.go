package service

import (
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCartService_DeleteItemsByUserId(t *testing.T) {
	type fields struct {
		productService *ProductServiceMock
		repository     *RepositoryMock
	}
	type args struct {
		userId int64
	}
	tests := []struct {
		name    string
		prepare func(f *fields, args args)
		args    args
		wantErr error
	}{
		{
			name: "Successful deletion of the cart",
			prepare: func(f *fields, args args) {
				f.repository.DeleteItemsByUserIdMock.Expect(args.userId).Return(nil)
			},
			args: args{
				userId: 31337,
			},
			wantErr: nil,
		},
		{
			name: "Bad user",
			prepare: func(f *fields, args args) {
			},
			args: args{
				userId: 0,
			},
			wantErr: ErrUserInvalid,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mc := minimock.NewController(t)

			f := fields{
				productService: NewProductServiceMock(mc),
				repository:     NewRepositoryMock(mc),
			}

			cartService := NewCartService(f.repository, f.productService)
			tt.prepare(&f, tt.args)
			err := cartService.DeleteItemsByUserId(tt.args.userId)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
