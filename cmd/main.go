package main

import (
	"currency-exchange/internal/config"
	httpserver "currency-exchange/internal/http"
	"database/sql"
	"log"
	"net/http"
	"time"

	"currency-exchange/internal/repository/db"
	"currency-exchange/internal/service"
	"currency-exchange/pkg/postgres"

	_ "github.com/lib/pq"
)

func main() {
	addr := config.HTTPAddr()

	pgConfig := config.PostgresConfigFromEnv()
	dsn := pgConfig.DSN()
	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("db open error: %v", err)
	}
	if err := postgres.PingWithRetry(dbConn, 10, 2*time.Second); err != nil {
		log.Fatalf("db ping error: %v", err)
	}

	currencyRepo := db.NewCurrencyRepository(dbConn)
	exchangeRepo := db.NewExchangeRepository(dbConn)

	currencyService := service.NewCurrencyService(currencyRepo)
	exchangeService := service.NewExchangeService(exchangeRepo, currencyRepo)

	handler := httpserver.LoggingMiddleware(httpserver.New(currencyService, exchangeService))

	log.Printf("http server listening on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("http server error: %v", err)
	}
}
