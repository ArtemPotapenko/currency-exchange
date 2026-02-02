package main

import (
	"context"
	"currency-exchange/internal/config"
	httpserver "currency-exchange/internal/http"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
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

	mux := http.NewServeMux()
	handler := httpserver.LoggingMiddleware(httpserver.New(mux, currencyService, exchangeService))

	log.Printf("http server listening on %s", addr)
	server := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server error: %v", err)
		}
	}()

	stopCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	<-stopCtx.Done()
	log.Printf("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("http server shutdown error: %v", err)
	}
}
