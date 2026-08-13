### Запуск докер контейнера с постгрёй
docker compose up -d

Флаг -d в команде docker compose up -d означает запуск контейнеров в фоновом режиме (detached mode).
### Запуск приложения 
go run cmd/main.go 

### Генерация из .proto файла grpc-сервера
protoc --go_out=./protos/gen --go-grpc_out=./protos/gen --proto_path=./protos/proto protos/proto/library.proto

### Генерация из .proto файла grpc-gateway со стоковыми эндпоинтами
protoc -I . --go_out=./protos/gen --go-grpc_out=./protos/gen --proto_path=./protos/proto --grpc-gateway_out=allow_delete_body=true,paths=source_relative,generate_unbound_methods=true:./protos/gen library.proto

### Генерация из .proto файла grpc-gateway с ручными эндпоинтами
protoc -I . --go_out=. --go_opt=module=goproject --go-grpc_out=. --go-grpc_opt=module=goproject --grpc-gateway_out=logtostderr=true,allow_delete_body=true:. --grpc-gateway_opt=module=goproject protos/proto/library.proto

### Генерация из .proto файла описания Kafka эвента
protoc -I . --go_out=. --go_opt=module=goproject protos/proto/kafka_events.proto

#### Команда для просмотра всех топиков Kafka
docker exec -it kafka /opt/kafka/bin/kafka-topics.sh --list --bootstrap-server localhost:9092