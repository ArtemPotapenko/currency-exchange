package db

import (
	"context"
	"currency-exchange/internal/errs"
	"currency-exchange/internal/pagination"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"currency-exchange/internal/entity"
)

type CurrencyRepository interface {
	Create(ctx context.Context, currency entity.Currency) (int64, error)
	GetByCode(ctx context.Context, code string) (entity.Currency, error)
	GetPage(ctx context.Context, page pagination.PageRequest) (pagination.Page[entity.Currency], error)
}

type CurrencyRepositoryDB struct {
	db *sql.DB
}

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
		return 0, fmt.Errorf("%w: db create currency", errs.ErrInternal)
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
			return entity.Currency{}, fmt.Errorf("%w: currency not found", errs.ErrNotFound)
		}
		log.Printf("currency_repository.get_by_id error: %v", err)
		return entity.Currency{}, fmt.Errorf("%w: db get currency by id", errs.ErrInternal)
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
			return entity.Currency{}, fmt.Errorf("%w: currency not found", errs.ErrNotFound)
		}
		log.Printf("currency_repository.get_by_code error: %v", err)
		return entity.Currency{}, fmt.Errorf("%w: db get currency by code", errs.ErrInternal)
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
		return nil, fmt.Errorf("%w: db get all currencies", errs.ErrInternal)
	}
	defer rows.Close()

	var currencies []entity.Currency
	for rows.Next() {
		currency, err := scanCurrency(rows)
		if err != nil {
			log.Printf("currency_repository.get_all scan_error: %v", err)
			return nil, fmt.Errorf("%w: db scan currency", errs.ErrInternal)
		}
		currencies = append(currencies, currency)
	}

	if err := rows.Err(); err != nil {
		log.Printf("currency_repository.get_all iterate_error: %v", err)
		return nil, fmt.Errorf("%w: db iterate currencies", errs.ErrInternal)
	}

	log.Printf("currency_repository.get_all ok count=%d", len(currencies))
	return currencies, nil
}

func (r *CurrencyRepositoryDB) GetPage(ctx context.Context, page pagination.PageRequest) (pagination.Page[entity.Currency], error) {
	log.Printf("currency_repository.get_page start page=%d size=%d", page.PageNumber, page.PageSize)
	if page.PageNumber < 1 || page.PageSize < 1 {
		log.Printf("currency_repository.get_page validation_error page=%d size=%d", page.PageNumber, page.PageSize)
		return pagination.Page[entity.Currency]{}, fmt.Errorf("%w: invalid page params", errs.ErrValidation)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM currencies`).Scan(&total); err != nil {
		log.Printf("currency_repository.get_page count_error: %v", err)
		return pagination.Page[entity.Currency]{}, fmt.Errorf("%w: db count currencies", errs.ErrInternal)
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
		return pagination.Page[entity.Currency]{}, fmt.Errorf("%w: db get currency page", errs.ErrInternal)
	}
	defer rows.Close()

	var currencies []entity.Currency
	for rows.Next() {
		currency, err := scanCurrency(rows)
		if err != nil {
			log.Printf("currency_repository.get_page scan_error: %v", err)
			return pagination.Page[entity.Currency]{}, fmt.Errorf("%w: db scan currency", errs.ErrInternal)
		}
		currencies = append(currencies, currency)
	}

	if err := rows.Err(); err != nil {
		log.Printf("currency_repository.get_page iterate_error: %v", err)
		return pagination.Page[entity.Currency]{}, fmt.Errorf("%w: db iterate currency page", errs.ErrInternal)
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
