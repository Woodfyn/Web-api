CREATE TABLE game
(
    id serial PRIMARY KEY,
    title varchar(255) NOT NULL,
    genre varchar(255) NOT NULL,
    evaluation int NOT NULL
);