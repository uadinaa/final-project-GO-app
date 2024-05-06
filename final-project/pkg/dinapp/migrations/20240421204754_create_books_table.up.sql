CREATE TABLE IF NOT EXISTS books (
    book_id BIGSERIAL PRIMARY KEY,
    book_title TEXT UNIQUE,
    createdAt timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updatedAt timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    description text,
    genre INT REFERENCES genres(genre_id) ON DELETE CASCADE,
    movie_rec BIGINT REFERENCES movies(movie_id) ON DELETE CASCADE
);

-- CREATE TABLE IF NOT EXISTS books (
--     book_id BIGSERIAL PRIMARY KEY,
--     book_title TEXT UNIQUE,
--     createdAt TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
--     updatedAt TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
--     description TEXT,
--     genre INT REFERENCES genres(genre_id) ON DELETE CASCADE,
--     movie_id BIGINT REFERENCES movies_table(movie_id) ON DELETE CASCADE
-- );
