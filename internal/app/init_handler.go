package app

import (
	"goproject/internal/handlers"
	booksRepoPackage "goproject/internal/repositories/books"
	booksInUseRepoPackage "goproject/internal/repositories/booksInUse"
	readersRepoPackage "goproject/internal/repositories/readers"
	booksServicePackage "goproject/internal/services/books"
	booksInUseServicePackage "goproject/internal/services/booksInUse"
	readersServicePackage "goproject/internal/services/readers"
	"goproject/internal/usecases"

	"github.com/jackc/pgx/v5/pgxpool"
)

func initHandler(pool *pgxpool.Pool) *handlers.Handler {
	// Репозитории
	bookRepo := booksRepoPackage.NewRepo(pool)
	readerRepo := readersRepoPackage.NewRepo(pool)
	bookInUseRepo := booksInUseRepoPackage.NewRepo(pool)

	// Сервисы
	bookService := booksServicePackage.NewService(bookRepo)
	readerService := readersServicePackage.NewService(readerRepo)
	bookInUseService := booksInUseServicePackage.NewService(bookInUseRepo)

	useCase := usecases.NewUseCase(bookService, readerService, bookInUseService)
	return handlers.NewHandler(useCase)
}
