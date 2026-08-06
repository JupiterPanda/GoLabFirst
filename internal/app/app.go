package app

import (
	"context"
	"goproject/internal/delivery/libgrpc"
	"goproject/internal/kafka"
	constants "goproject/internal/package"
	"goproject/internal/package/migrator"
	"goproject/protos/gen"
	"net"
	"net/http"
	"strings"

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
	useCases := initUseCase(pool)

	// Запускаем миграции
	err = migrator.Migrate(ctx, pool, constants.MigrationsPath)
	if err != nil {
		log.Fatalf("Migration failed: %v", err) // Завершаем, если миграции не применились
	}

	// Запускаем консюмер кафки на прослушивание топика
	consumer, err := kafka.NewConsumer(ctx, kafka.Config{
		Brokers: strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
		Topic:   os.Getenv("KAFKA_TOPIC"),
		GroupID: os.Getenv("KAFKA_GROUP_ID"),
	}, useCases)
	if err != nil {
		log.Fatalf("failed to init kafka consumer: %v", err)
	}
	defer consumer.Close()

	go func() {
		if err := consumer.Run(ctx); err != nil {
			log.Printf("kafka consumer exited with error: %v", err)
		}
	}()

	// Запускаем сервер для обработки http пакетов
	go func() {
		mux := runtime.NewServeMux()

		// Регистрируем HTTP-ручки для сервиса Library
		if err := gen.RegisterLibraryHandlerFromEndpoint(
			ctx,
			mux,
			os.Getenv("GRPC_SERVER_ENDPOINT"),
			[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
		); err != nil {
			log.Fatalf("failed to start HTTP gateway: %v", err)
		}

		log.Println("HTTP gateway listening on ", os.Getenv("HTTP_SERVER_PORT"))
		if err := http.ListenAndServe(os.Getenv("HTTP_SERVER_PORT"), mux); err != nil {
			log.Fatalf("failed to serve HTTP gateway: %v", err)
		}
	}()

	// Запускаем сервер для обработки grpc пакетов
	lis, err := net.Listen("tcp", os.Getenv("GRPC_SERVER_PORT"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	gen.RegisterLibraryServer(grpcServer, libgrpc.NewGRPCServer(useCases))
	log.Println("gRPC server listening on ", os.Getenv("GRPC_SERVER_PORT"))
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
