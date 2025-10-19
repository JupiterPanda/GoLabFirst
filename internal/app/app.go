package app

import (
	"context"
	constants "goproject/internal/package"
	"goproject/internal/package/migrator"

	"log"
	"os"

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
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer pool.Close()

	// Запускаем миграции
	err = migrator.Migrate(ctx, pool, constants.MigrationsPath)
	if err != nil {
		log.Fatalf("Migration failed: %v", err) // Завершаем, если миграции не применились
	}

	// Инициируем handler
	handler := initHandler(pool)
	log.Println("Все типы и бд проинициализированы")

	// Инициируем роутер
	router := initRouter(handler)
	err = router.Run("localhost:8080")
	if err != nil {
		return
	}
}
