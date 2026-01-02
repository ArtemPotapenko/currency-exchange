package repository

import (
	"context"
	"currency-exchange/internal/entity"

	"github.com/shopspring/decimal"
)

type ExchangeRepository interface {
	Create(ctx context.Context, rate entity.ExchangeRate) (int64, error)
	Update(ctx context.Context, rate entity.ExchangeRate) error
	GetByID(ctx context.Context, id int64) (entity.ExchangeRate, error)
	GetRate(ctx context.Context, baseId int64, targetId int64) (decimal.Decimal, error)
}
