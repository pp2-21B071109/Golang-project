package data

import (
	"database/sql"
	"time"
	
	"errors"
	"github.com/lib/pq"
	"greenlight.dimash.net/internal/validator" // New import
)

type RareCoin struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Genres    []string  `json:"genres,omitempty"`
	Price     int32     `json:"version"`
}

func ValidateMovie(v *validator.Validator, rarecoin *RareCoin) {
	v.Check(rarecoin.Title != "", "title", "must be provided")
	v.Check(len(rarecoin.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(rarecoin.Year != 0, "year", "must be provided")
	v.Check(rarecoin.Year >= 1888, "year", "must be greater than 1888")
	v.Check(rarecoin.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(rarecoin.Genres != nil, "genres", "must be provided")
	v.Check(len(rarecoin.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(rarecoin.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(rarecoin.Genres), "genres", "must not contain duplicate values")
}

type CoinModel struct {
	DB *sql.DB
}

// Add a placeholder method for inserting a new record in the movies table.
func (m CoinModel) Insert(coin *Coin) error {
	// Define the SQL query for inserting a new record in the movies table and returning
	// the system-generated data.
	query := `
	INSERT INTO movies (title, year, runtime, genres)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, version`
	// Create an args slice containing the values for the placeholder parameters from
	// the movie struct. Declaring this slice immediately next to our SQL query helps to
	// make it nice and clear *what values are being used where* in the query.
	args := []interface{}{coin.Title, coin.Year, pq.Array(coin.Genres), coin.Price}
	// Use the QueryRow() method to execute the SQL query on our connection pool,
	// passing in the args slice as a variadic parameter and scanning the system-
	// generated id, created_at and version values into the movie struct.
	return m.DB.QueryRow(query, args...).Scan(&coin.ID, &coin.CreatedAt, &coin.Version)
}

// Add a placeholder method for fetching a specific record from the movies table.
func (m CoinModel) Get(id int64) (*Coin, error) {
	// The PostgreSQL bigserial type that we're using for the movie ID starts
	// auto-incrementing at 1 by default, so we know that no movies will have ID values
	// less than that. To avoid making an unnecessary database call, we take a shortcut
	// and return an ErrRecordNotFound error straight away.
	if id < 1 {
	return nil, ErrRecordNotFound
	}
	query := `
SELECT id, created_at, title, year, runtime, genres, version
FROM movies
WHERE id = $1`

var coin Coin
// Execute the query using the QueryRow() method, passing in the provided id value
// as a placeholder parameter, and scan the response data into the fields of the
// Movie struct. Importantly, notice that we need to convert the scan target for the
// genres column using the pq.Array() adapter function again.
err := m.DB.QueryRow(query, id).Scan(
&movie.ID,
&movie.CreatedAt,
&movie.Title,
&movie.Year,
pq.Array(&movie.Genres),
&movie.Price,
)
// Handle any errors. If there was no matching movie found, Scan() will return
// a sql.ErrNoRows error. We check for this and return our custom ErrRecordNotFound
// error instead.
if err != nil {
switch {
case errors.Is(err, sql.ErrNoRows):
return nil, ErrRecordNotFound
default:
return nil, err
}
}
// Otherwise, return a pointer to the Movie struct.
return &movie, nil
// Add a placeholder method for updating a specific record in the movies table.
func (m Coin) Update(m coin *Coin) error {
	return nil
}

// Add a placeholder method for deleting a specific record from the movies table.
func (m CoinModel) Delete(id int64) error {
	return nil
}
