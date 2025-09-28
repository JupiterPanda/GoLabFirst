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
	// router.GET("/readers", getReaders)
	// router.POST("/reader", postReader)
	router.GET("/books", handler.GetAllBooks)
	// router.POST("/book", postBook)
	// router.GET("/book", getBookByTitle)
	router.GET("/reader/books", handler.GetReaderBooksSepGoodAndBad)
	router.PATCH("/rent", handler.RentBookByTitleAndReaderName)
	router.PATCH("/return", handler.ReturnBookByTitleAndReaderName)

	/* TODO bookinuse routes
	Create(ctx context.Context, bookInUse *models.BookInUse, readerId int) error
	Delete(ctx context.Context, readerId int, bookId int) error
	GetAll(ctx context.Context) ([]models.BookInUse, error)
	GetReadersIdsByBookId(ctx context.Context, bookId int) ([]int, error)
	GetBooksInUseByReaderId(ctx context.Context, readerId int) (map[int]time.Time, error)
	CountByReaderId(ctx context.Context, readerId int) (int, error) */

	/* TODO reader routes
	GetAll(ctx context.Context) ([]models.Reader, error)
	Create(ctx context.Context, reader *models.Reader) error
	GetIdByName(ctx context.Context, name string) (int, error)
	Delete(ctx context.Context, readerId int) error
	UpdateContactInfo(ctx context.Context, readerId int, phoneNumber string, address string) error */

	/* TODO book routes
	GetByTitle(ctx context.Context, title string) (*models.Book, error)
	GetAll(ctx context.Context) ([]models.Book, error)
	GetByID(ctx context.Context, id int) (*models.Book, error)
	GetIdByTitle(ctx context.Context, title string) (int, error)
	CheckCopiesByID(ctx context.Context, id int) error
	CheckCopies(ctx context.Context, book *models.Book) error
	Create(ctx context.Context, book *models.Book) error
	Delete(ctx context.Context, book *models.Book) error
	PlusCopyById(ctx context.Context, id int) error
	MinusCopyById(ctx context.Context, id int) error */

	err = router.Run("localhost:8080")
	if err != nil {
		return
	}
}
