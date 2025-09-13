-- Таблица книг
CREATE TABLE IF NOT EXISTS books (
    ID SERIAL PRIMARY KEY,
    Title VARCHAR(255) NOT NULL,
    Author VARCHAR(255) NOT NULL,
    Issue DATE NOT NULL,
    Copies INT NOT NULL
);

-- Таблица читателей
CREATE TABLE IF NOT EXISTS readers (
    ID SERIAL PRIMARY KEY,
    Name VARCHAR(255) NOT NULL,
    Number VARCHAR(50) NOT NULL,
    Address VARCHAR(255),
    DateOfBirth DATE NOT NULL
);

-- Связующая таблица: книги в пользовании
CREATE TABLE IF NOT EXISTS readerBooks (
    Book_ID INT NOT NULL REFERENCES Book(ID) ON DELETE CASCADE,
    Reader_ID INT NOT NULL REFERENCES Reader(ID),
    DateOfRent DATE NOT NULL,
    PRIMARY KEY (Book_ID, Reader_ID, DateOfRent)
);