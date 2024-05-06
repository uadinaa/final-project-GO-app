CREATE TABLE IF NOT EXISTS genres (
    genre_id SERIAL PRIMARY KEY,
    genre_title TEXT UNIQUE,
    createdAt timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updatedAt timestamp(0) with time zone NOT NULL DEFAULT NOW()
);
