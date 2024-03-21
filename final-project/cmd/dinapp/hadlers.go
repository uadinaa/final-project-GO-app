package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"final-project/pkg/dinapp/model"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	// "main.go/pkg/dinapp/model"
)

func (app *application) respondWithError(w http.ResponseWriter, code int, message string) {
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

	movie := &model.Movies{
		Title:            input.Title,
		Description:      input.Description,
		YearOfProduction: input.YearOfProduction,
		GenreId:          input.GenreId,
	}

	err = app.models.Movies.Insert(movie)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, movie)
}

func (app *application) getMoviesHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	param := params["movieId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid movie ID")
		return
	}

	movies, err := app.models.Movies.Get(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Menu with ID %d not found\n", id)
		}
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, movies)
}

func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	param := params["movieId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid movie ID")
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
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
		movie.GenreId = *input.GenreId
	}

	err = app.models.Movies.Update(movie)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	app.respondWithJSON(w, http.StatusOK, movie)
}

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	param := params["moviesId"]

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

func (app *application) readJSON(_ http.ResponseWriter, r *http.Request, dst interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		return err
	}
	return nil
}

func (app *application) createGenresHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Id    string `json:"id"`
		Title string `json:"title"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		log.Println(err)
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	genre := &model.Genres{
		Id:    input.Id,
		Title: input.Title,
	}

	err = app.models.Genres.InsertG(genre)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, genre)
}

func (app *application) getGenresHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	param := params["genreId"]

	genre, err := app.models.Genres.GetG(param)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, genre)
}

func (app *application) updateGenreHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	param := params["genreId"]

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

// func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	param := params["moviesId"]

// 	id, err := strconv.Atoi(param)
// 	if err != nil || id < 1 {
// 		app.respondWithError(w, http.StatusBadRequest, "Invalid movie ID")
// 		return
// 	}

// 	err = app.models.Movies.Delete(id)
// 	if err != nil {
// 		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
// 		return
// 	}
// 	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
//

// "main.go/pkg/dinapp/model"
// "github.com/uadinaa/final-project-GO-app/tree/main/final-project/pkg/dinapp/model"
// "github.com/uadinaa/final-project-GO-app/final-project/pkg/dinapp/model"

// model "github.com/uadinaa/final-project-GO-app/tree/main/model"
// model "command-line-arguments/Users/dinaabitova/code/golan/final-project/pkg/dinapp/model"
// model "command-line-arguments/Users/dinaabitova/code/golan/final-project/pkg/dinapp/model/movies.go"
// model "github.com/uadinaa/final-project-GO-app/tree/main/model"

// func createMusic(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	var music MusicMax
// 	_ = json.NewDecoder(r.Body).Decode(&music)
// 	music.ID = strconv.Itoa(rand.Intn(100)) //по сути нат зе бест чойс просто рандомно просто создает the id
// 	musics = append(musics, music)
// 	json.NewEncoder(w).Encode(music)
// }

// func getMusic(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json") //сделали его стрингом типа соны кайтару керек деп
// 	params := mux.Vars(r)
// 	for _, item := range musics { // item = iterator
// 		if item.ID == params["id"] {
// 			json.NewEncoder(w).Encode(item)
// 			return
// 		}
// 	}
// 	json.NewEncoder(w).Encode(&MusicMax{})
// }

// // to add new song to site
// func createMusic(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	var music MusicMax
// 	_ = json.NewDecoder(r.Body).Decode(&music)
// 	music.ID = strconv.Itoa(rand.Intn(100)) //по сути нат зе бест чойс просто рандомно просто создает the id
// 	musics = append(musics, music)
// 	json.NewEncoder(w).Encode(music)
// }

// // we can обновить инфу
// func updateMusics(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	params := mux.Vars(r)
// 	for index, item := range musics {
// 		if item.ID == params["id"] {
// 			musics = append(musics[:index], musics[index+1:]...)
// 			var music MusicMax
// 			_ = json.NewDecoder(r.Body).Decode(&music)
// 			music.ID = params["id"]
// 			musics = append(musics, music)
// 			json.NewEncoder(w).Encode(music)
// 			return
// 		}
// 	}
// }

// func deleteMusics(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	params := mux.Vars(r)
// 	for index, item := range musics {
// 		if item.ID == params["id"] {
// 			musics = append(musics[:index], musics[index+1:]...)
// 			break
// 		}
// 	}
// 	json.NewEncoder(w).Encode(musics)
// }
