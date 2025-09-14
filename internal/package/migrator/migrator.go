package migrator

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Migrate применяет все миграционные SQL-файлы из заданной папки к базе данных.
// Аргументы:
//
//	ctx - контекст для отмены и таймаутов операций
//	db - пул соединений с базой данных PostgreSQL
//	migrationsPath - путь к папке с миграциями (*.sql)
//
// Возвращает ошибку в случае неудачи применения миграций.
func Migrate(ctx context.Context, db *pgxpool.Pool, migrationsPath string) error {
	// Считываем список файлов в папке миграций
	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		return err // Ошибка чтения папки
	}

	// Проходим по всем файлам и для каждого применяем миграцию
	for _, file := range files {
		if file.IsDir() { // Пропускаем директории
			continue
		}
		filepath := migrationsPath + "/" + file.Name()

		// Читаем содержимое миграционного файла
		migrationBytes, err := os.ReadFile(filepath)
		if err != nil {
			return err // Ошибка чтения файла
		}
		migration := string(migrationBytes)

		// Логируем информацию о применяемом файле миграции
		log.Printf("Applying migration: %s\n", file.Name())

		// Выполняем SQL запрос миграции к базе данных
		_, err = db.Exec(ctx, migration)
		if err != nil {
			return err // Ошибка применения миграции
		}
	}
	return nil // Все миграции успешно применены
}
