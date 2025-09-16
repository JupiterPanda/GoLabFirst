package app

import (
	"context"
	"fmt"
	"goproject/internal/package/migrator"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	migrationsPath = "migrations" // Путь к папке с миграциями
)

// Run инициализирует подключение к базе, применяет миграции и запускает приложение
func Run() {
	ctx := context.Background()

	// Подключаемся к базе данных через пул соединений
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Printf("Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Запускаем миграции
	err = migrator.Migrate(ctx, pool, migrationsPath)
	if err != nil {
		log.Printf("Migration failed: %v", err) // Завершаем, если миграции не применились
	}

	fmt.Println("DB migrated")

	// Здесь можно инициализировать Gin и дальше работать с приложением через pool
	// ...
}
