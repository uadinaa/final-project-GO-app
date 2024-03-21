CREATE TABLE IF NOT EXISTS genres (
    genre_id SERIAL PRIMARY KEY,
    genre_title TEXT UNIQUE
);

CREATE TABLE IF NOT EXISTS movies (
    movie_id BIGSERIAL PRIMARY KEY,
    movie_title TEXT,
    description TEXT,
    year_of_production INT,
    genre_id INT REFERENCES genres(genre_id)
);