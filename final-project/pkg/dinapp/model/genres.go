package model

import (
	"context"
	"github.com/jmoiron/sqlx"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type Genres struct {
	Id        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Title     string `json:"title"`
}

type GenreModel struct {
	DB       *sqlx.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (gm GenreModel) InsertG(genre *Genres) error {
	query := `
		INSERT INTO genres(genre_title)
		VALUES ($1)
		RETURNING genre_id, createdat, updatedat
	`
	//, $2

	args := []interface{}{genre.Title}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := gm.DB.QueryRowContext(ctx, query, args...).Scan(&genre.Id, &genre.Title, &genre.CreatedAt, &genre.UpdatedAt)

	return row
	// if err := row.Scan(&genre.Id); err != nil {
	// 	// Handle any error that occurred during scanning
	// 	return err
	// }
	// return nil
}

func (gm GenreModel) GetG(id string) (*Genres, error) {
	query := `
		SELECT genre_id, genre_title
		FROM genres
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
		UPDATE genres 
		SET genre_title = $1
		WHERE genre_id = $2
		RETURNING updatedat
		`
	args := []interface{}{genre.Title, genre.Id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return gm.DB.QueryRowContext(ctx, query, args...).Scan(&genre.UpdatedAt) //ExecContext was before Query
}

func (gm GenreModel) DeleteG(id string) error {
	query := `
		DELETE FROM genres
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

	// _, err := u.DB.ExecContext(ctx, query, id)
	// return err
}
