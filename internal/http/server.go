package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"currency-exchange/internal/dto"
	apperror "currency-exchange/internal/error"
	"currency-exchange/internal/pagination"
	"currency-exchange/internal/service"

	"github.com/shopspring/decimal"
)

type CurrencyServer struct {
	currencyService *service.CurrencyService
	exchangeService *service.ExchangeService
	mux             *http.ServeMux
}

func New(currencyService *service.CurrencyService, exchangeService *service.ExchangeService) http.Handler {
	s := &CurrencyServer{
		currencyService: currencyService,
		exchangeService: exchangeService,
		mux:             http.NewServeMux(),
	}

	s.mux.HandleFunc("/currencies", s.handleCurrencies)
	s.mux.HandleFunc("/currencies/", s.handleCurrencyByCode)
	s.mux.HandleFunc("/rates", s.handleRates)
	s.mux.HandleFunc("/rates/", s.handleRateByID)
	s.mux.HandleFunc("/exchange", s.handleExchange)

	return s
}

func (s *CurrencyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *CurrencyServer) handleCurrencies(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.handleCurrenciesPost(w, r)
	case http.MethodGet:
		s.handleCurrenciesGet(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// @Summary List currencies
// @Tags currencies
// @Accept json
// @Produce json
// @Param pageNumber query int false "Page number" minimum(1)
// @Param pageSize query int false "Page size" minimum(1)
// @Success 200 {object} dto.CurrencyPageDto
// @Failure 400 {object} dto.ErrorDto
// @Failure 500 {object} dto.ErrorDto
// @Router /currencies [get]
func (s *CurrencyServer) handleCurrenciesGet(w http.ResponseWriter, r *http.Request) {
	pageNumber, pageSize, err := parsePageRequest(r)
	if err != nil {
		writeError(w, apperror.Validation("invalid pagination", err.Error()))
		return
	}
	page, err := s.currencyService.GetAllCurrencyPage(pagination.PageRequest{
		PageNumber: pageNumber,
		PageSize:   pageSize,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, page)
}

// @Summary Create currency
// @Tags currencies
// @Accept json
// @Produce json
// @Param request body dto.CreateCurrencyRequest true "Currency payload"
// @Success 201 {object} dto.CurrencyDto
// @Failure 400 {object} dto.ErrorDto
// @Failure 500 {object} dto.ErrorDto
// @Router /currencies [post]
func (s *CurrencyServer) handleCurrenciesPost(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCurrencyRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, apperror.Validation("invalid request", err.Error()))
		return
	}
	currency, err := s.currencyService.CreateCurrency(req.Code, req.FullName, req.Sign)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, currency)
}

// @Summary Get currency by code
// @Tags currencies
// @Accept json
// @Produce json
// @Param code path string true "Currency code"
// @Success 200 {object} dto.CurrencyDto
// @Failure 400 {object} dto.ErrorDto
// @Failure 404 {object} dto.ErrorDto
// @Failure 500 {object} dto.ErrorDto
// @Router /currencies/{code} [get]
func (s *CurrencyServer) handleCurrencyByCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	code := strings.TrimPrefix(r.URL.Path, "/currencies/")
	if code == "" {
		writeError(w, apperror.Validation("currency code is required", "empty code"))
		return
	}
	currency, err := s.currencyService.GetCurrencyByCode(code)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, currency)
}

func (s *CurrencyServer) handleRates(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.handleRatesPost(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// @Summary Create exchange rate
// @Tags rates
// @Accept json
// @Produce json
// @Param request body dto.CreateRateRequest true "Rate payload"
// @Success 201 {object} dto.ExchangeRateDto
// @Failure 400 {object} dto.ErrorDto
// @Failure 404 {object} dto.ErrorDto
// @Failure 500 {object} dto.ErrorDto
// @Router /rates [post]
func (s *CurrencyServer) handleRatesPost(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateRateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, apperror.Validation("invalid request", err.Error()))
		return
	}
	rate, err := s.exchangeService.CreateRate(req.BaseCode, req.TargetCode, req.Rate)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, rate)
}

func (s *CurrencyServer) handleRateByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/rates/")
	if idStr == "" {
		writeError(w, apperror.Validation("rate id is required", "empty id"))
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, apperror.Validation("invalid rate id", err.Error()))
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.handleRateByIDGet(w, r, id)
	case http.MethodPut:
		s.handleRateByIDPut(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// @Summary Get exchange rate by id
// @Tags rates
// @Accept json
// @Produce json
// @Param id path int true "Rate ID"
// @Success 200 {object} dto.ExchangeRateDto
// @Failure 400 {object} dto.ErrorDto
// @Failure 404 {object} dto.ErrorDto
// @Failure 500 {object} dto.ErrorDto
// @Router /rates/{id} [get]
func (s *CurrencyServer) handleRateByIDGet(w http.ResponseWriter, r *http.Request, id int64) {
	rate, err := s.exchangeService.GetRateByID(id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, rate)
}

// @Summary Update exchange rate
// @Tags rates
// @Accept json
// @Produce json
// @Param id path int true "Rate ID"
// @Param request body dto.UpdateRateRequest true "Rate payload"
// @Success 200 {object} dto.ExchangeRateDto
// @Failure 400 {object} dto.ErrorDto
// @Failure 404 {object} dto.ErrorDto
// @Failure 500 {object} dto.ErrorDto
// @Router /rates/{id} [put]
func (s *CurrencyServer) handleRateByIDPut(w http.ResponseWriter, r *http.Request, id int64) {
	var req dto.UpdateRateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, apperror.Validation("invalid request", err.Error()))
		return
	}
	rate, err := s.exchangeService.UpdateRate(id, req.BaseCode, req.TargetCode, req.Rate)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, rate)
}

// @Summary Exchange currency
// @Tags exchange
// @Accept json
// @Produce json
// @Param base query string true "Base currency code"
// @Param target query string true "Target currency code"
// @Param amount query string true "Amount to exchange"
// @Success 200 {object} dto.ExchangeDto
// @Failure 400 {object} dto.ErrorDto
// @Failure 404 {object} dto.ErrorDto
// @Failure 500 {object} dto.ErrorDto
// @Router /exchange [get]
func (s *CurrencyServer) handleExchange(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	baseCode := r.URL.Query().Get("base")
	targetCode := r.URL.Query().Get("target")
	amountStr := r.URL.Query().Get("amount")
	amount, err := decimal.NewFromString(amountStr)
	if err != nil {
		writeError(w, apperror.Validation("invalid amount", err.Error()))
		return
	}
	result, err := s.exchangeService.Exchange(baseCode, targetCode, amount)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func parsePageRequest(r *http.Request) (int32, int32, error) {
	pageNumberStr := r.URL.Query().Get("pageNumber")
	pageSizeStr := r.URL.Query().Get("pageSize")
	if pageNumberStr == "" && pageSizeStr == "" {
		return 1, 20, nil
	}
	pageNumber, err := strconv.ParseInt(pageNumberStr, 10, 32)
	if err != nil {
		return 0, 0, err
	}
	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil {
		return 0, 0, err
	}
	return int32(pageNumber), int32(pageSize), nil
}

func decodeJSON(r *http.Request, target any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(target)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, apperror.ErrValidation):
		var validationErr *apperror.ValidationError
		if errors.As(err, &validationErr) {
			writeJSON(w, http.StatusBadRequest, dto.ErrorDto{Message: validationErr.Message})
			return
		}
		writeJSON(w, http.StatusBadRequest, dto.ErrorDto{Message: "validation error"})
	case errors.Is(err, apperror.ErrNotFound):
		var notFoundErr *apperror.NotFoundError
		if errors.As(err, &notFoundErr) {
			writeJSON(w, http.StatusNotFound, dto.ErrorDto{Message: notFoundErr.Message})
			return
		}
		writeJSON(w, http.StatusNotFound, dto.ErrorDto{Message: "not found"})
	case errors.Is(err, apperror.ErrInternal):
		var internalErr *apperror.InternalError
		if errors.As(err, &internalErr) {
			writeJSON(w, http.StatusInternalServerError, dto.ErrorDto{Message: internalErr.Message})
			return
		}
		writeJSON(w, http.StatusInternalServerError, dto.ErrorDto{Message: "internal error"})
	default:
		writeJSON(w, http.StatusInternalServerError, dto.ErrorDto{Message: "internal error"})
	}
}
