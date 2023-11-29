package main

package main
import (
"context"
"database/sql"
"flag"
"os"
"time"
"greenlight.alexedwards.net/internal/data"
"greenlight.alexedwards.net/internal/jsonlog"
_ "github.com/lib/pq"
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
	rps float64
	burst int
	enabled bool
	}
	}
type application struct {
	config config
	logger *jsonlog.Logger
	logger *jsonlog.Logger
	models data.Models
	}
	func main() {
		var cfg config
		flag.IntVar(&cfg.port, "port", 4000, "API server port")
		flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
		flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")
		flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
		flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
		flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")
		flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
		flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
		flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")
		flag.Parse()
		logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
		db, err := openDB(cfg)
		if err != nil {
		logger.PrintFatal(err, nil)
		}
		defer db.Close()
		logger.PrintInfo("database connection pool established", nil)
		app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		}
		// Call app.serve() to start the server.
		err = app.serve()
		if err != nil {
		logger.PrintFatal(err, nil)
		}
		}
func (app *application) updatecoinHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the coin ID from the URL.
	id, err := app.readIDParam(r)
	if err != nil {
	app.notFoundResponse(w, r)
	return
	}
	// Fetch the existing coin record from the database, sending a 404 Not Found
	// response to the client if we couldn't find a matching record.
	coin, err := app.models.coins.Get(id)
	if err != nil {
	switch {
	case errors.Is(err, data.ErrRecordNotFound):
	app.notFoundResponse(w, r)
	default:
	app.serverErrorResponse(w, r, err)
	}
	return
	}
	}
	// Declare an input struct to hold the expected data from the client.
	var input struct {
	Title string `json:"title"`
	Year int32 `json:"year"`
	Runtime data.Runtime `json:"runtime"`
	Genres []string `json:"genres"`
	}
	// Read the JSON request body data into the input struct.
	err = app.readJSON(w, r, &input)
	if err != nil {
	app.badRequestResponse(w, r, err)
	return
	}
	// Copy the values from the request body to the appropriate fields of the coin
	// record.
	coin.Title = input.Title
	coin.Year = input.Year
	coin.Runtime = input.Runtime
	coin.Genres = input.Genres
	// Validate the updated coin record, sending the client a 422 Unprocessable Entity
	// response if any checks fail.
	v := validator.New()
	if data.Validatecoin(v, coin); !v.Valid() {
	app.failedValidationResponse(w, r, v.Errors)
	return
	}
	// Pass the updated coin record to our new Update() method.
	err = app.models.coins.Update(coin)
	if err != nil {
	app.serverErrorResponse(w, r, err)
	return
	}
	// Write the updated coin record in a JSON response.
	err = app.writeJSON(w, http.StatusOK, envelope{"coin": coin}, nil)
	if err != nil {
	app.serverErrorResponse(w, r, err)
	}
	