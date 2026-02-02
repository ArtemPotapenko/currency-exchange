package postgres

import (
	"database/sql"
	"log"
	"time"
)

func PingWithRetry(dbConn *sql.DB, attempts int, delay time.Duration) error {
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
