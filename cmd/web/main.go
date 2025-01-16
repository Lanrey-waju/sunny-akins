package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"os"
	"strconv"
	"sync"
	"text/template"
	"time"

	_ "github.com/lib/pq"

	"github.com/Lanrey-waju/sunny-akins/internal/database"
	"github.com/Lanrey-waju/sunny-akins/internal/mailer"
)

// config will hold environment specific variables
type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}

type application struct {
	config        config
	errorLog      *log.Logger
	infoLog       *log.Logger
	db            *database.Queries
	templateCache map[string]*template.Template
	mailer        mailer.Mailer
	wg            sync.WaitGroup
}

// openDB opens and verifies a connection to a database
func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	// create a logger for writing informational messages
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// create a logger for error messages
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	cfg := config{}

	portString := os.Getenv("SUNNY_PORT")
	port, err := strconv.Atoi(portString)
	if err != nil {
		errorLog.Fatalf("unable to convert port string to int: %s", err)
	}

	flag.IntVar(&cfg.port, "port", port, "API server port")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("SUNNY_DSN"), "PostgreSQL DSN")
	flag.StringVar(&cfg.env, "env", "development", "Environment(development|staging|production)")

	// smtp server config settings
	flag.StringVar(&cfg.smtp.host, "smtp-host", "smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 2525, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "58989729b79228", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", os.Getenv("SMTP_PASSWORD"), "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "from@example.com", "SMTP sender")

	flag.Parse()

	db, err := openDB(cfg)
	if err != nil {
		errorLog.Fatal(err)
	}
	infoLog.Println("database connection pool established")
	defer db.Close()

	dbQueries := database.New(db)

	// initialize a new template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatalln(err)
	}

	mailer := mailer.New(
		cfg.smtp.host,
		cfg.smtp.port,
		cfg.smtp.username,
		cfg.smtp.password,
		cfg.smtp.sender,
	)

	app := &application{
		config:        cfg,
		errorLog:      errorLog,
		infoLog:       infoLog,
		db:            dbQueries,
		templateCache: templateCache,
		mailer:        mailer,
	}

	err = app.serve()
	if err != nil {
		errorLog.Fatal(err)
	}
}
