package main

import (
	"errors"
	"net/http"

	"greenlight.dimash.net/internal/data"
	// New import
)

type MockCoinModel struct{}

func (m MockCoinModel) Insert(coin *Coin) error {
	// Mock the action...
}
func (m MockCoinModel) Get(id int64) (*Coin, error) {
	// Mock the action...
}
func (m MockCoinModel) Update(coin *Coin) error {
	// Mock the action...
}
func (m MockCoinModel) Delete(id int64) error {
	// Mock the action...
}

func (app *application) showMCoinHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
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
	err = app.writeJSON(w, http.StatusOK, envelope{"coin": coin}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
