package database

import (
	"context"
	"database/sql"
	"log"

	"github.com/Lanrey-waju/sunny-akins/internal/config"
)

// openDB opens and verifies a connection to a database
func OpenDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DB.Dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.DB.Timeout)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	log.Println("database connection pool established")
	return db, nil
}
