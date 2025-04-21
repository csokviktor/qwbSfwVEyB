-- migrate:up
-- Enable foreign key support (important for SQLite)
PRAGMA foreign_keys = ON;

-- Authors
CREATE TABLE authors (
    id TEXT PRIMARY KEY NOT NULL UNIQUE,
    name TEXT NOT NULL
);

-- Borrowers
CREATE TABLE borrowers (
    id TEXT PRIMARY KEY NOT NULL UNIQUE,
    name TEXT NOT NULL
);

-- Books
CREATE TABLE books (
    id TEXT PRIMARY KEY NOT NULL UNIQUE,
    title TEXT NOT NULL,
    author_id TEXT NOT NULL,
    borrower_id TEXT, -- NULL = available

    FOREIGN KEY (author_id) REFERENCES authors(id) ON DELETE CASCADE,
    FOREIGN KEY (borrower_id) REFERENCES borrowers(id) ON DELETE SET NULL
);

-- migrate:down
DROP TABLE IF EXISTS authors;
DROP TABLE IF EXISTS borrowers;
DROP TABLE IF EXISTS books;