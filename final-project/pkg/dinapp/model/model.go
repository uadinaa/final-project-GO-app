package model

import (
	"database/sql"
	"log"
	"os"
)

type Models struct {
	Movies MovieModel
	Genres GenreModel
}

func NewModels(db *sql.DB) Models {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	return Models{
		Genres: GenreModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Movies: MovieModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
	}
}

// type Genres struct {
// 	Id    string `json:"id"`
// 	Title string `json:"title"`
// }

// type Movies struct {
// 	Id               string `json:"id"`
// 	Title            string `json:"title"`
// 	Description      string `json:"description"`
// 	YearOfProduction int    `json:"yearOfProduction"`
// 	GenreId          string `json:"genreId"`
// }

// type GenreModel struct {
// 	DB       *sql.DB
// 	InfoLog  *log.Logger
// 	ErrorLog *log.Logger
// }

// type MovieModel struct {
// 	DB       *sql.DB
// 	InfoLog  *log.Logger
// 	ErrorLog *log.Logger
// }
