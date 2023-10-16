package data

import (
	"greenlight/internal/validator" // New import
	"time"
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
