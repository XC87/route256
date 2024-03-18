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
				userId: 1,
				sku: domain.Item{
					Sku_id: 773297411,
					Count:  3,
				},
			},
			wantErr: nil,
		},
		{
			name: "Check success add",
			args: args{
				userId: 2,
				sku: domain.Item{
					Sku_id: 773297411,
					Count:  3,
				},
			},
			wantErr: nil,
		},
	}
	memory := NewMemoryRepository()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			userId := tt.args.userId
			item := tt.args.sku
			err := memory.AddItem(userId, item)

			require.ErrorIs(t, err, tt.wantErr)
			if memory.cart[userId][item.Sku_id].Count != item.Count {
				t.Errorf("Expected count to be %d, but got %d", item.Count, memory.cart[userId][item.Sku_id].Count)
			}
		})
	}
}

func TestMemory_DeleteItem(t *testing.T) {
	type fields struct {
		memoryRepository *Memory
	}
	type args struct {
		userId int64
		sku    domain.Item
	}
	tests := []struct {
		name    string
		args    args
		prepare func(f *fields, args args)
		wantErr error
	}{
		{
			name: "Check success delete",
			args: args{
				userId: 1,
				sku: domain.Item{
					Sku_id: 773297411,
					Count:  3,
				},
			},
			prepare: func(f *fields, args args) {
				f.memoryRepository.AddItem(args.userId, args.sku)
			},
			wantErr: nil,
		},
		{
			name: "Check success delete",
			args: args{
				userId: 2,
				sku: domain.Item{
					Sku_id: 773297411,
					Count:  3,
				},
			},
			prepare: func(f *fields, args args) {
				f.memoryRepository.AddItem(args.userId, args.sku)
			},
			wantErr: nil,
		},
	}
	memory := NewMemoryRepository()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := fields{
				memoryRepository: memory,
			}
			tt.prepare(&f, tt.args)

			userId := tt.args.userId
			item := tt.args.sku

			err := memory.DeleteItem(tt.args.userId, tt.args.sku.Sku_id)

			require.ErrorIs(t, err, tt.wantErr)
			if _, ok := memory.cart[userId][item.Sku_id]; ok {
				t.Errorf("Expected item to be deleted")
			}
		})
	}
}

func TestMemory_DeleteItemsByUserId(t *testing.T) {
	type fields struct {
		memoryRepository *Memory
	}
	type args struct {
		userId int64
		sku    domain.Item
	}
	tests := []struct {
		name    string
		args    args
		prepare func(f *fields, args args)
		wantErr error
	}{
		{
			name: "Check success delete user cart",
			args: args{
				userId: 1,
				sku: domain.Item{
					Sku_id: 773297411,
					Count:  3,
				},
			},
			prepare: func(f *fields, args args) {
				f.memoryRepository.AddItem(args.userId, args.sku)
			},
			wantErr: nil,
		},
		{
			name: "Check success delete user cart",
			args: args{
				userId: 2,
				sku: domain.Item{
					Sku_id: 773297411,
					Count:  3,
				},
			},
			prepare: func(f *fields, args args) {
				f.memoryRepository.AddItem(args.userId, args.sku)
			},
			wantErr: nil,
		},
	}
	memory := NewMemoryRepository()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := fields{
				memoryRepository: memory,
			}
			tt.prepare(&f, tt.args)

			err := memory.DeleteItemsByUserId(tt.args.userId)
			require.ErrorIs(t, err, tt.wantErr)

			if _, ok := memory.cart[tt.args.userId]; ok {
				t.Errorf("Expected user cart to be deleted")
			}
		})
	}
}

func TestMemory_GetItemsByUserId(t *testing.T) {
	type fields struct {
		memoryRepository *Memory
	}
	type args struct {
		userId int64
		sku    domain.Item
	}
	tests := []struct {
		name    string
		args    args
		prepare func(f *fields, args args)
		wantErr error
	}{
		{
			name: "Check get user cart",
			args: args{
				userId: 1,
				sku: domain.Item{
					Sku_id: 773297411,
					Count:  3,
				},
			},
			prepare: func(f *fields, args args) {
				f.memoryRepository.AddItem(args.userId, args.sku)
			},
			wantErr: nil,
		},
		{
			name: "Check get user cart",
			args: args{
				userId: 2,
				sku: domain.Item{
					Sku_id: 773297411,
					Count:  3,
				},
			},
			prepare: func(f *fields, args args) {
				f.memoryRepository.AddItem(args.userId, args.sku)
			},
			wantErr: nil,
		},
	}

	memory := NewMemoryRepository()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := fields{
				memoryRepository: memory,
			}
			tt.prepare(&f, tt.args)

			itemsMap, err := memory.GetItemsByUserId(tt.args.userId)
			require.Equal(t, domain.ItemsMap{tt.args.sku.Sku_id: tt.args.sku}, itemsMap)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
