# Использовать официальный образ postgres последней версии
FROM postgres:latest

# Можно добавить переменные окружения по умолчанию (опционально)
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=qwerty
ENV POSTGRES_DB=mydatabase

# Открыть порт 5432
EXPOSE 5431