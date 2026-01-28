package service

import (
	"context"
	"currency-exchange/internal/dto"
	"currency-exchange/internal/entity"
	"currency-exchange/internal/errs"
	repository "currency-exchange/internal/repository/db"
	"errors"
	"fmt"
	"log"
)

type ExchangeService struct {
	exchangeRepository repository.ExchangeRepository
	currencyRepository repository.CurrencyRepository
}

func NewExchangeService(
	exchangeRepository repository.ExchangeRepository,
	currencyRepository repository.CurrencyRepository,
) *ExchangeService {
	return &ExchangeService{
		exchangeRepository: exchangeRepository,
		currencyRepository: currencyRepository,
	}
}

func (s *ExchangeService) CreateRate(ctx context.Context, request dto.CreateRateRequest) (dto.ExchangeRateDto, error) {
	log.Printf("exchange_service.create_rate start base=%s target=%s", request.BaseCode, request.TargetCode)
	if err := request.Validate(); err != nil {
		return dto.ExchangeRateDto{}, err
	}
	baseCurrency, err := s.currencyRepository.GetByCode(ctx, request.BaseCode)
	if err != nil {
		return dto.ExchangeRateDto{}, s.wrapCurrencyError("base currency", request.BaseCode, err)
	}
	targetCurrency, err := s.currencyRepository.GetByCode(ctx, request.TargetCode)
	if err != nil {
		return dto.ExchangeRateDto{}, s.wrapCurrencyError("target currency", request.TargetCode, err)
	}

	entityRate := entity.ExchangeRate{
		BaseCurrency:   baseCurrency,
		TargetCurrency: targetCurrency,
		Rate:           request.Rate,
	}
	id, err := s.exchangeRepository.Create(ctx, entityRate)
	if err != nil {
		log.Printf("exchange_service.create_rate error: %v", err)
		return dto.ExchangeRateDto{}, fmt.Errorf("%w: create exchange rate", errs.ErrInternal)
	}
	entityRate.ID = id
	log.Printf("exchange_service.create_rate ok id=%d", id)
	return mapRate(entityRate), nil
}

func (s *ExchangeService) UpdateRate(ctx context.Context, id int64, request dto.UpdateRateRequest) (dto.ExchangeRateDto, error) {
	log.Printf("exchange_service.update_rate start id=%d", id)
	if err := request.Validate(); err != nil {
		return dto.ExchangeRateDto{}, err
	}
	baseCurrency, err := s.currencyRepository.GetByCode(ctx, request.BaseCode)
	if err != nil {
		return dto.ExchangeRateDto{}, s.wrapCurrencyError("base currency", request.BaseCode, err)
	}
	targetCurrency, err := s.currencyRepository.GetByCode(ctx, request.TargetCode)
	if err != nil {
		return dto.ExchangeRateDto{}, s.wrapCurrencyError("target currency", request.TargetCode, err)
	}

	entityRate := entity.ExchangeRate{
		ID:             id,
		BaseCurrency:   baseCurrency,
		TargetCurrency: targetCurrency,
		Rate:           request.Rate,
	}
	if err := s.exchangeRepository.Update(ctx, entityRate); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			log.Printf("exchange_service.update_rate not_found id=%d", id)
			return dto.ExchangeRateDto{}, fmt.Errorf("%w: exchange rate not found", errs.ErrNotFound)
		}
		log.Printf("exchange_service.update_rate error: %v", err)
		return dto.ExchangeRateDto{}, fmt.Errorf("%w: update exchange rate", errs.ErrInternal)
	}

	log.Printf("exchange_service.update_rate ok id=%d", id)
	return mapRate(entityRate), nil
}

func (s *ExchangeService) GetRateByID(ctx context.Context, id int64) (dto.ExchangeRateDto, error) {
	log.Printf("exchange_service.get_rate_by_id start id=%d", id)
	rate, err := s.exchangeRepository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			log.Printf("exchange_service.get_rate_by_id not_found id=%d", id)
			return dto.ExchangeRateDto{}, fmt.Errorf("%w: exchange rate not found", errs.ErrNotFound)
		}
		log.Printf("exchange_service.get_rate_by_id error: %v", err)
		return dto.ExchangeRateDto{}, fmt.Errorf("%w: get exchange rate by id", errs.ErrInternal)
	}
	log.Printf("exchange_service.get_rate_by_id ok id=%d", rate.ID)
	return mapRate(rate), nil
}

func (s *ExchangeService) Exchange(
	ctx context.Context,
	request dto.ExchangeRequest,
) (dto.ExchangeDto, error) {
	log.Printf("exchange_service.exchange start base=%s target=%s amount=%s", request.BaseCode, request.TargetCode, request.Amount.String())
	if err := request.Validate(); err != nil {
		return dto.ExchangeDto{}, err
	}

	baseCurrency, err := s.currencyRepository.GetByCode(ctx, request.BaseCode)
	if err != nil {
		return dto.ExchangeDto{}, s.wrapCurrencyError("base currency", request.BaseCode, err)
	}

	targetCurrency, err := s.currencyRepository.GetByCode(ctx, request.TargetCode)
	if err != nil {
		return dto.ExchangeDto{}, s.wrapCurrencyError("target currency", request.TargetCode, err)
	}

	rate, err := s.exchangeRepository.GetRate(ctx, baseCurrency.ID, targetCurrency.ID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			log.Printf("exchange_service.exchange rate_not_found base_id=%d target_id=%d", baseCurrency.ID, targetCurrency.ID)
			return dto.ExchangeDto{}, fmt.Errorf("%w: exchange rate not found", errs.ErrNotFound)
		}
		log.Printf("exchange_service.exchange rate_error: %v", err)
		return dto.ExchangeDto{}, fmt.Errorf("%w: get exchange rate", errs.ErrInternal)
	}

	result := dto.ExchangeDto{
		ExchangeRate: dto.ExchangeRateDto{
			BaseCurrency: dto.CurrencyDto{
				ID:       baseCurrency.ID,
				Code:     baseCurrency.Code,
				FullName: baseCurrency.FullName,
				Sign:     baseCurrency.Sign,
			},
			TargetCurrency: dto.CurrencyDto{
				ID:       targetCurrency.ID,
				Code:     targetCurrency.Code,
				FullName: targetCurrency.FullName,
				Sign:     targetCurrency.Sign,
			},
			Rate: rate,
		},
		Amount:        request.Amount,
		ConvertAmount: request.Amount.Mul(rate),
	}
	log.Printf("exchange_service.exchange ok base=%s target=%s amount=%s converted=%s", request.BaseCode, request.TargetCode, request.Amount.String(), result.ConvertAmount.String())
	return result, nil
}

func (s *ExchangeService) wrapCurrencyError(currencyRole string, code string, err error) error {
	if errors.Is(err, errs.ErrNotFound) {
		log.Printf("exchange_service.exchange %s not_found code=%s", currencyRole, code)
		return fmt.Errorf("%w: %s not found", errs.ErrNotFound, currencyRole)
	}
	log.Printf("exchange_service.exchange %s error: %v", currencyRole, err)
	return fmt.Errorf("%w: get %s", errs.ErrInternal, currencyRole)
}

func mapRate(rate entity.ExchangeRate) dto.ExchangeRateDto {
	return dto.ExchangeRateDto{
		ID:             rate.ID,
		BaseCurrency:   mapCurrency(rate.BaseCurrency),
		TargetCurrency: mapCurrency(rate.TargetCurrency),
		Rate:           rate.Rate,
	}
}
