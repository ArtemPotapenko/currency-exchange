package dto

import (
	"fmt"
	"log"
	"unicode/utf8"

	"currency-exchange/internal/entity"
	"currency-exchange/internal/errs"
)

type CurrencyDto struct {
	ID       int64  `json:"id"`
	Code     string `json:"code"`
	FullName string `json:"fullName"`
	Sign     string `json:"sign"`
}

type CreateCurrencyRequest struct {
	Code     string `json:"code"`
	FullName string `json:"fullName"`
	Sign     string `json:"sign"`
}

func (r CreateCurrencyRequest) Validate() error {
	if len(r.Code) == 0 || utf8.RuneCountInString(r.Code) > entity.CurrencyCodeMaxLen {
		log.Printf("create_currency_request validation_error code=%s", r.Code)
		return fmt.Errorf("%w: invalid currency code", errs.ErrValidation)
	}
	if len(r.Sign) == 0 || utf8.RuneCountInString(r.Sign) > entity.CurrencySignMaxLen {
		log.Printf("create_currency_request validation_error sign=%s", r.Sign)
		return fmt.Errorf("%w: invalid currency sign", errs.ErrValidation)
	}
	if utf8.RuneCountInString(r.FullName) < entity.CurrencyFullNameMinLen ||
		utf8.RuneCountInString(r.FullName) > entity.CurrencyFullNameMaxLen {
		log.Printf("create_currency_request validation_error full_name=%s", r.FullName)
		return fmt.Errorf("%w: invalid currency full name", errs.ErrValidation)
	}
	return nil
}
