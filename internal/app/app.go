package app

import (
	"context"
	//handlers "goproject/internal/handler/handler"
	constants "goproject/internal/package"
	"goproject/internal/package/migrator"
	booksRepoPackage "goproject/internal/repositories/books"
	booksInUseRepoPackage "goproject/internal/repositories/booksinuse"
	readersRepoPackage "goproject/internal/repositories/readers"
	booksServicePackage "goproject/internal/services/books"
	booksInUseServicePackage "goproject/internal/services/booksinuse"
	readersServicePackage "goproject/internal/services/readers"
	"log"
	"os"

	//"github.com/gin-gonic/gin"
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

	// Хендлеры
	//bookHandler := handlers.NewBookHandler(bookService)
	//readerHandler := handlers.NewReaderHandler(readerService)
	//bookInUseHandler := handlers.NewBookInUseHandler(bookInUseService)

	// Роутер и маршруты
	//router := gin.Default()
	// router.GET()/POST() и т.д. — регистрация хендлеров

	//router.Run()
}
