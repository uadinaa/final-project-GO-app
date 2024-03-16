package model

import (
	"database/sql"
	"errors"
	"log"
)

type Genres struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

type GenreModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

var genres = []Genres{
	{
		Id:    "5",
		Title: "fantasy",
	},
	{
		Id:    "6",
		Title: "classic",
	},
}

func GetGenres() []Genres {
	return genres
}

func GetGenre(id string) (*Genres, error) {
	for _, g := range genres {
		if g.Id == id {
			return &g, nil
		}
	}
	return nil, errors.New("genre was not found")
}

func DeleteGenre(id string) (*Genres, error) {
	for i, g := range genres {
		if g.Id == id {
			if i != -1 {
				genres = append(genres[:i], genres[i+1:]...)
				return nil, nil
			}
		}
	}
	return nil, errors.New("genre was not found")
}
