package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os" // New import
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"greenlight.alexedwards.net/internal/data"
	"greenlight.alexedwards.net/internal/jsonlog"
	"greenlight.alexedwards.net/internal/mailer"
	"greenlight.alexedwards.net/internal/validator"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	limiter struct {
		enabled bool
		rps     float64
		burst   int
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
	// Add a cors struct and trustedOrigins field with the type []string.
	cors struct {
		trustedOrigins []string
	}
}

// Update the application struct to hold a new Mailer instance.
type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	// Read the SMTP server configuration settings into the config struct, using the
	// Mailtrap settings as the default values. IMPORTANT: If you're following along,
	// make sure to replace the default values for smtp-username and smtp-password
	// with your own Mailtrap credentials.
	flag.StringVar(&cfg.smtp.host, "smtp-host", "smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 2525, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "13aba2ce443ca4", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "3be67c83b7ebe2", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Greenlight <no-reply@greenlight.alexedwards.net>", "SMTP sender")
	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})
	flag.Parse()
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()
	logger.PrintInfo("database connection pool established", nil)
	// Initialize a new Mailer instance using the settings from the command line
	// flags, and add it to the application struct.
	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}
	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s sslmode=disable", cfg.db.dsn)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Set the maximum number of concurrently open connections (in-use + idle)
	// to the given value.
	db.SetMaxOpenConns(cfg.db.maxOpenConns)

	// Set the maximum number of concurrently idle connections to the given value.
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	// Set the maximum idle time for a connection.
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	// Ping the database to check the DSN provided by the user is valid and the
	// database server is reachable.
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (app *application) updateсoinHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the coin ID from the URL.
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Fetch the existing coin record from the database, sending a 404 Not Found
	// response to the client if we couldn't find a matching record.
	coin, err := app.models.Coins.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Declare an input struct to hold the expected data from the client.
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	// Read the JSON request body data into the input struct.
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Copy the values from the request body to the appropriate fields of the coin
	// record.
	// Copy the values from the request body to the appropriate fields of the coin
	// record.
	coin.Title = input.Title
	coin.Year = input.Year
	coin.Runtime = int32(input.Runtime)
	coin.Genres = input.Genres

	// Validate the updated coin record, sending the client a 422 Unprocessable Entity
	// response if any checks fail.
	v := validator.New()
	if data.ValidateCoin(v, coin); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Pass the updated coin record to our new Update() method.
	err = app.models.Coins.Update(coin)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Write the updated coin record in a JSON response.
	err = app.writeJSON(w, http.StatusOK, envelope{"coin": coin}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
