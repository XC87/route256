# Выполнить миграцию
goose -dir loms/migrations postgres "postgresql://postgres:password@loms-db:5432/loms?sslmode=disable" up

# Запустить сервер
/loms-server
