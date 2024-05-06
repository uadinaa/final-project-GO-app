CREATE TABLE IF NOT EXISTS movies (
    movie_id BIGSERIAL PRIMARY KEY,
    movie_title TEXT,
    createdAt timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updatedAt timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    description TEXT,
    year_of_production INT,
    genre_id INT REFERENCES genres(genre_id) ON DELETE CASCADE
);
