-- +goose Up
-- +goose StatementBegin
-- A single table to track the watchlist
CREATE TABLE Watchlist (
    watchlist_id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL UNIQUE,
    release_year INTEGER,
    genre TEXT,
    director TEXT,
    status TEXT CHECK(status IN ('watched', 'not watched', 'watching')) DEFAULT 'not watched',
    added_date DATE DEFAULT (date('now'))
);

-- Sample Data
INSERT INTO Watchlist (title, release_year, genre, director, status)
VALUES ('Inception', 2010, 'Science Fiction', 'Christopher Nolan', 'watched');

INSERT INTO Watchlist (title, release_year, genre, director, status)
VALUES ('The Matrix', 1999, 'Action', 'Lana Wachowski, Lilly Wachowski', 'watching');

INSERT INTO Watchlist (title, release_year, genre, director, status)
VALUES ('The Godfather', 1972, 'Crime', 'Francis Ford Coppola', 'not watched');

INSERT INTO Watchlist (title, release_year, genre, director, status)
VALUES ('Parasite', 2019, 'Thriller', 'Bong Joon Ho', 'watched');

INSERT INTO Watchlist (title, release_year, genre, director, status)
VALUES ('Avengers: Endgame', 2019, 'Action', 'Anthony Russo, Joe Russo', 'not watched');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
