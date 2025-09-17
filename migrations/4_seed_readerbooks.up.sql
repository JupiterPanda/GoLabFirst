INSERT INTO reader_books (book_id, reader_id, date_of_rent) VALUES
    (1, 1, '2025-09-01'),
    (2, 2, '2025-09-05'),
    (3, 3, '2025-09-10')
ON CONFLICT (book_id, reader_id) DO NOTHING;
