package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"sync"
	"text/template"
	"time"

	"github.com/Lanrey-waju/sunny-akins/internal/config"
	"github.com/Lanrey-waju/sunny-akins/internal/database"
	"github.com/Lanrey-waju/sunny-akins/internal/mailer"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
)

type application struct {
	config         config.Config
	db             *database.Queries
	templateCache  map[string]*template.Template
	mailer         mailer.Mailer
	sessionManager *scs.SessionManager
	wg             sync.WaitGroup
}

func main() {
	cfg := config.NewConfig()

	portString := os.Getenv("PORT")
	port, err := strconv.Atoi(portString)
	if err != nil {
		cfg.ErrorLog.Fatalf("unable to convert port string to int: %s", err)
	}

	smtpPortString := os.Getenv("PORT")
	SMTP_PORT, err := strconv.Atoi(smtpPortString)
	if err != nil {
		cfg.ErrorLog.Fatalf("unable to convert smtp port string to int: %s", err)
	}

	flag.IntVar(&cfg.Port, "port", port, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment(development|staging|production)")
	flag.StringVar(&cfg.DB.Dsn, "db-dsn", os.Getenv("DATABASE_URL"), "PostgreSQL DSN")
	fmt.Println(cfg.DB.Dsn)

	// smtp server config settings
	flag.StringVar(&cfg.Smtp.Host, "smtp-host", os.Getenv("SMTP_HOST"), "SMTP host")
	flag.IntVar(&cfg.Smtp.Port, "smtp-port", SMTP_PORT, "SMTP port")
	flag.StringVar(&cfg.Smtp.Username, "smtp-username", os.Getenv("SMTP_USERNAME"), "SMTP username")
	flag.StringVar(&cfg.Smtp.Password, "smtp-password", os.Getenv("SMTP_PASSWORD"), "SMTP password")
	flag.StringVar(&cfg.Smtp.Sender, "smtp-sender", "from@example.com", "SMTP sender")

	flag.Parse()

	fmt.Printf(
		"%v, %v, %v, %v\n",
		cfg.Smtp.Username,
		cfg.Smtp.Host,
		cfg.Smtp.Password,
		cfg.Smtp.Sender,
	)
	// initialize a new template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		cfg.ErrorLog.Fatalln(err)
	}

	mailer := mailer.New(
		cfg.Smtp.Host,
		cfg.Smtp.Port,
		cfg.Smtp.Username,
		cfg.Smtp.Password,
		cfg.Smtp.Sender,
	)

	db, err := cfg.NewDB()
	if err != nil {
		cfg.ErrorLog.Fatal(err)
	}

	dbQueries := database.New(db)

	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 12 * time.Second

	app := &application{
		config:         *cfg,
		db:             dbQueries,
		templateCache:  templateCache,
		mailer:         mailer,
		sessionManager: sessionManager,
	}

	err = app.serve()
	if err != nil {
		app.config.ErrorLog.Fatal(err)
	}
}
