# Carrega variáveis do .env
include .env
export

.PHONY: up down logs

# Sobe a infra (banco + redis) em background
up:
	docker compose up -d

# Derruba tudo
down:
	docker compose down

# Vê logs dos containers
logs:
	docker compose logs -f

# Roda a aplicação Go
run:
	go run cmd/api/main.go

DB_URL=postgres://${DB_USER}:${DB_PASSWORD}@localhost:${DB_PORT}/${DB_NAME}?sslmode=disable

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

migrate-up:
	migrate -path migrations -database "${DB_URL}" -verbose up

migrate-down:
	migrate -path migrations -database "${DB_URL}" -verbose down

migrate-version:
	migrate -path migrations -database "${DB_URL}" version

migrate-force:
	migrate -path migrations -database "${DB_URL}" force $(version)