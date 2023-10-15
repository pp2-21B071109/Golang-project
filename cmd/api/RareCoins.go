package main

import (
	"fmt"
	"greenlight/internal/data" // New import
	"net/http"
	"strconv"
	"time" // New import

	"github.com/julienschmidt/httprouter"
)

// Add a createMovieHandler for the "POST /v1/movies" endpoint. For now we simply
// return a plain-text placeholder response.
func (app *application) createRareCoinsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new RareCoin")
}

// Add a showMovieHandler for the "GET /v1/movies/:id" endpoint. For now, we retrieve
// the interpolated "id" parameter from the current URL and include it in a placeholder
// response.
func (app *application) showRareCoinsHandler(w http.ResponseWriter, r *http.Request) {
	// When httprouter is parsing a request, any interpolated URL parameters will be
	// stored in the request context. We can use the ParamsFromContext() function to
	// retrieve a slice containing these parameter names and values.
	params := httprouter.ParamsFromContext(r.Context())
	// We can then use the ByName() method to get the value of the "id" parameter from
	// the slice. In our project all movies will have a unique positive integer ID, but
	// the value returned by ByName() is always a string. So we try to convert it to a
	// base 10 integer (with a bit size of 64). If the parameter couldn't be converted,
	// or is less than 1, we know the ID is invalid so we use the http.NotFound()
	// function to return a 404 Not Found response.
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	RareCoin := data.RareCoin{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Шиллинг новая Ангия",
		Genres:    []string{"England", "12 pence"},
		Price:     414000,
	}
	// Encode the struct to JSON and send it as the HTTP response.
	err = app.writeJSON(w, http.StatusOK, movie, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
	func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
		id, err := app.readIDParam(r)
		if err != nil {
		http.NotFound(w, r)
		return
		}
		RareCoin := data.RareCoin{
			ID:        id,
			CreatedAt: time.Now(),
			Title:     "Шиллинг новая Ангия",
			Genres:    []string{"England", "12 pence"},
			Price:     414000,
		}
		}
		// Create an envelope{"movie": movie} instance and pass it to writeJSON(), instead
		// of passing the plain movie struct.
		err = app.writeJSON(w, http.StatusOK, envelope{"RareCoin": RareCoin}, nil)
		if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		}
		func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
			id, err := app.readIDParam(r)
			if err != nil {
			// Use the new notFoundResponse() helper.
			app.notFoundResponse(w, r)
			return
			}
			RareCoin := data.RareCoin{
				ID:        id,
				CreatedAt: time.Now(),
				Title:     "Шиллинг новая Ангия",
				Genres:    []string{"England", "12 pence"},
				Price:     414000,
			}
			err = app.writeJSON(w, http.StatusOK, envelope{"RareCoin": RareCoin}, nil)
			if err != nil {
			// Use the new serverErrorResponse() helper.
			app.serverErrorResponse(w, r, err)
			}
			}
			
		}
