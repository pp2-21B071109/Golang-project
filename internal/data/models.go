package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Coins  CoinModel
	Tokens TokenModel // Add a new Tokens field.
	Users  UserModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Coins:  CoinModel{DB: db},
		Tokens: TokenModel{DB: db}, // Initialize a new TokenModel instance.
		Users:  UserModel{DB: db},
	}
}
