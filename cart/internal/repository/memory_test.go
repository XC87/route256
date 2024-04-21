package repository

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
	"route256.ozon.ru/project/cart/internal/domain"
	"testing"
)

func BenchmarkInMemoryStorage(b *testing.B) {
	storage := NewMemoryRepository()
	b.ResetTimer()
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		storage.AddItem(ctx, 123, domain.Item{
			Sku_id: 773297411,
			Count:  0,
		})
	}
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestMemory_AddItem(t *testing.T) {
	type args struct {
		userId int64
		sku    domain.Item
	}
	ctx := context.Background()
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
			err := memory.AddItem(ctx, userId, item)

			require.ErrorIs(t, err, tt.wantErr)

			itemsMap, _ := memory.GetItemsByUserId(ctx, userId)
			if itemsMap[item.Sku_id].Count != item.Count {
				t.Errorf("Expected count to be %d, but got %d", item.Count, itemsMap[item.Sku_id].Count)
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
	ctx := context.Background()
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
				f.memoryRepository.AddItem(ctx, args.userId, args.sku)
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
				f.memoryRepository.AddItem(ctx, args.userId, args.sku)
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

			err := memory.DeleteItem(ctx, tt.args.userId, tt.args.sku.Sku_id)
			require.ErrorIs(t, err, tt.wantErr)

			itemsMap, _ := memory.GetItemsByUserId(ctx, userId)
			if _, ok := itemsMap[item.Sku_id]; ok {
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
	ctx := context.Background()
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
				f.memoryRepository.AddItem(ctx, args.userId, args.sku)
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
				f.memoryRepository.AddItem(ctx, args.userId, args.sku)
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

			err := memory.DeleteItemsByUserId(ctx, tt.args.userId)
			require.ErrorIs(t, err, tt.wantErr)

			itemsMap, _ := memory.GetItemsByUserId(ctx, tt.args.userId)
			assert.Empty(t, itemsMap)
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
	ctx := context.Background()
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
				f.memoryRepository.AddItem(ctx, args.userId, args.sku)
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
				f.memoryRepository.AddItem(ctx, args.userId, args.sku)
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

			itemsMap, err := memory.GetItemsByUserId(ctx, tt.args.userId)
			require.Equal(t, domain.ItemsMap{tt.args.sku.Sku_id: tt.args.sku}, itemsMap)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
