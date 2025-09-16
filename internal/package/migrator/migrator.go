package migrator

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Migrate запускает поочередно файлы из директории migrationsPath
func Migrate(ctx context.Context, db *pgxpool.Pool, migrationsPath string) error {
	// Считываем список файлов в папке миграций
	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		return err
	}

	// Проходим по всем файлам и для каждого применяем миграцию
	for _, file := range files {
		if file.IsDir() { // Пропускаем директории
			continue
		}
		filepath := migrationsPath + "/" + file.Name()

		migrationBytes, err := os.ReadFile(filepath)
		if err != nil {
			return err
		}
		migration := string(migrationBytes)

		log.Printf("Applying migration: %s\n", file.Name())

		// Выполняем SQL запрос миграции к базе данных
		_, err = db.Exec(ctx, migration)
		if err != nil {
			return err
		}
	}
	return nil
}
