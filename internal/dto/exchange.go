package dto

import "github.com/shopspring/decimal"

type CurrencyDto struct {
	ID       int64  `json:"id"`
	Code     string `json:"code"`
	FullName string `json:"fullName"`
	Sign     string `json:"sign"`
}

type ExchangeRateDto struct {
	ID             int64           `json:"id"`
	BaseCurrency   CurrencyDto     `json:"baseCurrency"`
	TargetCurrency CurrencyDto     `json:"targetCurrency"`
	Rate           decimal.Decimal `json:"rate" swaggertype:"string"`
}

type ExchangeDto struct {
	ExchangeRate  ExchangeRateDto `json:"exchangeRate"`
	Amount        decimal.Decimal `json:"amount" swaggertype:"string"`
	ConvertAmount decimal.Decimal `json:"convertAmount" swaggertype:"string"`
}

type CreateCurrencyRequest struct {
	Code     string `json:"code"`
	FullName string `json:"fullName"`
	Sign     string `json:"sign"`
}

type CreateRateRequest struct {
	BaseCode   string          `json:"baseCode"`
	TargetCode string          `json:"targetCode"`
	Rate       decimal.Decimal `json:"rate" swaggertype:"string"`
}

type UpdateRateRequest struct {
	BaseCode   string          `json:"baseCode"`
	TargetCode string          `json:"targetCode"`
	Rate       decimal.Decimal `json:"rate" swaggertype:"string"`
}
