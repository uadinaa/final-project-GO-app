package model

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Movies struct {
	Id               string `json:"id"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	YearOfProduction int    `json:"yearOfProduction"`
	GenreId          string `json:"genreId"`
}

type MovieModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m MovieModel) Insert(movies *Movies) error {
	query := `
		INSERT INTO movies_table (movie_id, movie_title, description, year_of_production, genre_id)
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING movie_id
		`
	args := []interface{}{movies.Title, movies.Description, movies.YearOfProduction, movies.GenreId}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&movies.Id)
}

func (m MovieModel) Get(id int) (*Movies, error) {
	query := `
		SELECT movie_id, movie_title, description, year_of_production, genre_id
		FROM movies_table
		WHERE movie_id = $1
		`
	var movies Movies
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&movies.Id, &movies.Title, &movies.Description, movies.YearOfProduction, movies.GenreId)
	if err != nil {
		return nil, fmt.Errorf("cannot retrive movie with id: %v, %w", id, err)
	}
	return &movies, nil
}

func (m MovieModel) Update(movies *Movies) error {
	query := `
		UPDATE movies_table 
		SET movie_title = $1, description = $2, year_of_production = $3, genre_id = $4
		WHERE movie_id = $5
		`
	args := []interface{}{movies.Title, movies.Description, movies.YearOfProduction, movies.GenreId, movies.Id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// _, err := m.DB.ExecContext(ctx, query, args...)
	// return err
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&movies.Id)
}

func (m MovieModel) Delete(id int) error {
	query := `
		DELETE FROM movies_table
		WHERE id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

// func main() {
// 	db, err := sql.Open("mysql",
// 		"user:password@tcp(127.0.0.1:3306)/hello")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()
// }
//как я поняла сюда и будем пихать дб типа команды для добавления, удаления и тд
