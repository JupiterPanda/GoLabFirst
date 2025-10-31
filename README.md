### Запуск докер контейнера с постгрёй
docker compose --env-file .env up -d

Флаг -d в команде docker compose up -d означает запуск контейнеров в фоновом режиме (detached mode).
### Запуск приложения 
go run cmd/main.go 