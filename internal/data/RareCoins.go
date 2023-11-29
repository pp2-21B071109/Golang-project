package data

import (
	"context" // New import
	"database/sql"
	"fmt"
	"time"

	"errors"

	"github.com/lib/pq"
	"greenlight.alexedwards.net/internal/validator"
)

type Coin struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   int32     `json:"runtime,omitempty"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}

func ValidateCoin(v *validator.Validator, coin *Coin) {
	v.Check(coin.Title != "", "title", "must be provided")
	v.Check(len(coin.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(coin.Year != 0, "year", "must be provided")
	v.Check(coin.Year >= 1888, "year", "must be greater than 1888")
	v.Check(coin.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(coin.Genres != nil, "genres", "must be provided")
	v.Check(len(coin.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(coin.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(coin.Genres), "genres", "must not contain duplicate values")
}

type CoinModel struct {
	DB *sql.DB
}

// Add a placeholder method for inserting a new record in the coins table.
func (m CoinModel) Insert(coin *Coin) error {
	query := `
	INSERT INTO coins (title, year, runtime, genres)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, version`
	args := []interface{}{coin.Title, coin.Year, coin.Runtime, pq.Array(coin.Genres)}
	// Create a context with a 3-second timeout.
	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Use QueryRowContext() and pass the context as the first argument.
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&coin.ID, &coin.CreatedAt, &coin.Version)
}
func (m CoinModel) Get(id int64) (*Coin, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	// Remove the pg_sleep(10) clause.
	query := `
	SELECT id, created_at, title, year, runtime, genres, version
	FROM coins
	WHERE id = $1`
	var coin Coin
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Remove &[]byte{} from the first Scan() destination.
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&coin.ID,
		&coin.CreatedAt,
		&coin.Title,
		&coin.Year,
		&coin.Runtime,
		pq.Array(&coin.Genres),
		&coin.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &coin, nil
}
func (m CoinModel) Update(coin *Coin) error {
	query := `
	UPDATE coins
	SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
	WHERE id = $5 AND version = $6
	RETURNING version`
	args := []interface{}{
		coin.Title,
		coin.Year,
		coin.Runtime,
		pq.Array(coin.Genres),
		coin.ID,
		coin.Version,
	}
	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Use QueryRowContext() and pass the context as the first argument.
	// Use QueryRowContext() and pass the context as the first argument.
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&coin.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}
func (m CoinModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `
	DELETE FROM coins
	WHERE id = $1`
	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Use ExecContext() and pass the context as the first argument.
	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (m CoinModel) GetAll(title string, genres []string, filters Filters) ([]*Coin, Metadata, error) {
	// Update the SQL query to include the window function which counts the total
	// (filtered) records.
	// (filtered) records.
	query := fmt.Sprintf(`
	SELECT count(*) OVER(), id, created_at, title, year, runtime, genres, version
	FROM coins
	WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
	AND (genres @> $2 OR $2 = '{}')
	ORDER BY %s %s, id ASC
	LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []interface{}{title, pq.Array(genres), filters.limit(), filters.offset()}
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err // Update this to return an empty Metadata struct.
	}
	defer rows.Close()
	// Declare a totalRecords variable.
	totalRecords := 0
	coins := []*Coin{}
	for rows.Next() {
		var coin Coin
		err := rows.Scan(
			&totalRecords, // Scan the count from the window function into totalRecords.
			&coin.ID,
			&coin.CreatedAt,
			&coin.Title,
			&coin.Year,
			&coin.Runtime,
			pq.Array(&coin.Genres),
			&coin.Version,
		)
		if err != nil {
			return nil, Metadata{}, err // Update this to return an empty Metadata struct.
		}
		coins = append(coins, &coin)
	}
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err // Update this to return an empty Metadata struct.
	}
	// Generate a Metadata struct, passing in the total record count and pagination
	// parameters from the client.
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	// Include the metadata struct when returning.
	return coins, metadata, nil
}
