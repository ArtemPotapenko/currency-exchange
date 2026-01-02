package repository

import (
	"context"
	"currency-exchange/internal/entity"
	"currency-exchange/internal/pagination"
)

type CurrencyRepository interface {
	Create(ctx context.Context, currency entity.Currency) (int64, error)
	GetByCode(ctx context.Context, code string) (entity.Currency, error)
	GetPage(ctx context.Context, page pagination.PageRequest) (pagination.Page[entity.Currency], error)
}
