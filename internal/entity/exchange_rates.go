package entity

import (
	"github.com/shopspring/decimal"
)

type ExchangeRate struct {
	ID             int64           `db:"id"`
	BaseCurrency   Currency        `db:"base_currency"`
	TargetCurrency Currency        `db:"target_currency"`
	Rate           decimal.Decimal `db:"rate"`
}
