package main
import (
"fmt"
"net/http"
"time"
"greenlight/internal/data"
"greenlight/internal/validator" // New import
)
func (app *application) createRareCoinsHandler(w http.ResponseWriter, r *http.Request) {
var input struct {
Title string `json:"title"`
Year int32 `json:"year"`
Runtime data.Runtime `json:"runtime"`
Genres []string `json:"genres"`
}
err := app.readJSON(w, r, &input)
if err != nil {
app.badRequestResponse(w, r, err)
return
}
// Initialize a new Validator instance.
v := validator.New()
// Use the Check() method to execute our validation checks. This will add the
// provided key and error message to the errors map if the check does not evaluate
// to true. For example, in the first line here we "check that the title is not
// equal to the empty string". In the second, we "check that the length of the title
// is less than or equal to 500 bytes" and so on.
v.Check(input.Title != "", "title", "must be provided")
v.Check(len(input.Title) <= 500, "title", "must not be more than 500 bytes long")
v.Check(input.Year != 0, "year", "must be provided")
v.Check(input.Year >= 1888, "year", "must be greater than 1888")
v.Check(input.Year <= int32(time.Now().Year()), "year", "must not be in the future")
v.Check(input.Runtime != 0, "runtime", "must be provided")
v.Check(input.Runtime > 0, "runtime", "must be a positive integer")
v.Check(input.Genres != nil, "genres", "must be provided")
v.Check(len(input.Genres) >= 1, "genres", "must contain at least 1 genre")
v.Check(len(input.Genres) <= 5, "genres", "must not contain more than 5 genres")
// Note that we're using the Unique helper in the line below to check that all
// values in the input.Genres slice are unique.
v.Check(validator.Unique(input.Genres), "genres", "must not contain duplicate values")
// Use the Valid() method to see if any of the checks failed. If they did, then use
// the failedValidationResponse() helper to send a response to the client, passing
// in the v.Errors map.
if !v.Valid() {
app.failedValidationResponse(w, r, v.Errors)
return
}
fmt.Fprintf(w, "%+v\n", input)
}

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
			func (app *application) createRareCoinHandler(w http.ResponseWriter, r *http.Request) {
				// Declare an anonymous struct to hold the information that we expect to be in the
				// HTTP request body (note that the field names and types in the struct are a subset
				// of the Movie struct that we created earlier). This struct will be our *target
				// decode destination*.
				var input struct {
				Title string `json:"title"`
				Year int32 `json:"year"`
				Genres []string `json:"genres"`
				}
				// Initialize a new json.Decoder instance which reads from the request body, and
				// then use the Decode() method to decode the body contents into the input struct.
				// Importantly, notice that when we call Decode() we pass a *pointer* to the input
				// struct as the target decode destination. If there was an error during decoding,
				// we also use our generic errorResponse() helper to send the client a 400 Bad
				// Request response containing the error message.
				err := json.NewDecoder(r.Body).Decode(&input)
				if err != nil {
				app.errorResponse(w, r, http.StatusBadRequest, err.Error())
				return
				}
				// Dump the contents of the input struct in a HTTP response.
				fmt.Fprintf(w, "%+v\n", input)
				}
				func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
					var input struct {
					Title string `json:"title"`
					Year int32 `json:"year"`
					Genres []string `json:"genres"`
					}
					err := app.readJSON(w, r, &input)
					if err != nil {
					// Use the new badRequestResponse() helper.
					app.badRequestResponse(w, r, err)
					return
					}
					fmt.Fprintf(w, "%+v\n", input)
					}
					func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
						var input struct {
						Title string `json:"title"`
						Year int32 `json:"year"`
						Runtime data.Runtime `json:"runtime"`
						Genres []string `json:"genres"`
						}
						err := app.readJSON(w, r, &input)
						if err != nil {
						app.badRequestResponse(w, r, err)
						return
						}
						fmt.Fprintf(w, "%+v\n", input)
						}
						func (app *application) createRareCoinsHandler(w http.ResponseWriter, r *http.Request) {
							var input struct {
							Title string `json:"title"`
							Year int32 `json:"year"`
							Genres []string `json:"genres"`
							}
							err := app.readJSON(w, r, &input)
							if err != nil {
							app.badRequestResponse(w, r, err)
							return
							}
							// Copy the values from the input struct to a new Movie struct.
							rarecoin := &data.RareCoin{
							Title: input.Title,
							Year: input.Year,
							Genres: input.Genres,
							}
							// Initialize a new Validator.
							v := validator.New()
							// Call the ValidateMovie() function and return a response containing the errors if
							// any of the checks fail.
							if data.ValidateRareCoin(v, rarecoin); !v.Valid() {
							app.failedValidationResponse(w, r, v.Errors)
							return
							}
							fmt.Fprintf(w, "%+v\n", input)
							}
		}
