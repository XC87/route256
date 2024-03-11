package repository

import (
	"github.com/stretchr/testify/require"
	"route256.ozon.ru/project/cart/internal/domain"
	"testing"
)

func BenchmarkInMemoryStorage(b *testing.B) {
	storage := NewMemoryRepository()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		storage.AddItem(123, domain.Item{
			Sku_id: 773297411,
			Count:  0,
		})
	}
}

func TestMemory_AddItem(t *testing.T) {
	type args struct {
		userId int64
		sku    domain.Item
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "Check success add",
			args: args{
				userId: 31337,
				sku: domain.Item{
					Sku_id: 773297411,
					Count:  3,
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memory := NewMemoryRepository()
			err := memory.AddItem(tt.args.userId, tt.args.sku)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestMemory_DeleteItem(t *testing.T) {
	type args struct {
		userId int64
		sku    domain.Item
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "Check success delete",
			args: args{
				userId: 31337,
				sku: domain.Item{
					Sku_id: 773297411,
					Count:  3,
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memory := NewMemoryRepository()
			err := memory.DeleteItem(tt.args.userId, tt.args.sku.Sku_id)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestMemory_DeleteItemsByUserId(t *testing.T) {
	type args struct {
		userId int64
		sku    domain.Item
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "Check success delete user cart",
			args: args{
				userId: 31337,
				sku: domain.Item{
					Sku_id: 773297411,
					Count:  3,
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memory := NewMemoryRepository()
			err := memory.DeleteItemsByUserId(tt.args.userId)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestMemory_GetItemsByUserId(t *testing.T) {
	type args struct {
		userId int64
		sku    domain.Item
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "Check get user cart",
			args: args{
				userId: 31337,
				sku: domain.Item{
					Sku_id: 773297411,
					Count:  3,
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memory := NewMemoryRepository()
			_, err := memory.GetItemsByUserId(tt.args.userId)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
