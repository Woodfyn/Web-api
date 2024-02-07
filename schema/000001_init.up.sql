CREATE TABLE games
(
    id serial PRIMARY KEY,
    title varchar(255) NOT NULL,
    genre varchar(255) NOT NULL,
    evaluation int NOT NULL
);

CREATE TABLE users
(
    id serial PRIMARY KEY,
    name varchar(255) NOT NULL,
    email varchar(255) NOT NULL,
    password varchar(255) NOT NULL,
    registered_at timestamp NOT NULL
);

CREATE TABLE refresh_tokens
(
    id serial PRIMARY KEY,
    user_id int NOT NULL,
    token varchar(255) NOT NULL,
    expires_at timestamp NOT NULL
);