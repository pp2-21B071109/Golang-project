package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	// Set the Movies field to be an interface containing the methods that both the
	// 'real' model and mock model need to support.
	Movies interface {
		Insert(coin *Coin) error
		Get(id int64) (*Coin, error)
		Update(coin *Coin) error
		Delete(id int64) error
	}
}

// For ease of use, we also add a New() method which returns a Models struct containing
// the initialized MovieModel.
func NewModels(db *sql.DB) Models {
	return Models{
		Coins: CoinModel{DB: db},
	}
}
