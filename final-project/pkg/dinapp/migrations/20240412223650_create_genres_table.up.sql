CREATE TABLE IF NOT EXISTS genres (
    genre_id SERIAL PRIMARY KEY,
    genre_title TEXT,
    createdAt timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updatedAt timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS movies (
    movie_id BIGSERIAL PRIMARY KEY,
    movie_title TEXT,
    createdAt timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updatedAt timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    description TEXT,
    year_of_production INT,
    genre  TEXT
    --genre_id INT REFERENCES genres(genre_id) ON DELETE CASCADE
);

-- CREATE TABLE IF NOT EXISTS movies_genres (
--     id SERIAL PRIMARY KEY,
--     genre_id INT REFERENCES genres(genre_id) ON DELETE CASCADE,
--     movie_id INT REFERENCES movies(movie_id) ON DELETE CASCADE
-- );