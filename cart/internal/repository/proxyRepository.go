package repository

import (
	"context"
	"go.opentelemetry.io/otel"
	"route256.ozon.ru/project/cart/internal/domain"
	"route256.ozon.ru/project/cart/internal/service"
)

type ProxyRepository struct {
	Repository service.Repository
}

func NewProxyRepository(rep service.Repository) *ProxyRepository {
	return &ProxyRepository{
		Repository: rep,
	}
}

func (p *ProxyRepository) AddItem(ctx context.Context, userId int64, item domain.Item) error {
	_, span := otel.Tracer("default").Start(ctx, "repository.AddItem")
	defer span.End()

	return p.Repository.AddItem(ctx, userId, item)
}

func (p *ProxyRepository) DeleteItem(ctx context.Context, userId int64, skuId int64) error {
	_, span := otel.Tracer("default").Start(ctx, "repository.DeleteItem")
	defer span.End()

	return p.Repository.DeleteItem(ctx, userId, skuId)
}

func (p *ProxyRepository) DeleteItemsByUserId(ctx context.Context, userId int64) error {
	_, span := otel.Tracer("default").Start(ctx, "repository.DeleteItemsByUserId")
	defer span.End()

	return p.Repository.DeleteItemsByUserId(ctx, userId)
}

func (p *ProxyRepository) GetItemsByUserId(ctx context.Context, userId int64) (domain.ItemsMap, error) {
	_, span := otel.Tracer("default").Start(ctx, "repository.GetItemsByUserId")
	defer span.End()

	return p.Repository.GetItemsByUserId(ctx, userId)
}
