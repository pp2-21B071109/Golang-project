package main

import (
	"fmt"
	"net/http"
	"time"

	"greenlight.dimash.net/internal/data"
	// New import
)

func (app *application) createRareCoinsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string   `json:"title"`
		Year   int32    `json:"year"`
		Genres []string `json:"genres"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showRareCoinsHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		// Use the new notFoundResponse() helper.
		app.notFoundResponse(w, r)
		return
	}
	rareCoin := data.RareCoin{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Шиллинг новая Ангия",
		Genres:    []string{"England", "12 pence"},
		Price:     414000,
	}

	// Create an envelope{"RareCoin": rareCoin} instance.
	envelope := envelope{"RareCoin": rareCoin}

	// Pass the envelope to writeJSON.
	err = app.writeJSON(w, http.StatusOK, envelope, nil)
	if err != nil {
		// Use the new serverErrorResponse() helper.
		app.serverErrorResponse(w, r, err)
	}
}
