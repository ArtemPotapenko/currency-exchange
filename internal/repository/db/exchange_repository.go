package db

import (
	"context"
	apperror "currency-exchange/internal/error"
	"currency-exchange/internal/repository"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"currency-exchange/internal/entity"

	"github.com/shopspring/decimal"
)

type ExchangeRepositoryDB struct {
	db *sql.DB
}

var _ repository.ExchangeRepository = (*ExchangeRepositoryDB)(nil)

func NewExchangeRepository(db *sql.DB) *ExchangeRepositoryDB {
	return &ExchangeRepositoryDB{db: db}
}

func (r *ExchangeRepositoryDB) Create(ctx context.Context, rate entity.ExchangeRate) (int64, error) {
	log.Printf("exchange_repository.create start base_id=%d target_id=%d", rate.BaseCurrency.ID, rate.TargetCurrency.ID)
	row := r.db.QueryRowContext(
		ctx,
		`INSERT INTO exchange_rates (base_currency_id, target_currency_id, rate)
		 VALUES ($1, $2, $3)
		 RETURNING id`,
		rate.BaseCurrency.ID,
		rate.TargetCurrency.ID,
		rate.Rate,
	)

	var id int64
	if err := row.Scan(&id); err != nil {
		log.Printf("exchange_repository.create error: %v", err)
		return 0, apperror.Internal("db create exchange rate", err.Error())
	}

	log.Printf("exchange_repository.create ok id=%d", id)
	return id, nil
}

func (r *ExchangeRepositoryDB) Update(ctx context.Context, rate entity.ExchangeRate) error {
	log.Printf("exchange_repository.update start id=%d", rate.ID)
	result, err := r.db.ExecContext(
		ctx,
		`UPDATE exchange_rates
		 SET base_currency_id = $1,
		     target_currency_id = $2,
		     rate = $3
		 WHERE id = $4`,
		rate.BaseCurrency.ID,
		rate.TargetCurrency.ID,
		rate.Rate,
		rate.ID,
	)
	if err != nil {
		log.Printf("exchange_repository.update error: %v", err)
		return apperror.Internal("db update exchange rate", err.Error())
	}

	affected, err := result.RowsAffected()
	if err != nil {
		log.Printf("exchange_repository.update rows_affected_error: %v", err)
		return apperror.Internal("db check exchange rate update", err.Error())
	}
	if affected == 0 {
		log.Printf("exchange_repository.update not_found id=%d", rate.ID)
		return apperror.NotFound("exchange rate not found", "id="+fmt.Sprint(rate.ID))
	}

	log.Printf("exchange_repository.update ok id=%d", rate.ID)
	return nil
}

func (r *ExchangeRepositoryDB) GetByID(ctx context.Context, id int64) (entity.ExchangeRate, error) {
	log.Printf("exchange_repository.get_by_id start id=%d", id)
	row := r.db.QueryRowContext(
		ctx,
		`SELECT er.id,
		        er.rate,
		        bc.id, bc.code, bc.full_name, bc.sign,
		        tc.id, tc.code, tc.full_name, tc.sign
		 FROM exchange_rates er
		 JOIN currencies bc ON bc.id = er.base_currency_id
		 JOIN currencies tc ON tc.id = er.target_currency_id
		 WHERE er.id = $1`,
		id,
	)

	rate, err := scanExchangeRates(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("exchange_repository.get_by_id not_found id=%d", id)
			return entity.ExchangeRate{}, apperror.NotFound("exchange rate not found", "id="+fmt.Sprint(id))
		}
		log.Printf("exchange_repository.get_by_id error: %v", err)
		return entity.ExchangeRate{}, apperror.Internal("db get exchange rate by id", err.Error())
	}

	log.Printf("exchange_repository.get_by_id ok id=%d", rate.ID)
	return rate, nil
}

func (r *ExchangeRepositoryDB) GetRate(ctx context.Context, baseId int64, targetId int64) (decimal.Decimal, error) {
	log.Printf("exchange_repository.get_rate start base_id=%d target_id=%d", baseId, targetId)
	var rate decimal.Decimal
	if err := r.db.QueryRowContext(
		ctx,
		`WITH normalized_rates AS (
			SELECT base_currency_id AS from_id,
			       target_currency_id AS to_id,
			       rate::numeric AS rate
			FROM exchange_rates
			UNION ALL
			SELECT target_currency_id AS from_id,
			       base_currency_id AS to_id,
			       1 / rate::numeric AS rate
			FROM exchange_rates
		),
		direct_rate AS (
			SELECT rate, 0 AS priority
			FROM normalized_rates
			WHERE from_id = $1 AND to_id = $2
		),
		one_hop_rate AS (
			SELECT r1.rate * r2.rate AS rate, 1 AS priority
			FROM normalized_rates r1
			JOIN normalized_rates r2 ON r1.to_id = r2.from_id
			WHERE r1.from_id = $1 AND r2.to_id = $2
		)
		SELECT rate
		FROM (
			SELECT * FROM direct_rate
			UNION ALL
			SELECT * FROM one_hop_rate
		) candidate_rates
		ORDER BY priority
		LIMIT 1`,
		baseId,
		targetId,
	).Scan(&rate); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("exchange_repository.get_rate not_found base_id=%d target_id=%d", baseId, targetId)
			return decimal.Decimal{}, apperror.NotFound("exchange rate not found", "base_id="+fmt.Sprint(baseId)+" target_id="+fmt.Sprint(targetId))
		}
		log.Printf("exchange_repository.get_rate error: %v", err)
		return decimal.Decimal{}, apperror.Internal("db get exchange rate", err.Error())
	}

	log.Printf("exchange_repository.get_rate ok rate=%s", rate.String())
	return rate, nil
}

func scanExchangeRates(scanner rowScanner) (entity.ExchangeRate, error) {
	var rate entity.ExchangeRate
	if err := scanner.Scan(
		&rate.ID,
		&rate.Rate,
		&rate.BaseCurrency.ID,
		&rate.BaseCurrency.Code,
		&rate.BaseCurrency.FullName,
		&rate.BaseCurrency.Sign,
		&rate.TargetCurrency.ID,
		&rate.TargetCurrency.Code,
		&rate.TargetCurrency.FullName,
		&rate.TargetCurrency.Sign,
	); err != nil {
		return entity.ExchangeRate{}, err
	}

	return rate, nil
}
