-- Создание таблиц по схемам

CREATE TABLE IF NOT EXISTS books (
    ID SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    issue DATE NOT NULL,
    copies INT NOT NULL
);

CREATE TABLE IF NOT EXISTS readers (
    ID SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    number VARCHAR(50) NOT NULL,
    address VARCHAR(255),
    date_of_birth DATE NOT NULL
);

CREATE TABLE IF NOT EXISTS reader_books (
    book_id INT NOT NULL REFERENCES books(ID) ON DELETE CASCADE,
    reader_id INT NOT NULL REFERENCES readers(ID),
    date_of_rent DATE NOT NULL,
    PRIMARY KEY (book_id, reader_id)
);

-- Добавление уникальных индексов, если их еще нет
CREATE UNIQUE INDEX IF NOT EXISTS unique_book_title ON books(title);
CREATE UNIQUE INDEX IF NOT EXISTS unique_reader_number ON readers(number);