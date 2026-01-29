package dto

import (
	"fmt"
	"log"

	"currency-exchange/internal/entity"
	"currency-exchange/internal/errs"

	"github.com/shopspring/decimal"
)

type ExchangeRateDto struct {
	ID             int64           `json:"id"`
	BaseCurrency   CurrencyDto     `json:"baseCurrency"`
	TargetCurrency CurrencyDto     `json:"targetCurrency"`
	Rate           decimal.Decimal `json:"rate"`
}

type ExchangeDto struct {
	ExchangeRate  ExchangeRateDto `json:"exchangeRate"`
	Amount        decimal.Decimal `json:"amount"`
	ConvertAmount decimal.Decimal `json:"convertAmount"`
}

type CreateRateRequest struct {
	BaseCode   string          `json:"baseCode"`
	TargetCode string          `json:"targetCode"`
	Rate       decimal.Decimal `json:"rate"`
}

type UpdateRateRequest struct {
	BaseCode   string          `json:"baseCode"`
	TargetCode string          `json:"targetCode"`
	Rate       decimal.Decimal `json:"rate"`
}

type ExchangeRequest struct {
	BaseCode   string          `json:"baseCode"`
	TargetCode string          `json:"targetCode"`
	Amount     decimal.Decimal `json:"amount"`
}

func (r CreateRateRequest) Validate() error {
	if err := validateRatePrecision(r.Rate); err != nil {
		log.Printf("create_rate_request validation_error: %v", err)
		return err
	}
	return nil
}

func (r UpdateRateRequest) Validate() error {
	if err := validateRatePrecision(r.Rate); err != nil {
		log.Printf("update_rate_request validation_error: %v", err)
		return err
	}
	return nil
}

func (r ExchangeRequest) Validate() error {
	if r.BaseCode == "" || r.TargetCode == "" {
		log.Printf("exchange_request validation_error: empty code base=%s target=%s", r.BaseCode, r.TargetCode)
		return fmt.Errorf("%w: currency codes are required", errs.ErrValidation)
	}
	if r.Amount.LessThanOrEqual(decimal.Zero) {
		log.Printf("exchange_request validation_error: non_positive amount=%s", r.Amount.String())
		return fmt.Errorf("%w: amount must be greater than zero", errs.ErrValidation)
	}
	if err := validateAmountPrecision(r.Amount); err != nil {
		log.Printf("exchange_request validation_error: %v", err)
		return err
	}
	return nil
}

func validateRatePrecision(rate decimal.Decimal) error {
	if rate.Exponent() < -entity.ExchangeRateMaxScale {
		return fmt.Errorf("%w: invalid exchange rate precision", errs.ErrValidation)
	}
	return nil
}

func validateAmountPrecision(amount decimal.Decimal) error {
	if amount.Exponent() < -entity.ExchangeAmountMaxScale {
		return fmt.Errorf("%w: invalid amount precision", errs.ErrValidation)
	}
	return nil
}
