package migrator

import (
	"context"
	"database/sql"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

////go:embed migrations/*.sql
//var embedMigrations embed.FS

// Migrate запускает поочередно файлы из директории migrationsPath
func Migrate(ctx context.Context, dbpool *pgxpool.Pool, migrationsPath string) error {

	//var embedMigrations embed.FS
	//goose.SetBaseFS(embedMigrations)
	log.Println("starting migrations")
	if err := goose.SetDialect(string(goose.DialectPostgres)); err != nil {
		panic(err)
	}

	db := stdlib.OpenDBFromPool(dbpool)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Was NOT closed database connection: %v", err)
		}
	}(db)
	if err := goose.UpContext(ctx, db, migrationsPath); err != nil {
		panic(err)
	}

	log.Println("DB migrated")
	return nil
}
