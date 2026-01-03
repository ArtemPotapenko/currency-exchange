package main

import (
	"context"
	httpserver "currency-exchange/internal/http"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"currency-exchange/internal/repository/db"
	"currency-exchange/internal/service"

	_ "github.com/lib/pq"
)

func main() {
	addr := getEnvOrDefault("HTTP_ADDR", ":8080")

	dsn := postgresDSNFromEnv()
	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("db open error: %v", err)
	}
	if err := pingWithRetry(dbConn, 10, 2*time.Second); err != nil {
		log.Fatalf("db ping error: %v", err)
	}

	currencyRepo := db.NewCurrencyRepository(dbConn)
	exchangeRepo := db.NewExchangeRepository(dbConn)

	ctx := context.Background()
	currencyService := service.NewCurrencyService(ctx, currencyRepo)
	exchangeService := service.NewExchangeService(ctx, exchangeRepo, currencyRepo)

	handler := httpserver.LoggingMiddleware(httpserver.New(currencyService, exchangeService))

	log.Printf("http server listening on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("http server error: %v", err)
	}
}

func pingWithRetry(dbConn *sql.DB, attempts int, delay time.Duration) error {
	var err error
	for i := 0; i < attempts; i++ {
		err = dbConn.Ping()
		if err == nil {
			return nil
		}
		log.Printf("db ping failed (attempt %d/%d): %v", i+1, attempts, err)
		time.Sleep(delay)
	}
	return err
}

func postgresDSNFromEnv() string {
	host := getEnvOrDefault("PG_HOST", "localhost")
	port := getEnvOrDefault("PG_PORT", "5432")
	user := getEnvOrDefault("PG_USER", "postgres")
	password := getEnvOrDefault("PG_PASSWORD", "postgres")
	dbname := getEnvOrDefault("PG_DBNAME", "postgres")
	sslmode := getEnvOrDefault("PG_SSLMODE", "disable")

	return "host=" + host +
		" port=" + port +
		" user=" + user +
		" password=" + password +
		" dbname=" + dbname +
		" sslmode=" + sslmode
}

func getEnvOrDefault(key string, deafaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return deafaultValue
	}
	return value
}
