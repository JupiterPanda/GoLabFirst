package app

import (
	"context"
	"goproject/internal/delivery/libgrpc"
	constants "goproject/internal/package"
	"goproject/internal/package/migrator"
	"goproject/protos/gen/librarypb"
	"net"

	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
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

	/*	// Инициируем handler
		useCase := initUseCase(pool)
		handler := handlers.NewHandler(useCase)
		log.Println("Все структуры, типы и бд проинициализированы")

		// Инициируем роутер
		router := initRouter(handler)
		err = router.Run("localhost:8080")
		if err != nil {
			return
		}*/

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	librarypb.RegisterLibraryServer(grpcServer, libgrpc.NewGRPCServer(initUseCase(pool)))

	log.Println("gRPC server listening on :8080")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
