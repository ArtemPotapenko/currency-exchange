package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"currency-exchange/internal/dto"
	"currency-exchange/internal/errs"
	"currency-exchange/internal/pagination"
	"currency-exchange/internal/service"

	"github.com/shopspring/decimal"
)

type CurrencyServer struct {
	currencyService *service.CurrencyService
	exchangeService *service.ExchangeService
	mux             *http.ServeMux
}

func New(mux *http.ServeMux, currencyService *service.CurrencyService, exchangeService *service.ExchangeService) http.Handler {
	s := &CurrencyServer{
		currencyService: currencyService,
		exchangeService: exchangeService,
		mux:             mux,
	}

	s.mux.HandleFunc("GET /currencies", s.handleCurrenciesList)
	s.mux.HandleFunc("POST /currencies", s.handleCurrenciesCreate)
	s.mux.HandleFunc("GET /currencies/", s.handleCurrencyGetByCode)
	s.mux.HandleFunc("POST /rates", s.handleRatesCreate)
	s.mux.HandleFunc("GET /rates/", s.handleRateGetByID)
	s.mux.HandleFunc("PUT /rates/", s.handleRateUpdateByID)
	s.mux.HandleFunc("GET /exchange", s.handleExchangeGet)

	return mux
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
func (s *CurrencyServer) handleCurrenciesList(w http.ResponseWriter, r *http.Request) {
	pageNumber, pageSize, err := parsePageRequest(r)
	if err != nil {
		writeError(w, fmt.Errorf("%w: invalid pagination", errs.ErrValidation))
		return
	}
	page, err := s.currencyService.GetAllCurrencyPage(
		r.Context(),
		pagination.PageRequest{
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
func (s *CurrencyServer) handleCurrenciesCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCurrencyRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, fmt.Errorf("%w: invalid request", errs.ErrValidation))
		return
	}
	currency, err := s.currencyService.CreateCurrency(r.Context(), req)
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
func (s *CurrencyServer) handleCurrencyGetByCode(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/currencies/")
	if code == "" {
		writeError(w, fmt.Errorf("%w: currency code is required", errs.ErrValidation))
		return
	}
	currency, err := s.currencyService.GetCurrencyByCode(r.Context(), code)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, currency)
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
func (s *CurrencyServer) handleRatesCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateRateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, fmt.Errorf("%w: invalid request", errs.ErrValidation))
		return
	}
	rate, err := s.exchangeService.CreateRate(r.Context(), req)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, rate)
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
func (s *CurrencyServer) handleRateGetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r.URL.Path, "/rates/")
	if err != nil {
		writeError(w, err)
		return
	}
	rate, err := s.exchangeService.GetRateByID(r.Context(), id)
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
func (s *CurrencyServer) handleRateUpdateByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r.URL.Path, "/rates/")
	if err != nil {
		writeError(w, err)
		return
	}
	var req dto.UpdateRateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, fmt.Errorf("%w: invalid request", errs.ErrValidation))
		return
	}
	rate, err := s.exchangeService.UpdateRate(r.Context(), id, req)
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
func (s *CurrencyServer) handleExchangeGet(w http.ResponseWriter, r *http.Request) {
	baseCode := r.URL.Query().Get("base")
	targetCode := r.URL.Query().Get("target")
	amountStr := r.URL.Query().Get("amount")
	amount, err := decimal.NewFromString(amountStr)
	if err != nil {
		writeError(w, fmt.Errorf("%w: invalid amount", errs.ErrValidation))
		return
	}
	result, err := s.exchangeService.Exchange(
		r.Context(),
		dto.ExchangeRequest{
			BaseCode:   baseCode,
			TargetCode: targetCode,
			Amount:     amount,
		})
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

func parseIDFromPath(path string, prefix string) (int64, error) {
	idStr := strings.TrimPrefix(path, prefix)
	if idStr == "" {
		return 0, fmt.Errorf("%w: rate id is required", errs.ErrValidation)
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: invalid rate id", errs.ErrValidation)
	}
	return id, nil
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
	case errors.Is(err, errs.ErrValidation):
		writeJSON(w, http.StatusBadRequest, dto.ErrorDto{Message: errorMessage(err, errs.ErrValidation)})
	case errors.Is(err, errs.ErrNotFound):
		writeJSON(w, http.StatusNotFound, dto.ErrorDto{Message: errorMessage(err, errs.ErrNotFound)})
	case errors.Is(err, errs.ErrInternal):
		writeJSON(w, http.StatusInternalServerError, dto.ErrorDto{Message: "internal error"})
	default:
		writeJSON(w, http.StatusInternalServerError, dto.ErrorDto{Message: "internal error"})
	}
}

func errorMessage(err error, sentinel error) string {
	msg := err.Error()
	prefix := sentinel.Error() + ": "
	if strings.HasPrefix(msg, prefix) {
		return strings.TrimPrefix(msg, prefix)
	}
	return msg
}
