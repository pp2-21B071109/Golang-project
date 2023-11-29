package main

import (
	"errors"
	"net/http"

	"greenlight.alexedwards.net/internal/validator"
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

func (app *application) deleteCoinHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the Coin ID from the URL.
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Delete the Coin from the database, sending a 404 Not Found response to the
	// client if there isn't a matching record.
	err = app.models.Coins.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Return a 200 OK status code along with a success message.
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "Coin successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) updateCoinHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Retrieve the coin record as normal.
	coin, err := app.models.Coins.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		err = app.models.Coins.Update(coin)
if err != nil {
switch {
case errors.Is(err, data.ErrEditConflict):
app.editConflictResponse(w, r)
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
	}
	// Use pointers for the Title, Year and Runtime fields.
	var input struct {
		Title   *string       `json:"title"`
		Year    *int32        `json:"year"`
		Runtime *data.Runtime `json:"runtime"`
		Genres  []string      `json:"genres"`
	}
	// Decode the JSON as normal.
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// If the input.Title value is nil then we know that no corresponding "title" key/
	// value pair was provided in the JSON request body. So we move on and leave the
	// Coin record unchanged. Otherwise, we update the Coin record with the new title
	// value. Importantly, because input.Title is a now a pointer to a string, we need
	// to dereference the pointer using the * operator to get the underlying value
	// before assigning it to our Coin record.
	if input.Title != nil {
		coin.Title = *input.Title
	}
	// We also do the same for the other fields in the input struct.
	if input.Year != nil {
		coin.Year = *input.Year
	}
	if input.Runtime != nil {
		coin.Runtime = *input.Runtime
	}
	if input.Genres != nil {
		coin.Genres = input.Genres // Note that we don't need to dereference a slice.
	}
	v := validator.New()
	if data.ValidateCoin(v, coin); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Coins.Update(coin)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"Coin": coin}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	func (app *application) listCoinsHandler(w http.ResponseWriter, r *http.Request) {
		var input struct {
		Title string
		Genres []string
		data.Filters
		}
		v := validator.New()
		qs := r.URL.Query()
		input.Title = app.readString(qs, "title", "")
		input.Genres = app.readCSV(qs, "genres", []string{})
		input.Filters.Page = app.readInt(qs, "page", 1, v)
		input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
		input.Filters.Sort = app.readString(qs, "sort", "id")
		input.Filters.SortSafelist = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}
		if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
		}
		// Accept the metadata struct as a return value.
		coins, metadata, err := app.models.Coins.GetAll(input.Title, input.Genres, input.Filters)
		if err != nil {
		app.serverErrorResponse(w, r, err)
		return
		}
		// Include the metadata in the response envelope.
		err = app.writeJSON(w, http.StatusOK, envelope{"coins": coins, "metadata": metadata}, nil)
		if err != nil {
		app.serverErrorResponse(w, r, err)
		}
		}