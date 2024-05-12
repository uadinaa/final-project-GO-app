package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"final-project/pkg/dinapp/model"

	"final-project/pkg/dinapp/validator"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	// "main.go/pkg/dinapp/model"
)

func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, headers http.Header) error {

	js, err := json.MarshalIndent(data, "", "\t")

	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) readIDParam(r *http.Request) (int64, error) {
	param := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}

// func (app *application) readUsernameParam(r *http.Request) (string, error) {
// 	username := mux.Vars(r)["username"]
// 	return username, nil
// }

func (app *application) respondWithError(w http.ResponseWriter, code int, message string) {
	// app.logError(r, message)

	app.respondWithJSON(w, code, map[string]string{"error": message})

}

func (app *application) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// func (app *application) readJSON(_ http.ResponseWriter, r *http.Request, dst interface{}) error {
// 	dec := json.NewDecoder(r.Body)
// 	dec.DisallowUnknownFields()

//		err := dec.Decode(dst)
//		if err != nil {
//			return err
//		}
//		return nil
//	}
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {

	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON at (charcter %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q",
					unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)",
				unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		// For anything else, return the error message as-is.
		default:
			return err
		}
	}
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

//
//
//
// MOVIE
//
//
//

func (app *application) createMoviesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title            string `json:"title"`
		Description      string `json:"description"`
		YearOfProduction int    `json:"yearOfProduction"`
		GenreId          string `json:"genreId"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		log.Println(err)
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	log.Println(input)
	movie := &model.Movies{
		Title:            input.Title,
		Description:      input.Description,
		YearOfProduction: input.YearOfProduction,
		Genre:          input.GenreId,
	}

	err = app.models.Movies.Insert(movie)
	if err != nil {
		log.Println(err)
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	headers := make(http.Header)
	// headers.Set("location", fmt.Sprintf("/v1/movies/%d", movie.Id))
	headers.Set("location", fmt.Sprintf("/v1/movie/%v", movie.Id))

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, headers)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	app.respondWithJSON(w, http.StatusCreated, movie)
	app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, headers)
}

// func (app *application) createGenresHandler(w http.ResponseWriter, r *http.Request) {
// 	var input struct {
// 		Title string `json:"title"`
// 	}

// 	err := app.readJSON(w, r, &input)
// 	if err != nil {
// 		log.Println(err)
// 		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
// 		return
// 	}

// 	genre := &model.Genres{
// 		Title: input.Title,
// 	}

// 	err = app.models.Genres.InsertG(genre)
// 	if err != nil {
// 		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
// 		return
// 	}

// 	// app.respondWithJSON(w, http.StatusCreated, genre)
// 	app.writeJSON(w, http.StatusCreated, envelope{"genre": genre}, nil)
// 	// fmt.Fprintln(w, "status: available")
// 	// fmt.Fprintf(w, "enviroment: %s\n", app.config.env)
// }

func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Title            string
		Description      string
		YearOfProduction int
		GenreId          string
		model.Filters
	}

	v := validator.New()

	qs := r.URL.Query()


	input.Title = app.readString(qs, "movie_title", "")

	input.YearOfProduction = app.readInt(qs, "yearOfProduction", 2000, v)

	input.GenreId = app.readString(qs, "genreId", "")

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")

	input.Filters.SortSafelist = []string{
		"movie_id", "movie_title", "description", "genreId", "yearOfProduction",
		"-movie_id", "-movie_title", "-description", "-genreId", "-yearOfProduction"}

	if model.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if model.ValidateFilters(v, input.Filters); !v.Valid() {
		// Convert validation errors into a single error message
		errorMessage := "Validation errors: "
		for _, err := range v.Errors {
			errorMessage += err + "; "
		}
		app.serverErrorResponse(w, r, errors.New(errorMessage))

		return
	}

	// movies, err := app.models.Movies.GetAll(input.Title, input.YearOfProduction, input.Filters) //, input.GenreId

	movies, metadata, err := app.models.Movies.GetAll(input.Title, input.YearOfProduction, input.GenreId, input.Filters) //input.GenreId,

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	app.writeJSON(w, http.StatusOK, envelope{"movies": movies, "metadata": metadata}, nil)

	// err = app.writeJSON(w, http.StatusOK, envelope{"movies": movies}, nil)

	// if err != nil {
	// 	app.serverErrorResponse(w, r, err)
	// }
}

func (app *application) getMoviesHandler(w http.ResponseWriter, r *http.Request) {
	// params := mux.Vars(r)
	// param := params["movieId"]

	// id, err := strconv.Atoi(param)
	id, err := app.readIDParam(r)

	if err != nil { //|| id < 1
		app.respondWithError(w, http.StatusBadRequest, "Invalid movie ID")
		return
	}

	movies, err := app.models.Movies.Get(int(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Movie with ID %d not found\n", id)
			app.respondWithError(w, http.StatusNotFound, "404 Not Found")
			return
		}
		app.serverErrorResponse(w, r, err)
		app.respondWithError(w, http.StatusNotFound, "404pp Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, movies)
}

// func (app *application) getBookHandler(w http.ResponseWriter, r *http.Request) {
// 	id, err := app.readIDParam(r)
// 	if err != nil {
// 		app.notFoundResponse(w, r)
// 		return
// 	}
// 	book, err := app.models.Books.Get(id)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, model.ErrRecordNotFound):
// 			app.notFoundResponse(w, r)
// 		default:
// 			app.serverErrorResponse(w, r, err)
// 		}
// 		return
// 	}
// 	err = app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 	}
// }

func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	param := params["movieId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid movie ID")
		return
	}

	movie, err := app.models.Movies.Get(int(id))
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404uu Not Found")
		return
	}

	var input struct {
		Title            *string `json:"title"`
		Description      *string `json:"description"`
		YearOfProduction *int    `json:"yearOfProduction"`
		GenreId          *string `json:"genreId"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Title != nil {
		movie.Title = *input.Title
	}

	if input.Description != nil {
		movie.Description = *input.Description
	}

	if input.YearOfProduction != nil {
		movie.YearOfProduction = *input.YearOfProduction
	}

	if input.GenreId != nil {
		movie.Genre = *input.GenreId
	}

	// v := validator.New()
	// if data.ValidateMovie(v, movie); !v.Valid() {
	// 	app.failedValidationResponse(w, r, v.Errors)
	// 	return
	// }

	err = app.models.Movies.Update(movie)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	app.respondWithJSON(w, http.StatusOK, movie)
}

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	param := params["movieId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid movie ID")
		return
	}

	err = app.models.Movies.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

//
//
//
// GENRE
//
//
//

func (app *application) createGenresHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string `json:"title"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		log.Println(err)
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	genre := &model.Genres{
		Title: input.Title,
	}

	err = app.models.Genres.InsertG(genre)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	// app.respondWithJSON(w, http.StatusCreated, genre)
	app.writeJSON(w, http.StatusCreated, envelope{"genre": genre}, nil)
	// fmt.Fprintln(w, "status: available")
	// fmt.Fprintf(w, "enviroment: %s\n", app.config.env)
}

func (app *application) getGenresHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	param := params["genreId"]

	// genre, err := app.models.Genres.GetG(param)
	genreID, err := strconv.Atoi(param)
	if err != nil || genreID < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid genre ID")
		return
	}

	genre, err := app.models.Genres.GetG(param)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("genre with ID %d not found\n", genreID)
		}
		app.respondWithError(w, http.StatusNotFound, "404pp Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, genre)

}

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"welcoming message": "Hello! Welcome to Diina's API about movies and books",
			"environment":       app.config.env,
		},
	}
	// time.Sleep(1 * time.Second)
	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) updateGenreHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	param := params["genreId"]

	// genreID, err := strconv.Atoi(param)
	genre, err := app.models.Genres.GetG(param)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	var input struct {
		Title string `json:"title"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	genre.Title = input.Title

	err = app.models.Genres.UpdateG(genre)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	app.respondWithJSON(w, http.StatusOK, genre)
}

func (app *application) deleteGenreHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	param := params["genreId"]
	// genreID, err := strconv.Atoi(param)

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid movie ID")
		return
	}
	// err = app.models.Movies.Delete(id)

	err = app.models.Genres.DeleteG(param)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
