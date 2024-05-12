package model

import (
	"context"
	"github.com/jmoiron/sqlx"
	"fmt"
	"log"
	"time"


)

type Movies struct {
	Id               string `json:"id"`
	CreatedAt        string `json:"createdAt"`
	UpdatedAt        string `json:"updatedAt"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	YearOfProduction int    `json:"year_of_production"`
	Genre          string `json:"genre"`
}

type MovieModel struct {
	DB       *sqlx.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m MovieModel) Insert(movies *Movies) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []interface{}{movies.Genre}

	var id int 
	query := `
		INSERT INTO genres (genre_title)
		VALUES ($1) 
		RETURNING genre_id
		`
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return err
	}
	
	query = `
		INSERT INTO movies (movie_title, description, year_of_production, genre)
		VALUES ($1, $2, $3, $4) 
		RETURNING movie_id
		`
	args = []interface{}{movies.Title, movies.Description, movies.YearOfProduction, movies.Genre}


	return m.DB.QueryRowContext(ctx, query, args...).Scan(&id)
}

func (m MovieModel) Get(id int) (*Movies, error) {
	query := `
		SELECT * 
		FROM movies
		WHERE movie_id = $1
		`
	//movie_id, movie_title, description, year_of_production, genre_id

	var movie Movies

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&movie.Id, &movie.Title, &movie.CreatedAt, &movie.UpdatedAt, &movie.Description, &movie.YearOfProduction, &movie.Genre)

	if err != nil {
		return nil, fmt.Errorf("cannot retrive movie with id: %v, %w", id, err)
	}

	return &movie, nil
}

func (m MovieModel) GetAll(title string, yearOfProduction int, genreId string, filters Filters) ([]*Movies, Metadata, error) { //сюда надо добавить фильтра или сортка данные, убрала китаптыкин жанры

	query := fmt.Sprintf(`
		SELECT count(*) OVER(), movie_id, createdat, updatedat, movie_title, description, year_of_production, genre_id
		FROM  movies
		WHERE (LOWER(title) = LOWER($4) OR $4 = '') 
		AND (year_of_production >= $6 OR $6 = 0)
		AND (LOWER(genre_id) = LOWER($7) OR $7 = '')
		ORDER BY %s %s, movie_id ASC
		LIMIT $8 OFFSET $9`, filters.sortColumn(), filters.sortDirection())
	//to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = ''

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// args := []interface{}{title, yearOfProduction, genreId, filters.limit(), filters.offset()}

	rows, err := m.DB.QueryContext(ctx, query, title, yearOfProduction, genreId, filters.limit(), filters.offset())

	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	movies := []*Movies{}

	for rows.Next() {
		var movie Movies

		err := rows.Scan(
			&totalRecords,
			&movie.Id,
			&movie.CreatedAt,
			&movie.UpdatedAt,
			&movie.Title,
			&movie.Description,
			&movie.YearOfProduction,
			&movie.Genre,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		movies = append(movies, &movie)
	}

	if err := rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return movies, metadata, nil
}

func (m MovieModel) Update(movie *Movies) error {
	query := `
		UPDATE movies
		SET movie_title = $1, description = $2, year_of_production = $3, genre_id = $4
		WHERE movie_id = $5
		RETURNING updatedAt
		`

	args := []interface{}{movie.Title, movie.Description, movie.YearOfProduction, movie.Genre, movie.Id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&movie.UpdatedAt)
}

func (m MovieModel) Delete(id int) error {
	query := `
		DELETE FROM movies
		WHERE movie_id = $1
		`

	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	// _, err := m.DB.ExecContext(ctx, query, id)
	// return err
}
