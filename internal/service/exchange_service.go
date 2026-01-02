package service

import (
	"context"
	"currency-exchange/internal/dto"
	"currency-exchange/internal/entity"
	apperror "currency-exchange/internal/error"
	"currency-exchange/internal/repository"
	"errors"
	"fmt"
	"log"

	"github.com/shopspring/decimal"
)

type ExchangeService struct {
	ctx                context.Context
	exchangeRepository repository.ExchangeRepository
	currencyRepository repository.CurrencyRepository
}

func NewExchangeService(
	ctx context.Context,
	exchangeRepository repository.ExchangeRepository,
	currencyRepository repository.CurrencyRepository,
) *ExchangeService {
	return &ExchangeService{
		ctx:                ctx,
		exchangeRepository: exchangeRepository,
		currencyRepository: currencyRepository,
	}
}

func (s *ExchangeService) CreateRate(baseCode string, targetCode string, rate decimal.Decimal) (dto.ExchangeRateDto, error) {
	log.Printf("exchange_service.create_rate start base=%s target=%s", baseCode, targetCode)
	baseCurrency, err := s.currencyRepository.GetByCode(s.ctx, baseCode)
	if err != nil {
		return dto.ExchangeRateDto{}, s.wrapCurrencyError("base currency", baseCode, err)
	}
	targetCurrency, err := s.currencyRepository.GetByCode(s.ctx, targetCode)
	if err != nil {
		return dto.ExchangeRateDto{}, s.wrapCurrencyError("target currency", targetCode, err)
	}

	entityRate := entity.ExchangeRate{
		BaseCurrency:   baseCurrency,
		TargetCurrency: targetCurrency,
		Rate:           rate,
	}
	id, err := s.exchangeRepository.Create(s.ctx, entityRate)
	if err != nil {
		log.Printf("exchange_service.create_rate error: %v", err)
		return dto.ExchangeRateDto{}, apperror.Internal("create exchange rate", err.Error())
	}
	entityRate.ID = id
	log.Printf("exchange_service.create_rate ok id=%d", id)
	return mapRate(entityRate), nil
}

func (s *ExchangeService) UpdateRate(id int64, baseCode string, targetCode string, rate decimal.Decimal) (dto.ExchangeRateDto, error) {
	log.Printf("exchange_service.update_rate start id=%d", id)
	baseCurrency, err := s.currencyRepository.GetByCode(s.ctx, baseCode)
	if err != nil {
		return dto.ExchangeRateDto{}, s.wrapCurrencyError("base currency", baseCode, err)
	}
	targetCurrency, err := s.currencyRepository.GetByCode(s.ctx, targetCode)
	if err != nil {
		return dto.ExchangeRateDto{}, s.wrapCurrencyError("target currency", targetCode, err)
	}

	entityRate := entity.ExchangeRate{
		ID:             id,
		BaseCurrency:   baseCurrency,
		TargetCurrency: targetCurrency,
		Rate:           rate,
	}
	if err := s.exchangeRepository.Update(s.ctx, entityRate); err != nil {
		var notFoundErr *apperror.NotFoundError
		if errors.As(err, &notFoundErr) {
			log.Printf("exchange_service.update_rate not_found id=%d", id)
			return dto.ExchangeRateDto{}, apperror.NotFound("exchange rate not found", "id="+fmt.Sprint(id))
		}
		log.Printf("exchange_service.update_rate error: %v", err)
		return dto.ExchangeRateDto{}, apperror.Internal("update exchange rate", err.Error())
	}

	log.Printf("exchange_service.update_rate ok id=%d", id)
	return mapRate(entityRate), nil
}

func (s *ExchangeService) GetRateByID(id int64) (dto.ExchangeRateDto, error) {
	log.Printf("exchange_service.get_rate_by_id start id=%d", id)
	rate, err := s.exchangeRepository.GetByID(s.ctx, id)
	if err != nil {
		var notFoundErr *apperror.NotFoundError
		if errors.As(err, &notFoundErr) {
			log.Printf("exchange_service.get_rate_by_id not_found id=%d", id)
			return dto.ExchangeRateDto{}, apperror.NotFound("exchange rate not found", "id="+fmt.Sprint(id))
		}
		log.Printf("exchange_service.get_rate_by_id error: %v", err)
		return dto.ExchangeRateDto{}, apperror.Internal("get exchange rate by id", err.Error())
	}
	log.Printf("exchange_service.get_rate_by_id ok id=%d", rate.ID)
	return mapRate(rate), nil
}

func (s *ExchangeService) Exchange(
	baseCode string,
	targetCode string,
	amount decimal.Decimal,
) (dto.ExchangeDto, error) {
	log.Printf("exchange_service.exchange start base=%s target=%s amount=%s", baseCode, targetCode, amount.String())
	if baseCode == "" || targetCode == "" {
		log.Printf("exchange_service.exchange validation_error: empty code")
		return dto.ExchangeDto{}, apperror.Validation("currency codes are required", "base or target code is empty")
	}
	if amount.LessThanOrEqual(decimal.Zero) {
		log.Printf("exchange_service.exchange validation_error: non_positive amount=%s", amount.String())
		return dto.ExchangeDto{}, apperror.Validation("amount must be greater than zero", "amount="+amount.String())
	}

	baseCurrency, err := s.currencyRepository.GetByCode(s.ctx, baseCode)
	if err != nil {
		return dto.ExchangeDto{}, s.wrapCurrencyError("base currency", baseCode, err)
	}

	targetCurrency, err := s.currencyRepository.GetByCode(s.ctx, targetCode)
	if err != nil {
		return dto.ExchangeDto{}, s.wrapCurrencyError("target currency", targetCode, err)
	}

	rate, err := s.exchangeRepository.GetRate(s.ctx, baseCurrency.ID, targetCurrency.ID)
	if err != nil {
		var notFoundErr *apperror.NotFoundError
		if errors.As(err, &notFoundErr) {
			log.Printf("exchange_service.exchange rate_not_found base_id=%d target_id=%d", baseCurrency.ID, targetCurrency.ID)
			return dto.ExchangeDto{}, apperror.NotFound("exchange rate not found", "base_id="+fmt.Sprint(baseCurrency.ID)+" target_id="+fmt.Sprint(targetCurrency.ID))
		}
		log.Printf("exchange_service.exchange rate_error: %v", err)
		return dto.ExchangeDto{}, apperror.Internal("get exchange rate", err.Error())
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
		Amount:        amount,
		ConvertAmount: amount.Mul(rate),
	}
	log.Printf("exchange_service.exchange ok base=%s target=%s amount=%s converted=%s", baseCode, targetCode, amount.String(), result.ConvertAmount.String())
	return result, nil
}

func (s *ExchangeService) wrapCurrencyError(currencyRole string, code string, err error) error {
	var notFoundErr *apperror.NotFoundError
	if errors.As(err, &notFoundErr) {
		log.Printf("exchange_service.exchange %s not_found code=%s", currencyRole, code)
		return apperror.NotFound(currencyRole+" not found", "code="+code)
	}
	log.Printf("exchange_service.exchange %s error: %v", currencyRole, err)
	return apperror.Internal("get "+currencyRole, err.Error())
}

func mapRate(rate entity.ExchangeRate) dto.ExchangeRateDto {
	return dto.ExchangeRateDto{
		ID:             rate.ID,
		BaseCurrency:   mapCurrency(rate.BaseCurrency),
		TargetCurrency: mapCurrency(rate.TargetCurrency),
		Rate:           rate.Rate,
	}
}
