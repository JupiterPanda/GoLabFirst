### Запуск докер контейнера с постгрёй
docker compose --env-file .env up -d

Флаг -d в команде docker compose up -d означает запуск контейнеров в фоновом режиме (detached mode).
### Запуск приложения 
go run cmd/main.go 

### Генерация из .proto файла grpc-сервера
protoc --go_out=./protos/gen --go-grpc_out=./protos/gen --proto_path=./protos/proto protos/proto/library.proto

### Генерация из .proto файла grpc-gateway со стоковыми эндпоинтами
protoc -I . --go_out=./protos/gen --go-grpc_out=./protos/gen --proto_path=./protos/proto --grpc-gateway_out=paths=source_relative,generate_unbound_methods=true:./protos/gen/librarypb library.proto

protoc -I . --proto_path=./protos/proto --openapiv2_out=./protos/gen --openapiv2_opt=logtostderr=true library.proto

### Генерация из .proto файла grpc-gateway с ручными эндпоинтами
protoc -I . --proto_path=./protos/proto --grpc-gateway_out ./protos/gen --grpc-gateway_opt paths=source_relative library.proto

### Генерация из .proto файла OpenAPI
protoc -I . --proto_path=./protos/proto --openapiv2_out ./protos/gen library.proto


protoc -I . --proto_path=./protos/proto --grpc-gateway_out ./protos/gen --grpc-gateway_opt paths=source_relative,allow_delete_body=true library.proto
