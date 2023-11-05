package main

import (
	"context"      // New import
	"database/sql" // New import
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"greenlight.dimash.net/internal/data" // New import
	// Import the pq driver so that it can register itself with the database/sql
	// package. Note that we alias this import to the blank identifier, to stop the Go
	// compiler complaining that the package isn't being used.
	_ "github.com/lib/pq"
)

const version = "1.0.0"


type application struct {
	config config
	logger *log.Logger
	models data.Models
	}
type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")
	flag.Parse()
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	db, err := openDB(cfg)
	if err != nil {
	logger.Fatal(err)
	}
	defer db.Close()
	logger.Printf("database connection pool established")
	// Use the data.NewModels() function to initialize a Models struct, passing in the
	// connection pool as a parameter.
	app := &application{
	config: cfg,
	logger: logger,
	models: data.NewModels(db),
	}
	srv := &http.Server{
	Addr: fmt.Sprintf(":%d", cfg.port),
	Handler: app.routes(),
	IdleTimeout: time.Minute,
	ReadTimeout: 10 * time.Second,
	WriteTimeout: 30 * time.Second,
	}
	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)
	}

db, err := openDB(cfg)
if err != nil {
logger.Fatal(err)
}
defer db.Close()
logger.Printf("database connection pool established")
migrationDriver, err := postgres.WithInstance(db, &postgres.Config{})
if err != nil {
logger.PrintFatal(err, nil)
}
migrator, err := migrate.NewWithDatabaseInstance("file:///path/to/your/migrations", "postgres", migrationDriver)
if err != nil {
logger.PrintFatal(err, nil)
}
err = migrator.Up()
if err != nil && err != migrate.ErrNoChange {
logger.PrintFatal(err, nil)
}
logger.Printf("database migrations applied")
}
