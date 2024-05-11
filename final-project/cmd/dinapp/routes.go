package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {

	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)
	r.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)
	r.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)

	// r.HandleFunc("/api/v1/healthcheck", app.healthcheckHandler).Methods("GET")

	v1 := r.PathPrefix("/api/v1").Subrouter()

	v1.HandleFunc("/healthcheck", app.healthcheckHandler).Methods("GET")

	v1.HandleFunc("/movie/{id}", app.getMoviesHandler).Methods("GET")
	v1.HandleFunc("/movie", app.createMoviesHandler).Methods("POST")
	v1.HandleFunc("/movie/{movieId}", app.updateMovieHandler).Methods("PUT")
	v1.HandleFunc("/movie/{movieId}", app.deleteMovieHandler).Methods("DELETE")

	v1.HandleFunc("/movies", app.listMoviesHandler).Methods("GET")

	v1.HandleFunc("/genre/{genreId}", app.getGenresHandler).Methods("GET")
	v1.HandleFunc("/genre", app.createGenresHandler).Methods("POST")
	v1.HandleFunc("/genre/{genreId}", app.updateGenreHandler).Methods("PUT")
	v1.HandleFunc("/genre/{genreId}", app.deleteGenreHandler).Methods("DELETE")

	v1.HandleFunc("/users", app.registerUserHandler).Methods("POST")
	v1.HandleFunc("/users/activated", app.activateUserHandler).Methods("PUT")

	return r
	// app.recoverPanic(app.rateLimit(r))
}
