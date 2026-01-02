package db

import (
	"context"
	apperror "currency-exchange/internal/error"
	"currency-exchange/internal/pagination"
	"currency-exchange/internal/repository"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"currency-exchange/internal/entity"
)

type CurrencyRepositoryDB struct {
	db *sql.DB
}

var _ repository.CurrencyRepository = (*CurrencyRepositoryDB)(nil)

func NewCurrencyRepository(db *sql.DB) *CurrencyRepositoryDB {
	return &CurrencyRepositoryDB{db: db}
}

func (r *CurrencyRepositoryDB) Create(ctx context.Context, currency entity.Currency) (int64, error) {
	log.Printf("currency_repository.create start code=%s", currency.Code)
	row := r.db.QueryRowContext(
		ctx,
		`INSERT INTO currencies (code, full_name, sign)
		 VALUES ($1, $2, $3)
		 RETURNING id`,
		currency.Code,
		currency.FullName,
		currency.Sign,
	)

	var id int64
	if err := row.Scan(&id); err != nil {
		log.Printf("currency_repository.create error: %v", err)
		return 0, apperror.Internal("db create currency", err.Error())
	}

	log.Printf("currency_repository.create ok id=%d", id)
	return id, nil
}

func (r *CurrencyRepositoryDB) GetByID(ctx context.Context, id int64) (entity.Currency, error) {
	log.Printf("currency_repository.get_by_id start id=%d", id)
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, code, full_name, sign
		 FROM currencies
		 WHERE id = $1`,
		id,
	)

	currency, err := scanCurrency(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("currency_repository.get_by_id not_found id=%d", id)
			return entity.Currency{}, apperror.NotFound("currency not found", "id="+fmt.Sprint(id))
		}
		log.Printf("currency_repository.get_by_id error: %v", err)
		return entity.Currency{}, apperror.Internal("db get currency by id", err.Error())
	}

	log.Printf("currency_repository.get_by_id ok id=%d", currency.ID)
	return currency, nil
}

func (r *CurrencyRepositoryDB) GetByCode(ctx context.Context, code string) (entity.Currency, error) {
	log.Printf("currency_repository.get_by_code start code=%s", code)
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, code, full_name, sign
		 FROM currencies
		 WHERE code = $1`,
		code,
	)

	currency, err := scanCurrency(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("currency_repository.get_by_code not_found code=%s", code)
			return entity.Currency{}, apperror.NotFound("currency not found", "code="+code)
		}
		log.Printf("currency_repository.get_by_code error: %v", err)
		return entity.Currency{}, apperror.Internal("db get currency by code", err.Error())
	}

	log.Printf("currency_repository.get_by_code ok id=%d", currency.ID)
	return currency, nil
}

func (r *CurrencyRepositoryDB) GetAll(ctx context.Context) ([]entity.Currency, error) {
	log.Printf("currency_repository.get_all start")
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, code, full_name, sign
		 FROM currencies
		 ORDER BY id`,
	)
	if err != nil {
		log.Printf("currency_repository.get_all error: %v", err)
		return nil, apperror.Internal("db get all currencies", err.Error())
	}
	defer rows.Close()

	var currencies []entity.Currency
	for rows.Next() {
		currency, err := scanCurrency(rows)
		if err != nil {
			log.Printf("currency_repository.get_all scan_error: %v", err)
			return nil, apperror.Internal("db scan currency", err.Error())
		}
		currencies = append(currencies, currency)
	}

	if err := rows.Err(); err != nil {
		log.Printf("currency_repository.get_all iterate_error: %v", err)
		return nil, apperror.Internal("db iterate currencies", err.Error())
	}

	log.Printf("currency_repository.get_all ok count=%d", len(currencies))
	return currencies, nil
}

func (r *CurrencyRepositoryDB) GetPage(ctx context.Context, page pagination.PageRequest) (pagination.Page[entity.Currency], error) {
	log.Printf("currency_repository.get_page start page=%d size=%d", page.PageNumber, page.PageSize)
	if page.PageNumber < 1 || page.PageSize < 1 {
		log.Printf("currency_repository.get_page validation_error page=%d size=%d", page.PageNumber, page.PageSize)
		return pagination.Page[entity.Currency]{}, apperror.Validation("invalid page params", fmt.Sprintf(
			"invalid page params: pageNumber=%d pageSize=%d",
			page.PageNumber,
			page.PageSize,
		))
	}

	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM currencies`).Scan(&total); err != nil {
		log.Printf("currency_repository.get_page count_error: %v", err)
		return pagination.Page[entity.Currency]{}, apperror.Internal("db count currencies", err.Error())
	}

	limit := int64(page.PageSize)
	offset := int64(page.PageNumber-1) * int64(page.PageSize)
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, code, full_name, sign
		 FROM currencies
		 ORDER BY id
		 LIMIT $1 OFFSET $2`,
		limit,
		offset,
	)
	if err != nil {
		log.Printf("currency_repository.get_page query_error: %v", err)
		return pagination.Page[entity.Currency]{}, apperror.Internal("db get currency page", err.Error())
	}
	defer rows.Close()

	var currencies []entity.Currency
	for rows.Next() {
		currency, err := scanCurrency(rows)
		if err != nil {
			log.Printf("currency_repository.get_page scan_error: %v", err)
			return pagination.Page[entity.Currency]{}, apperror.Internal("db scan currency", err.Error())
		}
		currencies = append(currencies, currency)
	}

	if err := rows.Err(); err != nil {
		log.Printf("currency_repository.get_page iterate_error: %v", err)
		return pagination.Page[entity.Currency]{}, apperror.Internal("db iterate currency page", err.Error())
	}

	log.Printf("currency_repository.get_page ok total=%d", total)
	return pagination.Page[entity.Currency]{
		Items:      currencies,
		PageNumber: page.PageNumber,
		PageSize:   page.PageSize,
		Total:      total,
	}, nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanCurrency(scanner rowScanner) (entity.Currency, error) {
	var currency entity.Currency
	if err := scanner.Scan(
		&currency.ID,
		&currency.Code,
		&currency.FullName,
		&currency.Sign,
	); err != nil {
		return entity.Currency{}, err
	}

	return currency, nil
}
