package app

import (
	"context"
	"goproject/internal/delivery/libgrpc"
	constants "goproject/internal/package"
	"goproject/internal/package/migrator"
	"goproject/protos/gen"
	"net"
	"net/http"

	"log"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Run инициализирует подключение к базе, применяет миграции и запускает приложение
func Run() {
	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Подключаемся к базе данных через пул соединений
	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer pool.Close()

	// Запускаем миграции
	// TODO: переехать на goose
	err = migrator.Migrate(ctx, pool, constants.MigrationsPath)
	if err != nil {
		log.Fatalf("Migration failed: %v", err) // Завершаем, если миграции не применились
	}

	go func() {
		mux := runtime.NewServeMux()

		// Зарегистрировать HTTP-ручки для сервиса Library
		if err := gen.RegisterLibraryHandlerFromEndpoint(
			ctx,
			mux,
			os.Getenv("GPRC_SERVER_ENDPOINT"),
			[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
		); err != nil {
			log.Fatalf("failed to start HTTP gateway: %v", err)
		}

		log.Println("HTTP gateway listening on ", os.Getenv("HTTP_SERVER_PORT"))
		if err := http.ListenAndServe(os.Getenv("HTTP_SERVER_PORT"), mux); err != nil {
			log.Fatalf("failed to serve HTTP gateway: %v", err)
		}
	}()

	lis, err := net.Listen("tcp", os.Getenv("GPRC_SERVER_PORT"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	gen.RegisterLibraryServer(grpcServer, libgrpc.NewGRPCServer(initUseCase(pool)))
	log.Println("gRPC server listening on ", os.Getenv("GPRC_SERVER_PORT"))
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
