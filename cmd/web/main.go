package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"

	_ "github.com/lib/pq"

	"github.com/Lanrey-waju/sunny-akins/internal/database"
)

// config will hold environment specific variables
type config struct {
	port int
	db   struct {
		dsn string
	}
}

type application struct {
	config        config
	errorLog      *log.Logger
	infoLog       *log.Logger
	db            *database.Queries
	templateCache map[string]*template.Template
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

	app := &application{
		config:        cfg,
		errorLog:      errorLog,
		infoLog:       infoLog,
		db:            dbQueries,
		templateCache: templateCache,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID(tls.X25519, tls.CurveP256),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	infoLog.Printf("Starting server on %s", srv.Addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}
