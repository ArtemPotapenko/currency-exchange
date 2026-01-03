package service

import (
	"context"
	"currency-exchange/internal/dto"
	"currency-exchange/internal/entity"
	apperror "currency-exchange/internal/error"
	"currency-exchange/internal/pagination"
	"currency-exchange/internal/repository"
	"errors"
	"fmt"
	"log"
	"unicode/utf8"
)

type CurrencyService struct {
	ctx                context.Context
	currencyRepository repository.CurrencyRepository
}

func NewCurrencyService(ctx context.Context, currencyRepository repository.CurrencyRepository) *CurrencyService {
	return &CurrencyService{ctx: ctx, currencyRepository: currencyRepository}
}

func (c *CurrencyService) CreateCurrency(code string, fullName string, sign string) (dto.CurrencyDto, error) {
	log.Printf("currency_service.create_currency start code=%s", code)
	if len(code) == 0 || utf8.RuneCountInString(code) > entity.CurrencyCodeMaxLen {
		log.Printf("currency_service.create_currency validation_error code=%s", code)
		return dto.CurrencyDto{}, apperror.Validation(
			"invalid currency code",
			"code must be 1.."+fmt.Sprint(entity.CurrencyCodeMaxLen)+" symbols",
		)
	}
	if len(sign) == 0 || utf8.RuneCountInString(sign) > entity.CurrencySignMaxLen {
		log.Printf("currency_service.create_currency validation_error sign=%s", sign)
		return dto.CurrencyDto{}, apperror.Validation(
			"invalid currency sign",
			"sign must be 1.."+fmt.Sprint(entity.CurrencySignMaxLen)+" symbols",
		)
	}
	if utf8.RuneCountInString(fullName) < entity.CurrencyFullNameMinLen ||
		utf8.RuneCountInString(fullName) > entity.CurrencyFullNameMaxLen {
		log.Printf("currency_service.create_currency validation_error full_name=%s", fullName)
		return dto.CurrencyDto{}, apperror.Validation(
			"invalid currency full name",
			"full name length must be "+
				fmt.Sprint(entity.CurrencyFullNameMinLen)+".."+
				fmt.Sprint(entity.CurrencyFullNameMaxLen)+" symbols",
		)
	}

	currency := entity.Currency{
		Code:     code,
		FullName: fullName,
		Sign:     sign,
	}
	id, err := c.currencyRepository.Create(c.ctx, currency)
	if err != nil {
		log.Printf("currency_service.create_currency error: %v", err)
		return dto.CurrencyDto{}, apperror.Internal("create currency", err.Error())
	}
	currency.ID = id
	log.Printf("currency_service.create_currency ok id=%d", id)
	return mapCurrency(currency), nil
}

func (c *CurrencyService) GetCurrencyByCode(code string) (dto.CurrencyDto, error) {
	log.Printf("currency_service.get_currency_by_code start code=%s", code)
	if code == "" {
		log.Printf("currency_service.get_currency_by_code validation_error: empty code")
		return dto.CurrencyDto{}, apperror.Validation("currency code is required", "empty code")
	}

	currency, err := c.currencyRepository.GetByCode(c.ctx, code)
	if err != nil {
		var notFoundErr *apperror.NotFoundError
		if errors.As(err, &notFoundErr) {
			log.Printf("currency_service.get_currency_by_code not_found code=%s", code)
			return dto.CurrencyDto{}, apperror.NotFound("currency not found", "code="+code)
		}
		log.Printf("currency_service.get_currency_by_code error: %v", err)
		return dto.CurrencyDto{}, apperror.Internal("get currency by code", err.Error())
	}

	log.Printf("currency_service.get_currency_by_code ok id=%d", currency.ID)
	return mapCurrency(currency), nil
}

func (c *CurrencyService) GetAllCurrencyPage(request pagination.PageRequest) (pagination.Page[dto.CurrencyDto], error) {
	log.Printf("currency_service.get_all_currency_page start page=%d size=%d", request.PageNumber, request.PageSize)
	if request.PageNumber < 1 || request.PageSize < 1 {
		log.Printf("currency_service.get_all_currency_page validation_error page=%d size=%d", request.PageNumber, request.PageSize)
		return pagination.Page[dto.CurrencyDto]{}, apperror.Validation(
			"pageNumber and pageSize must be greater than zero",
			"pageNumber or pageSize less than 1",
		)
	}

	page, err := c.currencyRepository.GetPage(c.ctx, request)
	if err != nil {
		log.Printf("currency_service.get_all_currency_page error: %v", err)
		return pagination.Page[dto.CurrencyDto]{}, apperror.Internal("get currency page", err.Error())
	}

	items := make([]dto.CurrencyDto, 0, len(page.Items))
	for _, currency := range page.Items {
		items = append(items, dto.CurrencyDto{
			ID:       currency.ID,
			Code:     currency.Code,
			FullName: currency.FullName,
			Sign:     currency.Sign,
		})
	}

	log.Printf("currency_service.get_all_currency_page ok total=%d", page.Total)
	return pagination.Page[dto.CurrencyDto]{
		Items:      items,
		PageNumber: page.PageNumber,
		PageSize:   page.PageSize,
		Total:      page.Total,
	}, nil
}

func mapCurrency(currency entity.Currency) dto.CurrencyDto {
	return dto.CurrencyDto{
		ID:       currency.ID,
		Code:     currency.Code,
		FullName: currency.FullName,
		Sign:     currency.Sign,
	}
}
