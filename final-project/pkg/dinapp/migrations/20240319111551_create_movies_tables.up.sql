create table if not exists genres_table(genre_id serial primary key, genre_title text unique);

create table if not exists movies_table(movie_id BIGSERIAL PRIMARY KEY,  movie_title text, description text, year_of_production int, genre_id INT REFERENCES genres(genre_id));

