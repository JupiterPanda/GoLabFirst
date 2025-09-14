package app

import (
	"context"
	"fmt"
	"goproject/internal/package/migrator"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Run инициализирует подключение к базе, применяет миграции и запускает приложение
func Run() {
	ctx := context.Background() // контекст с отменой и таймаутом можно передать сюда

	// URL подключения к базе PostgreSQL на localhost с параметром sslmode=disable
	dbUrl := "postgresql://postgres:qwerty@localhost:5432/postgres?sslmode=disable"

	// Обработка и разбор конфигурации подключения для pgxpool
	poolConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		log.Printf("Unable to parse database config: %v", err) // Завершаем, если ошибка
	}

	poolConfig.MaxConns = 10 // Максимальное количество соединений в пуле

	// Подключаемся к базе данных через пул соединений
	dbpool, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		log.Printf("Unable to connect to database: %v", err) // Завершаем, если ошибка соединения
	}
	defer dbpool.Close() // Закрываем пул соединений при выходе

	// Путь к папке с миграциями
	migrationsPath := "migrations"

	// Запускаем миграции
	err = migrator.Migrate(ctx, dbpool, migrationsPath)
	if err != nil {
		log.Printf("Migration failed: %v", err) // Завершаем, если миграции не применились
	}

	fmt.Println("DB migrated")

	// Здесь можно инициализировать Gin и дальше работать с приложением через dbpool
	// ...
}
