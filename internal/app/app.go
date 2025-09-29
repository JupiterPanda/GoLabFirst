package app

import (
	"context"
	"fmt"

	"goproject/internal/handlers"
	constants "goproject/internal/package"
	"goproject/internal/package/migrator"
	booksRepoPackage "goproject/internal/repositories/books"
	booksInUseRepoPackage "goproject/internal/repositories/booksinuse"
	readersRepoPackage "goproject/internal/repositories/readers"
	booksServicePackage "goproject/internal/services/books"
	booksInUseServicePackage "goproject/internal/services/booksinuse"
	readersServicePackage "goproject/internal/services/readers"
	"goproject/internal/usecases"
	"log"
	"os"

	//"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// Run инициализирует подключение к базе, применяет миграции и запускает приложение
func Run() {
	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Подключаемся к базе данных через пул соединений
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Printf("Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	defer pool.Close()

	// Запускаем миграции
	err = migrator.Migrate(ctx, pool, constants.MigrationsPath)
	if err != nil {
		log.Printf("Migration failed: %v", err) // Завершаем, если миграции не применились
	}

	// Репозитории
	bookRepo := booksRepoPackage.NewRepo(pool)
	readerRepo := readersRepoPackage.NewRepo(pool)
	bookInUseRepo := booksInUseRepoPackage.NewRepo(pool)

	// Сервисы
	bookService := booksServicePackage.NewService(bookRepo)
	readerService := readersServicePackage.NewService(readerRepo)
	bookInUseService := booksInUseServicePackage.NewService(bookInUseRepo)

	useCase := usecases.NewUseCase(bookService, readerService, bookInUseService)
	handler := handlers.NewHandler(useCase)

	fmt.Println("Все типы проинициализированы", handler)

	// Роутер и маршруты
	router := gin.Default()

	// Чтение, аренда, возврат
	router.GET("/reader/books", handler.GetReaderBooksSepGoodAndBad)
	router.PATCH("/rent", handler.RentBookByTitleAndReaderName)
	router.PATCH("/return", handler.ReturnBookByTitleAndReaderName)

	// Книги CRUD
	router.GET("/books", handler.GetAllBooks)
	router.GET("/book/title", handler.GetBookByTitle)
	router.GET("/book/id/title", handler.GetBookIdByTitle)
	router.GET("/book/id", handler.GetBookByID)
	router.POST("/book", handler.CreateBook)
	router.DELETE("/book", handler.DeleteBook)

	// Копии книги
	router.POST("/book/check/id", handler.CheckCopiesOfBookByID)
	router.POST("/book/check", handler.CheckCopiesOfBook)
	router.PATCH("/book/minus", handler.MinusCopyOfBookById)
	router.PATCH("/book/plus", handler.PlusCopyOfBookById)

	// Операции с InUse
	router.POST("/reader/book", handler.CreateBookInUse)
	router.GET("/reader/book", handler.GetAllBooksInUse)
	router.GET("/reader/book/count", handler.CountBookInUseByReaderId)
	router.GET("/reader/book/id", handler.GetBooksInUseByReaderId)
	router.DELETE("/reader/book", handler.DeleteBookInUse)

	// CRUD для читателей
	router.GET("/readers", handler.GetAllReaders)
	router.GET("/reader/id", handler.GetReaderIdByName)
	router.POST("/reader", handler.CreateReader)
	router.DELETE("/reader", handler.DeleteReader)
	router.PATCH("/reader/contact", handler.UpdateReaderContactInfo)
	router.GET("/reader/id/book", handler.GetReadersIdsByBookId)

	err = router.Run("localhost:8080")
	if err != nil {
		return
	}
}
