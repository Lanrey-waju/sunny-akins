package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Lanrey-waju/sunny-akins/internal/config"
	"github.com/Lanrey-waju/sunny-akins/internal/database"
	"github.com/pressly/goose"
)

func main() {
	cfg := config.NewConfig()

	flag.StringVar(&cfg.DB.Dsn, "db-dsn", os.Getenv("DATABASE_URL"), "PostgreSQL DSN")
	flag.StringVar(&cfg.DB.Version, "version", "", "target version for migration")
	flag.StringVar(&cfg.DB.Command, "command", "up", "goose command (up|down|status|version)")
	flag.DurationVar(&cfg.DB.Timeout, "timeout", 39*time.Second, "timeout for migration operation")

	flag.Parse()

	db, err := database.OpenDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	schemaPath := "app/sql/schema"

	// Execute migration command
	switch cfg.DB.Command {
	case "up":
		cfg.InfoLog.Println("running all pending migrations...")
		if err := goose.Up(db, schemaPath); err != nil {
			cfg.ErrorLog.Fatalf("failed to migrate up: %v", err)
		}

	case "down":
		cfg.InfoLog.Println("rolling back last migration...")
		if err := goose.Down(db, schemaPath); err != nil {
			cfg.ErrorLog.Fatalf("failed to migrate down: %v", err)
		}

	case "status":
		cfg.InfoLog.Println("getting migration status...")
		if err := goose.Status(db, schemaPath); err != nil {
			cfg.ErrorLog.Fatalf("failed to get status: %v", err)
		}

	case "version":
		if cfg.DB.Version == "" {
			log.Fatal("version flag required for version command")
		}
		targetVersion, err := strconv.ParseInt(cfg.DB.Version, 10, 64)
		if err != nil {
			cfg.ErrorLog.Fatalf("invalid version number: %v", err)
		}
		log.Printf("migrating to version %d...", targetVersion)
		if err := goose.UpTo(db, schemaPath, targetVersion); err != nil {
			cfg.ErrorLog.Fatalf("failed to migrate to version %d: %v", targetVersion, err)
		}

	default:
		cfg.ErrorLog.Fatalf("unknown command: %s", cfg.DB.Command)
	}

	cfg.InfoLog.Printf("successfully completed migration command: %s", cfg.DB.Command)
	defer db.Close()
}
