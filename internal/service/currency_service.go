package service

import (
	"context"
	"currency-exchange/internal/dto"
	"currency-exchange/internal/entity"
	"currency-exchange/internal/errs"
	"currency-exchange/internal/pagination"
	repository "currency-exchange/internal/repository/db"
	"errors"
	"fmt"
	"log"
)

type CurrencyService struct {
	currencyRepository repository.CurrencyRepository
}

func NewCurrencyService(currencyRepository repository.CurrencyRepository) *CurrencyService {
	return &CurrencyService{currencyRepository: currencyRepository}
}

func (c *CurrencyService) CreateCurrency(ctx context.Context, request dto.CreateCurrencyRequest) (dto.CurrencyDto, error) {
	log.Printf("currency_service.create_currency start code=%s", request.Code)
	if err := request.Validate(); err != nil {
		return dto.CurrencyDto{}, err
	}
	currency := entity.Currency{
		Code:     request.Code,
		FullName: request.FullName,
		Sign:     request.Sign,
	}
	id, err := c.currencyRepository.Create(ctx, currency)
	if err != nil {
		log.Printf("currency_service.create_currency error: %v", err)
		return dto.CurrencyDto{}, fmt.Errorf("%w: create currency", errs.ErrInternal)
	}
	currency.ID = id
	log.Printf("currency_service.create_currency ok id=%d", id)
	return mapCurrency(currency), nil
}

func (c *CurrencyService) GetCurrencyByCode(ctx context.Context, code string) (dto.CurrencyDto, error) {
	log.Printf("currency_service.get_currency_by_code start code=%s", code)
	if code == "" {
		log.Printf("currency_service.get_currency_by_code validation_error: empty code")
		return dto.CurrencyDto{}, fmt.Errorf("%w: currency code is required", errs.ErrValidation)
	}

	currency, err := c.currencyRepository.GetByCode(ctx, code)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			log.Printf("currency_service.get_currency_by_code not_found code=%s", code)
			return dto.CurrencyDto{}, fmt.Errorf("%w: currency not found", errs.ErrNotFound)
		}
		log.Printf("currency_service.get_currency_by_code error: %v", err)
		return dto.CurrencyDto{}, fmt.Errorf("%w: get currency by code", errs.ErrInternal)
	}

	log.Printf("currency_service.get_currency_by_code ok id=%d", currency.ID)
	return mapCurrency(currency), nil
}

func (c *CurrencyService) GetAllCurrencyPage(ctx context.Context, request pagination.PageRequest) (pagination.Page[dto.CurrencyDto], error) {
	log.Printf("currency_service.get_all_currency_page start page=%d size=%d", request.PageNumber, request.PageSize)
	if request.PageNumber < 1 || request.PageSize < 1 {
		log.Printf("currency_service.get_all_currency_page validation_error page=%d size=%d", request.PageNumber, request.PageSize)
		return pagination.Page[dto.CurrencyDto]{}, fmt.Errorf(
			"%w: pageNumber and pageSize must be greater than zero",
			errs.ErrValidation,
		)
	}

	page, err := c.currencyRepository.GetPage(ctx, request)
	if err != nil {
		log.Printf("currency_service.get_all_currency_page error: %v", err)
		return pagination.Page[dto.CurrencyDto]{}, fmt.Errorf("%w: get currency page", errs.ErrInternal)
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
