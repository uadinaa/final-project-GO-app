package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
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

func (gm GenreModel) InsertG(genre *Genres) error {
	query := `
		INSERT INTO genres_table (genre_id, genre_title)
		VALUES ($1, $2)
		RETURNING genre_id
	`
	row := gm.DB.QueryRowContext(context.Background(), query, genre.Id, genre.Title)

	if err := row.Scan(&genre.Id); err != nil {
		// Handle any error that occurred during scanning
		return err
	}
	return nil
}

func (gm GenreModel) GetG(id string) (*Genres, error) {
	query := `
		SELECT genre_id, genre_title
		FROM genres_table
		WHERE genre_id = $1
		`
	var genre Genres
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := gm.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&genre.Id, &genre.Title)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve genre with id: %v, %w", id, err)
	}
	return &genre, nil
}

func (gm GenreModel) UpdateG(genre *Genres) error {
	query := `
		UPDATE genres_table 
		SET genre_title = $1
		WHERE genre_id = $2
		`
	args := []interface{}{genre.Title, genre.Id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := gm.DB.ExecContext(ctx, query, args...)
	return err
}

func (gm GenreModel) DeleteG(id string) error {
	query := `
		DELETE FROM genres_table
		WHERE genre_id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := gm.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("genre not found")
	}
	return nil
}

// var genres = []Genres{
// 	{
// 		Id:    "5",
// 		Title: "fantasy",
// 	},
// 	{
// 		Id:    "6",
// 		Title: "classic",
// 	},
// }

// func GetGenres() []Genres {
// 	return genres
// }

// func GetGenre(id string) (*Genres, error) {
// 	for _, g := range genres {
// 		if g.Id == id {
// 			return &g, nil
// 		}
// 	}
// 	return nil, errors.New("genre was not found")
// }

// func DeleteGenre(id string) (*Genres, error) {
// 	for i, g := range genres {
// 		if g.Id == id {
// 			if i != -1 {
// 				genres = append(genres[:i], genres[i+1:]...)
// 				return nil, nil
// 			}
// 		}
// 	}
// 	return nil, errors.New("genre was not found")
// }
