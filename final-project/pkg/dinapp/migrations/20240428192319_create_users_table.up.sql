CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
    users_id bigserial PRIMARY KEY,
    createdAt timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    email text UNIQUE NOT NULL,
    password bytea NOT NULL,
    activated bool NOT NULL,
    token text NOT NULL
);
